package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/mpacket"
)

// botAI handles bot behavior and movement (Phase 2+)
type botAI struct {
	bot *Player

	// Movement state
	movementEnabled bool
	moveDirection   int8 // -1 = left, 0 = stopped, 1 = right
	moveSpeed       int16
	lastMoveTime    time.Time
	nextActionTime  time.Time

	// Physics state
	velocityX float64 // Horizontal velocity
	velocityY float64 // Vertical velocity (for gravity/jumping)
	onGround  bool    // Is bot currently on a foothold?

	// Map boundaries (set when bot enters map)
	mapMinX, mapMaxX int16
	mapMinY, mapMaxY int16

	// Movement pattern
	walkDuration  time.Duration
	pauseDuration time.Duration
	shouldJump    bool
	jumpCooldown  time.Time

	// Physics constants
	gravity       float64 // Pixels per second squared
	jumpVelocity  float64 // Initial upward velocity on jump
	maxFallSpeed  float64 // Terminal velocity
	groundCheckDistance float64 // How far below to check for ground
}

// newBotAI creates a new AI controller for a bot
func newBotAI(bot *Player) *botAI {
	return &botAI{
		bot:             bot,
		movementEnabled: false,
		moveDirection:   0,
		moveSpeed:       100,             // Default walk speed (pixels/second)
		walkDuration:    time.Second * 3,
		pauseDuration:   time.Second * 2,
		nextActionTime:  time.Now().Add(time.Second * 2),
		
		// Physics constants (based on MapleStory physics)
		gravity:       670.0,  // Pixels per second squared
		jumpVelocity:  -350.0, // Initial upward velocity (negative = up)
		maxFallSpeed:  670.0,  // Terminal velocity
		groundCheckDistance: 10.0, // Check 10 pixels below for ground
		
		velocityX: 0,
		velocityY: 0,
		onGround:  true, // Assume spawned on ground
	}
}

// EnableMovement activates bot movement (Phase 2)
func (ai *botAI) EnableMovement(minX, maxX, minY, maxY int16) {
	ai.movementEnabled = true
	ai.mapMinX = minX
	ai.mapMaxX = maxX
	ai.mapMinY = minY
	ai.mapMaxY = maxY
	ai.nextActionTime = time.Now().Add(ai.pauseDuration)
	log.Printf("Bot %s movement enabled (bounds: x=%d-%d, y=%d-%d)",
		ai.bot.Name, minX, maxX, minY, maxY)
}

// DisableMovement stops bot movement
func (ai *botAI) DisableMovement() {
	ai.movementEnabled = false
	ai.moveDirection = 0
}

// Update is called periodically to update bot AI (should be called from a ticker)
func (ai *botAI) Update() {
	if !ai.movementEnabled || ai.bot.inst == nil {
		return
	}

	now := time.Now()

	// Check if it's time for next action
	if now.Before(ai.nextActionTime) {
		return
	}

	// Decide what to do next
	if ai.moveDirection == 0 {
		// Currently stopped, start moving
		ai.startWalking()
	} else {
		// Currently moving, stop
		ai.stopWalking()
	}
}

// startWalking makes the bot start walking in a random direction
func (ai *botAI) startWalking() {
	// Choose random direction based on position
	if ai.bot.pos.x <= ai.mapMinX+100 {
		ai.moveDirection = 1 // Force right if near left edge
	} else if ai.bot.pos.x >= ai.mapMaxX-100 {
		ai.moveDirection = -1 // Force left if near right edge
	} else {
		// Random direction
		if ai.bot.rng.Intn(2) == 0 {
			ai.moveDirection = -1
		} else {
			ai.moveDirection = 1
		}
	}

	// Random chance to jump while walking
	ai.shouldJump = ai.bot.rng.Intn(3) == 0 // 33% chance

	ai.nextActionTime = time.Now().Add(ai.walkDuration)
	ai.lastMoveTime = time.Now()
}

// stopWalking makes the bot stop moving
func (ai *botAI) stopWalking() {
	ai.moveDirection = 0
	ai.shouldJump = false
	ai.nextActionTime = time.Now().Add(ai.pauseDuration)
}

// PerformMovement executes physics-based movement and broadcasts to other players
func (ai *botAI) PerformMovement() {
	if !ai.movementEnabled || ai.bot.inst == nil {
		return
	}

	now := time.Now()
	deltaTime := now.Sub(ai.lastMoveTime).Seconds() // Use seconds for physics
	if deltaTime <= 0 || deltaTime > 0.5 { // Cap delta time to prevent huge jumps
		deltaTime = 0.1 // Default to 100ms
	}
	ai.lastMoveTime = now

	oldPos := ai.bot.pos

	// Apply gravity
	if !ai.onGround {
		ai.velocityY += ai.gravity * deltaTime
		// Cap fall speed
		if ai.velocityY > ai.maxFallSpeed {
			ai.velocityY = ai.maxFallSpeed
		}
	}

	// Horizontal movement
	if ai.moveDirection != 0 && ai.onGround {
		ai.velocityX = float64(ai.moveSpeed) * float64(ai.moveDirection)
	} else if ai.onGround {
		ai.velocityX = 0
	}

	// Calculate new position
	newX := float64(oldPos.x) + (ai.velocityX * deltaTime)
	newY := float64(oldPos.y) + (ai.velocityY * deltaTime)

	// Clamp X to map boundaries
	if newX < float64(ai.mapMinX) {
		newX = float64(ai.mapMinX)
		ai.moveDirection = 1 // Bounce off left wall
		ai.velocityX = 0
	} else if newX > float64(ai.mapMaxX) {
		newX = float64(ai.mapMaxX)
		ai.moveDirection = -1 // Bounce off right wall
		ai.velocityX = 0
	}

	// Ground collision detection
	// Check if there's a foothold at or below the new position
	testGroundPos := newPos(int16(newX), int16(newY), 0)
	groundPos := ai.bot.inst.fhHist.getFinalPosition(testGroundPos)

	// Check if we're on or very close to a foothold
	distanceToGround := float64(groundPos.y) - newY
	
	if distanceToGround >= 0 && distanceToGround <= ai.groundCheckDistance {
		// We're on the ground or very close to it
		ai.onGround = true
		ai.velocityY = 0
		newY = float64(groundPos.y) // Snap to foothold
		
		// Check if we should jump
		doJump := ai.shouldJump && now.After(ai.jumpCooldown) && ai.onGround
		if doJump {
			ai.velocityY = ai.jumpVelocity // Apply upward velocity
			ai.onGround = false
			ai.jumpCooldown = now.Add(time.Second * 2)
			ai.shouldJump = false
		}
	} else if distanceToGround < 0 {
		// We're above where we should be, need to fall
		ai.onGround = false
	} else {
		// We're too far above ground (falling)
		ai.onGround = false
	}

	// Update stance (facing direction)
	var stance byte
	if ai.moveDirection < 0 {
		stance = 1 // Facing left
	} else {
		stance = 0 // Facing right
	}

	// Create new position
	newPosition := newPos(int16(newX), int16(newY), groundPos.foothold)
	
	// Update bot position
	ai.bot.pos = newPosition
	ai.bot.stance = stance

	// Build movement packet
	moveData := ai.buildMovementPacket(oldPos, ai.bot.pos, stance, !ai.onGround, int16(deltaTime*1000))

	// Broadcast to other players
	ai.bot.inst.movePlayer(ai.bot.ID, moveData, ai.bot)
}

// buildMovementPacket creates movement packet data for the bot
func (ai *botAI) buildMovementPacket(fromPos, toPos pos, stance byte, jump bool, duration int16) []byte {
	p := mpacket.NewPacket()

	// Original position
	p.WriteInt16(fromPos.x)
	p.WriteInt16(fromPos.y)

	// Movement type and fragment
	var mType byte
	if jump {
		mType = 1 // Jump movement type
	} else {
		mType = 0 // Normal movement type
	}

	// Number of fragments (1 for simple movement)
	p.WriteByte(1)

	// Movement fragment
	p.WriteByte(mType)
	p.WriteInt16(toPos.x)
	p.WriteInt16(toPos.y)

	// Velocity (simplified)
	vx := int16(ai.moveSpeed * int16(ai.moveDirection))
	vy := int16(0)
	if jump {
		vy = -150 // Jump velocity
	}

	p.WriteInt16(vx)
	p.WriteInt16(vy)
	p.WriteInt16(toPos.foothold)
	p.WriteByte(stance)
	p.WriteInt16(duration)

	return p
}

func directionStr(dir int8) string {
	switch dir {
	case -1:
		return "left"
	case 1:
		return "right"
	default:
		return "stopped"
	}
}

// GetMovementState returns current movement state for debugging
func (ai *botAI) GetMovementState() string {
	if !ai.movementEnabled {
		return "disabled"
	}
	if ai.moveDirection == 0 {
		return "stopped"
	}
	return directionStr(ai.moveDirection)
}
