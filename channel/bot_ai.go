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

	// Map boundaries (set when bot enters map)
	mapMinX, mapMaxX int16
	mapMinY, mapMaxY int16

	// Movement pattern
	walkDuration   time.Duration
	pauseDuration  time.Duration
	shouldJump     bool
	jumpCooldown   time.Time
}

// newBotAI creates a new AI controller for a bot
func newBotAI(bot *Player) *botAI {
	return &botAI{
		bot:             bot,
		movementEnabled: false,
		moveDirection:   0,
		moveSpeed:       100, // Default walk speed
		walkDuration:    time.Second * 3,
		pauseDuration:   time.Second * 2,
		nextActionTime:  time.Now().Add(time.Second * 2),
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

	log.Printf("Bot %s starts walking %s at position (%d, %d)",
		ai.bot.Name, directionStr(ai.moveDirection), ai.bot.pos.x, ai.bot.pos.y)
}

// stopWalking makes the bot stop moving
func (ai *botAI) stopWalking() {
	ai.moveDirection = 0
	ai.shouldJump = false
	ai.nextActionTime = time.Now().Add(ai.pauseDuration)

	log.Printf("Bot %s stops walking at position (%d, %d)",
		ai.bot.Name, ai.bot.pos.x, ai.bot.pos.y)
}

// PerformMovement executes movement and broadcasts to other players
func (ai *botAI) PerformMovement() {
	if !ai.movementEnabled || ai.moveDirection == 0 || ai.bot.inst == nil {
		return
	}

	now := time.Now()
	deltaTime := now.Sub(ai.lastMoveTime).Milliseconds()
	if deltaTime <= 0 {
		return
	}
	ai.lastMoveTime = now

	// Calculate movement
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

	// Update stance (facing direction)
	var stance byte
	if ai.moveDirection < 0 {
		stance = 5 // Walking left
	} else {
		stance = 2 // Walking right
	}

	// Check if we should jump
	doJump := ai.shouldJump && now.After(ai.jumpCooldown)
	if doJump {
		ai.jumpCooldown = now.Add(time.Second * 2) // Jump cooldown
	}

	// Calculate new position with foothold
	oldPos := ai.bot.pos
	tempPos := newPos(newX, oldPos.y, oldPos.foothold)

	// Get correct foothold for new position
	if ai.bot.inst != nil {
		newPos := ai.bot.inst.fhHist.getFinalPosition(tempPos)
		ai.bot.pos = newPos
	} else {
		ai.bot.pos = tempPos
	}

	ai.bot.stance = stance

	// Build movement packet
	moveData := ai.buildMovementPacket(oldPos, ai.bot.pos, stance, doJump, int16(deltaTime))

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
