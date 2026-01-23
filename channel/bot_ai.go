package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/mpacket"
)

// PlayerState represents the bot's current state (from MapleStory client PlayerStates)
type PlayerState int

const (
	StateStanding PlayerState = iota
	StateWalking
	StateFalling
	StateJumping
)

// botAI handles bot behavior and movement (Phase 2+)
// Implements MapleStory client-style physics system
type botAI struct {
	bot *Player

	// AI decision making
	movementEnabled bool
	nextActionTime  time.Time
	walkDuration    time.Duration
	pauseDuration   time.Duration

	// Player state (from MapleStory client)
	state         PlayerState
	facingLeft    bool // true = left, false = right
	lastMoveTime  time.Time

	// Physics state (PhysicsObject from MapleStory client)
	hspeed    float64 // Horizontal velocity
	vspeed    float64 // Vertical velocity
	x         float64 // Precise X position (float for sub-pixel accuracy)
	y         float64 // Precise Y position
	fhid      int16   // Current foothold ID
	fhslope   float64 // Current foothold slope
	fhlayer   int8    // Current foothold layer
	onground  bool    // Is on ground
	canjump   bool    // Can jump (prevents double jump)
	
	// Map boundaries (set when bot enters map)
	mapMinX, mapMaxX int16
	mapMinY, mapMaxY int16
}

// newBotAI creates a new AI controller for a bot
func newBotAI(bot *Player) *botAI {
	return &botAI{
		bot:             bot,
		movementEnabled: false,
		walkDuration:    time.Second * 3,
		pauseDuration:   time.Second * 2,
		nextActionTime:  time.Now().Add(time.Second * 2),
		
		// Initialize physics state
		state:    StateStanding,
		facingLeft: false,
		hspeed:   0,
		vspeed:   0,
		x:        float64(bot.pos.x),
		y:        float64(bot.pos.y),
		fhid:     bot.pos.foothold,
		fhslope:  0,
		fhlayer:  0,
		onground: true,
		canjump:  true,
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
	ai.state = StateStanding
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
	if ai.state == StateStanding {
		// Currently stopped, start walking
		ai.startWalking()
	} else if ai.state == StateWalking {
		// Currently walking, stop
		ai.stopWalking()
	}
}

// startWalking makes the bot start walking in a random direction
func (ai *botAI) startWalking() {
	// Choose random direction based on position
	if ai.bot.pos.x <= ai.mapMinX+100 {
		ai.facingLeft = false // Force right if near left edge
	} else if ai.bot.pos.x >= ai.mapMaxX-100 {
		ai.facingLeft = true // Force left if near right edge
	} else {
		// Random direction
		ai.facingLeft = ai.bot.rng.Intn(2) == 0
	}

	ai.state = StateWalking
	
	// Random chance to jump while walking
	if ai.bot.rng.Intn(3) == 0 { // 33% chance
		ai.tryJump()
	}

	ai.nextActionTime = time.Now().Add(ai.walkDuration)
	ai.lastMoveTime = time.Now()
}

// stopWalking makes the bot stop moving
func (ai *botAI) stopWalking() {
	ai.state = StateStanding
	ai.nextActionTime = time.Now().Add(ai.pauseDuration)
}

// tryJump attempts to make the bot jump
func (ai *botAI) tryJump() {
	if ai.onground && ai.canjump && (ai.state == StateStanding || ai.state == StateWalking) {
		ai.state = StateJumping
		ai.vspeed = JUMPFORCE
		ai.canjump = false
	}
}

// Physics constants from MapleStory client (from Physics.cpp)
const (
	GRAVFORCE      = 0.35  // Gravity acceleration per frame
	FRICTION       = 0.3   // Ground friction
	WALKFORCE      = 0.7   // Walking acceleration
	WALKSPEED      = 1.5   // Maximum walk speed
	JUMPFORCE      = -5.5  // Initial jump force (negative = upward)
	MAXVERTSPEED   = 8.0   // Terminal velocity (max fall speed)
	GROUNDTHRESHOLD = 5.0  // Distance tolerance for ground detection (pixels)
)

// PerformMovement executes one physics update cycle (from MapleStory client move_object)
// Order: update_fh → apply_physics → move
func (ai *botAI) PerformMovement() {
	if !ai.movementEnabled || ai.bot.inst == nil {
		return
	}

	now := time.Now()
	deltaTime := now.Sub(ai.lastMoveTime).Seconds()
	if deltaTime <= 0 || deltaTime > 0.5 {
		deltaTime = 0.1 // Default to 100ms if first call or too large
	}
	ai.lastMoveTime = now

	oldPos := ai.bot.pos

	// === STEP 1: Update foothold (from FootholdTree.cpp) ===
	ai.updateFoothold()

	// === STEP 2: Apply physics based on state (from Physics.cpp) ===
	ai.applyPhysics(deltaTime)

	// === STEP 3: Move (apply velocities to position) ===
	ai.move(deltaTime)

	// === STEP 4: Update foothold again after movement ===
	ai.updateFoothold()

	// === STEP 5: Check state transitions ===
	ai.updateState()

	// === STEP 6: Sync bot position to player struct ===
	ai.bot.pos.x = int16(ai.x)
	ai.bot.pos.y = int16(ai.y)
	ai.bot.pos.foothold = ai.fhid

	// Update stance (facing direction)
	var stance byte
	if ai.facingLeft {
		stance = 1
	} else {
		stance = 0
	}
	ai.bot.stance = stance

	// Build movement packet
	moveData := ai.buildMovementPacket(oldPos, ai.bot.pos, stance, !ai.onground, int16(deltaTime*1000))

	// Broadcast to other players
	ai.bot.inst.movePlayer(ai.bot.ID, moveData, ai.bot)
}

// applyPhysics calculates forces and accelerations (from Physics.cpp move_normal)
func (ai *botAI) applyPhysics(dt float64) {
	// Reset accelerations
	hacc := 0.0
	vacc := 0.0

	if ai.onground {
		// === ON GROUND PHYSICS ===
		
		// Apply horizontal walking force
		if ai.state == StateWalking {
			walkdir := 1.0
			if ai.facingLeft {
				walkdir = -1.0
			}
			hacc += WALKFORCE * walkdir
		}
		
		// Apply friction
		if ai.hspeed != 0 {
			friction := FRICTION
			if ai.hspeed > 0 {
				hacc -= friction
			} else {
				hacc += friction
			}
		}
		
		// Apply slope force (simplified - assuming flat ground for now)
		// TODO: Get actual slope from foothold and apply slope forces
		
	} else {
		// === IN AIR PHYSICS ===
		
		// Apply gravity
		vacc += GRAVFORCE
		
		// Minimal air resistance on horizontal movement
		if ai.hspeed != 0 {
			if ai.hspeed > 0 {
				hacc -= FRICTION * 0.1 // Less friction in air
			} else {
				hacc += FRICTION * 0.1
			}
		}
	}
	
	// Update velocities
	ai.hspeed += hacc
	ai.vspeed += vacc
	
	// Apply speed caps
	if ai.hspeed > WALKSPEED {
		ai.hspeed = WALKSPEED
	} else if ai.hspeed < -WALKSPEED {
		ai.hspeed = -WALKSPEED
	}
	
	// Enforce terminal velocity
	if ai.vspeed > MAXVERTSPEED {
		ai.vspeed = MAXVERTSPEED
	}
	
	// Stop if speed is very small (prevents jitter)
	if abs(ai.hspeed) < 0.01 {
		ai.hspeed = 0
	}
}

// move applies velocities to position (from PhysicsObject.move())
func (ai *botAI) move(dt float64) {
	// Update position based on velocities
	ai.x += ai.hspeed
	ai.y += ai.vspeed
	
	// Clamp to map boundaries
	if ai.x < float64(ai.mapMinX) {
		ai.x = float64(ai.mapMinX)
		ai.hspeed = 0
		ai.facingLeft = false // Bounce - face right
	} else if ai.x > float64(ai.mapMaxX) {
		ai.x = float64(ai.mapMaxX)
		ai.hspeed = 0
		ai.facingLeft = true // Bounce - face left
	}
	
	if ai.y < float64(ai.mapMinY) {
		ai.y = float64(ai.mapMinY)
		ai.vspeed = 0
	} else if ai.y > float64(ai.mapMaxY) {
		ai.y = float64(ai.mapMaxY)
		ai.vspeed = 0
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// updateFoothold updates the bot's foothold and ground state (from FootholdTree.cpp update_fh)
func (ai *botAI) updateFoothold() {
	if ai.bot.inst == nil {
		return
	}

	// Get foothold at current position
	testPos := newPos(int16(ai.x), int16(ai.y), 0)
	groundPos := ai.bot.inst.fhHist.getFinalPosition(testPos)

	// Check if valid foothold exists
	if groundPos.foothold == 0 {
		// No foothold - in air
		ai.onground = false
		ai.fhid = 0
		ai.fhslope = 0.0
		ai.canjump = false
		return
	}

	// Calculate vertical distance to ground
	distToGround := ai.y - float64(groundPos.y)

	if abs(distToGround) <= GROUNDTHRESHOLD {
		// On or very close to ground - snap to it
		ai.onground = true
		ai.fhid = groundPos.foothold
		ai.y = float64(groundPos.y) // Snap to ground
		ai.vspeed = 0
		ai.canjump = true
		
		// TODO: Calculate actual slope from foothold data
		ai.fhslope = 0.0
		
	} else if distToGround < 0 {
		// Above ground - falling or jumping
		ai.onground = false
		ai.canjump = false
		
	} else {
		// Below ground (clipped through) - emergency snap up
		ai.onground = true
		ai.fhid = groundPos.foothold
		ai.y = float64(groundPos.y)
		ai.vspeed = 0
		ai.canjump = true
	}
}

// updateState transitions between player states based on physics
func (ai *botAI) updateState() {
	switch ai.state {
	case StateJumping:
		// Transition from jumping to falling when moving downward
		if ai.vspeed > 0 {
			ai.state = StateFalling
		}
		// If landed, will transition in next cycle after updateFoothold
		if ai.onground {
			if ai.hspeed != 0 {
				ai.state = StateWalking
			} else {
				ai.state = StateStanding
			}
		}
		
	case StateFalling:
		// Transition to standing/walking when landed
		if ai.onground {
			if ai.hspeed != 0 {
				ai.state = StateWalking
			} else {
				ai.state = StateStanding
			}
		}
		
	case StateWalking:
		// Check if we walked off an edge
		if !ai.onground {
			ai.state = StateFalling
		}
		// Check if stopped moving
		if ai.hspeed == 0 {
			ai.state = StateStanding
		}
		
	case StateStanding:
		// Check if fell off edge while standing
		if !ai.onground {
			ai.state = StateFalling
		}
		// Check if started moving
		if ai.hspeed != 0 {
			ai.state = StateWalking
		}
	}
}

// buildMovementPacket creates movement packet data for the bot
func (ai *botAI) buildMovementPacket(fromPos, toPos pos, stance byte, inAir bool, duration int16) []byte {
	p := mpacket.NewPacket()

	// Original position
	p.WriteInt16(fromPos.x)
	p.WriteInt16(fromPos.y)

	// Movement type based on state
	var mType byte
	switch ai.state {
	case StateJumping, StateFalling:
		mType = 1 // Jump/fall movement type
	case StateWalking:
		mType = 0 // Walking movement type
	default:
		mType = 0 // Standing
	}

	// Number of fragments (1 for simple movement)
	p.WriteByte(1)

	// Movement fragment
	p.WriteByte(mType)
	p.WriteInt16(toPos.x)
	p.WriteInt16(toPos.y)

	// Velocity (use actual physics velocity)
	vx := int16(ai.hspeed * 50) // Scale for packet
	vy := int16(ai.vspeed * 50)

	p.WriteInt16(vx)
	p.WriteInt16(vy)
	p.WriteInt16(toPos.foothold)
	p.WriteByte(stance)
	p.WriteInt16(duration)

	return p
}

// GetMovementState returns current movement state for debugging
func (ai *botAI) GetMovementState() string {
	if !ai.movementEnabled {
		return "disabled"
	}
	switch ai.state {
	case StateStanding:
		return "standing"
	case StateWalking:
		if ai.facingLeft {
			return "walking left"
		}
		return "walking right"
	case StateJumping:
		return "jumping"
	case StateFalling:
		return "falling"
	default:
		return "unknown"
	}
}
