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

	// Physics state (based on MapleStory client physics)
	hspeed float64 // Horizontal speed
	vspeed float64 // Vertical speed
	hforce float64 // Horizontal force
	vforce float64 // Vertical force
	hacc   float64 // Horizontal acceleration
	vacc   float64 // Vertical acceleration
	
	fhslope  float64 // Current foothold slope
	onground bool    // Is bot on ground

	// Map boundaries (set when bot enters map)
	mapMinX, mapMaxX int16
	mapMinY, mapMaxY int16

	// Movement pattern
	walkDuration  time.Duration
	pauseDuration time.Duration
	shouldJump    bool
	jumpCooldown  time.Time
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
		
		// Physics state
		hspeed: 0,
		vspeed: 0,
		hforce: 0,
		vforce: 0,
		hacc:   0,
		vacc:   0,
		fhslope:  0,
		onground: true, // Assume spawned on ground
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

// Physics constants from MapleStory client
const (
	GRAVFORCE    = 2.0  // Increased for more immediate gravity
	FRICTION     = 0.5
	SLOPEFACTOR  = 0.1
	GROUNDSLIP   = 3.0
	WALKFORCE    = 0.14 // Force applied when walking
	JUMPFORCE    = -15.0 // Increased jump force to compensate for stronger gravity
)

// PerformMovement executes movement using simple X movement from 4eed6f0 + physics Y movement
func (ai *botAI) PerformMovement() {
	if !ai.movementEnabled || ai.bot.inst == nil {
		return
	}

	now := time.Now()
	deltaTime := now.Sub(ai.lastMoveTime).Milliseconds()
	if deltaTime <= 0 {
		deltaTime = 100 // Default to 100ms if first call
	}
	ai.lastMoveTime = now

	oldPos := ai.bot.pos

	// === X MOVEMENT: Simple distance-based (from 4eed6f0 - this worked!) ===
	if ai.moveDirection != 0 {
		distance := int16(float64(ai.moveSpeed) * float64(deltaTime) / 1000.0)
		newX := ai.bot.pos.x + (distance * int16(ai.moveDirection))

		// Clamp to map boundaries
		if newX < ai.mapMinX {
			newX = ai.mapMinX
			ai.moveDirection = 1 // Bounce off left wall
		} else if newX > ai.mapMaxX {
			newX = ai.mapMaxX
			ai.moveDirection = -1 // Bounce off right wall
		}

		ai.bot.pos.x = newX
	}

	// === Y MOVEMENT: Physics-based with gravity (current working logic) ===
	// Update foothold information
	ai.updateFoothold()

	// Apply physics for vertical movement
	ai.applyPhysics()

	// Apply vertical speed to position
	newY := float64(ai.bot.pos.y) + ai.vspeed
	ai.bot.pos.y = int16(newY)

	// Update foothold after movement (snap to ground if needed)
	ai.updateFoothold()

	// Handle jumping
	if ai.shouldJump && now.After(ai.jumpCooldown) && ai.onground {
		ai.vforce = JUMPFORCE
		ai.shouldJump = false
		ai.jumpCooldown = now.Add(time.Second * 2)
	}

	// Update stance (facing direction)
	var stance byte
	if ai.moveDirection < 0 {
		stance = 1 // Facing left
	} else {
		stance = 0 // Facing right
	}
	ai.bot.stance = stance

	// Build movement packet
	moveData := ai.buildMovementPacket(oldPos, ai.bot.pos, stance, !ai.onground, int16(deltaTime))

	// Broadcast to other players
	ai.bot.inst.movePlayer(ai.bot.ID, moveData, ai.bot)
}

// applyPhysics applies physics calculations for Y movement only (X uses simple distance)
func (ai *botAI) applyPhysics() {
	ai.vacc = 0.0

	if ai.onground {
		// On ground - apply vertical forces (e.g., jump)
		ai.vacc += ai.vforce
	} else {
		// In air - apply gravity
		ai.vacc += GRAVFORCE
		ai.vacc += ai.vforce // Also apply any jump force that's active
	}

	// Reset vertical force after using it
	ai.vforce = 0.0

	// Update vertical speed
	ai.vspeed += ai.vacc
}

// updateFoothold updates the bot's foothold and checks if on ground
func (ai *botAI) updateFoothold() {
	if ai.bot.inst == nil {
		return
	}

	// Get foothold at current position
	testPos := newPos(ai.bot.pos.x, ai.bot.pos.y, 0)
	groundPos := ai.bot.inst.fhHist.getFinalPosition(testPos)

	// Check if valid foothold
	if groundPos.foothold == 0 {
		ai.onground = false
		ai.fhslope = 0.0
		return
	}

	// Update foothold
	ai.bot.pos.foothold = groundPos.foothold
	
	// Calculate ground level for this foothold
	// The foothold system returns the Y position where ground is
	groundY := groundPos.y

	// Check if we're on or near the ground (larger tolerance to catch fast falling)
	distanceToGround := ai.bot.pos.y - groundY
	
	if distanceToGround >= -10 && distanceToGround <= 10 {
		// On ground or very close - snap to ground
		ai.onground = true
		ai.bot.pos.y = groundY // Snap to ground
		ai.vspeed = 0
		
		// Calculate slope (simplified - would need actual foothold data)
		ai.fhslope = 0.0 // TODO: Get actual slope from foothold
	} else if ai.bot.pos.y < groundY {
		// Above ground - falling
		ai.onground = false
	} else {
		// Below ground (fell past foothold) - snap up to ground
		ai.bot.pos.y = groundY
		ai.onground = true
		ai.vspeed = 0
	}
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
