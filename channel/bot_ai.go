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
	randomSeed      int64 // Per-bot random seed for unique behavior

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
	// Create a unique random seed for this bot based on spawn time and character ID
	randomSeed := time.Now().UnixNano() + int64(bot.ID)
	
	// Use the bot's existing RNG to generate HIGHLY VARIED durations per bot
	// Wide range ensures each bot has distinctly different timing
	walkDur := time.Second * time.Duration(1+bot.rng.Intn(8))   // 1-8 seconds (8x variance)
	pauseDur := time.Second * time.Duration(1+bot.rng.Intn(6))  // 1-6 seconds (6x variance)
	
	return &botAI{
		bot:             bot,
		movementEnabled: false,
		walkDuration:    walkDur,
		pauseDuration:   pauseDur,
		nextActionTime:  time.Now().Add(pauseDur),
		randomSeed:      randomSeed,
		
		// Initialize physics state
		state:    StateStanding,
		facingLeft: bot.rng.Intn(2) == 0, // Random initial facing
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
	if now.After(ai.nextActionTime) || now.Equal(ai.nextActionTime) {
		// Decide what to do next
		if ai.state == StateStanding {
			// Currently stopped, start walking
			ai.startWalking()
		} else if ai.state == StateWalking {
			// Currently walking, stop
			ai.stopWalking()
		}
	}
}

// startWalking makes the bot start walking in a random direction
func (ai *botAI) startWalking() {
	// Choose TRULY RANDOM direction with UNIQUE per-bot bias
	// This prevents all bots from walking the same way
	// Each bot has a persistent direction preference based on their random seed
	directionRoll := ai.bot.rng.Intn(100) // 0-99
	
	// Use bot's randomSeed to create unique preference (some bots prefer left, some right)
	preferenceThreshold := int((ai.randomSeed % 60) + 20) // 20-79 unique per bot
	
	// Near edges, add small bias but don't force
	if ai.bot.pos.x <= ai.mapMinX+100 {
		preferenceThreshold -= 15 // Slightly prefer right
	} else if ai.bot.pos.x >= ai.mapMaxX-100 {
		preferenceThreshold += 15 // Slightly prefer left
	}
	
	// Decide direction based on preference
	if directionRoll < preferenceThreshold {
		ai.facingLeft = true
	} else {
		ai.facingLeft = false
	}

	ai.state = StateWalking
	
	// HIGHLY VARIED jump chance per bot (10-50% range)
	jumpChance := ai.bot.rng.Intn(10) // 0-9
	jumpThreshold := 2 + ai.bot.rng.Intn(3) // Each bot has threshold of 2-4 (20-40% chance)
	if jumpChance < jumpThreshold {
		ai.tryJump()
	}

	// LARGE variation in walk duration for each walk cycle
	variation := time.Millisecond * time.Duration(ai.bot.rng.Intn(3000)) // +/- up to 3 seconds
	ai.nextActionTime = time.Now().Add(ai.walkDuration + variation)
	ai.lastMoveTime = time.Now()
}

// stopWalking makes the bot stop moving
func (ai *botAI) stopWalking() {
	ai.state = StateStanding
	// LARGE variation in pause duration for each stop
	variation := time.Millisecond * time.Duration(ai.bot.rng.Intn(2000)) // +/- up to 2 seconds
	ai.nextActionTime = time.Now().Add(ai.pauseDuration + variation)
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
// Note: Speeds are tuned for 10 FPS server update rate
const (
	GRAVFORCE      = 0.35  // Gravity acceleration per frame
	FRICTION       = 0.3   // Ground friction
	WALKFORCE      = 1.5   // Walking acceleration (increased for more responsive movement)
	WALKSPEED      = 10.0  // Maximum walk speed (increased for visible movement at 10 FPS)
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
	// If bot is standing still (not Walking/Jumping/Falling), don't apply physics forces
	if ai.state == StateStanding {
		// Apply strong friction to stop quickly
		if ai.hspeed != 0 {
			friction := FRICTION * 3.0 // Stronger friction when standing
			if ai.hspeed > 0 {
				ai.hspeed -= friction
				if ai.hspeed < 0 {
					ai.hspeed = 0
				}
			} else {
				ai.hspeed += friction
				if ai.hspeed > 0 {
					ai.hspeed = 0
				}
			}
		}
		return
	}
	
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
		
		// Apply gravity - NO horizontal air resistance for straight drops
		vacc += GRAVFORCE
		
		// No air resistance on horizontal movement - let it maintain momentum
		// This prevents the "sliding down hill" appearance
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
// Implements wall and edge collision detection from FootholdTree.cpp
func (ai *botAI) move(dt float64) {
	// Store current position before movement
	crntX := ai.x
	crntY := ai.y
	
	// Calculate next position
	nextX := ai.x + ai.hspeed
	nextY := ai.y + ai.vspeed
	
	// === HORIZONTAL COLLISION DETECTION (from FootholdTree.cpp limit_movement) ===
	if ai.hspeed != 0 {
		// Check for walls and edges
		movingLeft := ai.hspeed < 0
		
		// Get wall or edge position in movement direction
		wallOrEdgeInfo := ai.getWallOrEdge(movingLeft, nextY)
		wallOrEdge := wallOrEdgeInfo.position
		isEdge := wallOrEdgeInfo.isEdge
		platformHeight := wallOrEdgeInfo.platformHeight
		
		// Check if we'll collide with wall/edge
		collision := false
		if movingLeft {
			collision = crntX >= wallOrEdge && nextX <= wallOrEdge
		} else {
			collision = crntX <= wallOrEdge && nextX >= wallOrEdge
		}
		
		if collision {
			// Stop at wall/edge
			nextX = wallOrEdge
			ai.hspeed = 0
			
			// If it's an edge (not a wall), intelligently decide what to do
			if isEdge && ai.onground && ai.canjump {
				heightDiff := platformHeight - crntY
				
				// If platform is above us (climbing up), ALWAYS try to jump
				if heightDiff < 0 && heightDiff > -80 { // Platform is up to 80 pixels higher
					// Jump to try to reach the platform
					ai.vspeed = JUMPFORCE
					ai.canjump = false
					ai.state = StateJumping
					log.Printf("Bot %s jumping up to platform (heightDiff: %.1f)", ai.bot.Name, heightDiff)
					// Don't reverse direction - we're jumping forward
					return // Exit early so we don't reverse direction
				} else if heightDiff > 0 && heightDiff < 100 { // Platform is slightly below (drop down)
					// Just walk off and fall naturally
					log.Printf("Bot %s walking off edge to drop down (heightDiff: %.1f)", ai.bot.Name, heightDiff)
					// Don't reverse direction, let gravity handle it
					return // Exit early so we don't reverse direction
				} else {
					// Can't reach platform or too far down, reverse direction
					log.Printf("Bot %s reversing at unreachable platform (heightDiff: %.1f)", ai.bot.Name, heightDiff)
					ai.facingLeft = !ai.facingLeft
					// Stop walking state to prevent walking in place
					if ai.state == StateWalking {
						ai.state = StateStanding
					}
				}
			} else {
				// It's a wall or can't jump, reverse direction
				log.Printf("Bot %s reversing at wall (isEdge: %v, onground: %v, canjump: %v)", 
					ai.bot.Name, isEdge, ai.onground, ai.canjump)
				ai.facingLeft = !ai.facingLeft
				// Stop walking state to prevent walking in place
				if ai.state == StateWalking {
					ai.state = StateStanding
				}
			}
		}
	}
	
	// === VERTICAL COLLISION DETECTION (from FootholdTree.cpp limit_movement) ===
	if ai.vspeed != 0 && ai.vspeed > 0 { // Only check when falling (vspeed > 0)
		// Check ground collision at current and next X positions
		// This handles landing on sloped platforms correctly
		groundAtCrnt := ai.getGroundBelow(crntX)
		groundAtNext := ai.getGroundBelow(nextX)
		
		// Check if we're crossing through the ground
		collision := crntY <= groundAtCrnt && nextY >= groundAtNext
		
		if collision {
			// Land on ground
			nextY = groundAtNext
			ai.vspeed = 0
		}
	}
	
	// Clamp to absolute map boundaries as final safety
	if nextX < float64(ai.mapMinX) {
		nextX = float64(ai.mapMinX)
		ai.hspeed = 0
		ai.facingLeft = false
	} else if nextX > float64(ai.mapMaxX) {
		nextX = float64(ai.mapMaxX)
		ai.hspeed = 0
		ai.facingLeft = true
	}
	
	if nextY < float64(ai.mapMinY) {
		nextY = float64(ai.mapMinY)
		ai.vspeed = 0
	} else if nextY > float64(ai.mapMaxY) {
		nextY = float64(ai.mapMaxY)
		ai.vspeed = 0
	}
	
	// Apply movement
	ai.x = nextX
	ai.y = nextY
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// wallOrEdgeInfo contains collision detection information
type wallOrEdgeInfo struct {
	position       float64 // X position of the wall/edge
	isEdge         bool    // true if it's an edge (can jump), false if it's a wall
	platformHeight float64 // Y position of the platform (if it's an edge)
}

// getWallOrEdge returns the wall or edge position in the given direction
// Implements FootholdTree.cpp get_wall() and get_edge() logic
func (ai *botAI) getWallOrEdge(left bool, nextY float64) wallOrEdgeInfo {
	if ai.bot.inst == nil || ai.fhid == 0 {
		// No foothold data, use map boundaries (treat as wall)
		if left {
			return wallOrEdgeInfo{
				position: float64(ai.mapMinX),
				isEdge:   false,
				platformHeight: nextY,
			}
		}
		return wallOrEdgeInfo{
			position: float64(ai.mapMaxX),
			isEdge:   false,
			platformHeight: nextY,
		}
	}
	
	// Try to detect walls by checking adjacent footholds
	// In the real client, this checks if adjacent footholds are walls (vertical)
	// For now, we'll do edge detection by trying to get ground at a test position
	
	if left {
		// Check left - try to move left and see if we can find ground
		testX := ai.x - 50 // Test 50 pixels to the left
		testPos := newPos(int16(testX), int16(nextY), 0)
		groundPos := ai.bot.inst.fhHist.getFinalPosition(testPos)
		
		// If no ground found or ground is way below, treat as edge
		if groundPos.foothold == 0 {
			// No ground - it's an edge (or wall)
			return wallOrEdgeInfo{
				position:       ai.x,
				isEdge:         true,
				platformHeight: float64(ai.mapMaxY), // Unknown platform
			}
		}
		
		yDiff := abs(float64(groundPos.y) - nextY)
		if yDiff > 100 {
			// Ground is far away - treat as edge
			return wallOrEdgeInfo{
				position:       ai.x,
				isEdge:         true,
				platformHeight: float64(groundPos.y),
			}
		}
		
		// Ground is close - can continue walking
		return wallOrEdgeInfo{
			position:       float64(ai.mapMinX),
			isEdge:         false,
			platformHeight: float64(groundPos.y),
		}
	} else {
		// Check right - try to move right and see if we can find ground
		testX := ai.x + 50 // Test 50 pixels to the right
		testPos := newPos(int16(testX), int16(nextY), 0)
		groundPos := ai.bot.inst.fhHist.getFinalPosition(testPos)
		
		// If no ground found or ground is way below, treat as edge
		if groundPos.foothold == 0 {
			// No ground - it's an edge (or wall)
			return wallOrEdgeInfo{
				position:       ai.x,
				isEdge:         true,
				platformHeight: float64(ai.mapMaxY), // Unknown platform
			}
		}
		
		yDiff := abs(float64(groundPos.y) - nextY)
		if yDiff > 100 {
			// Ground is far away - treat as edge
			return wallOrEdgeInfo{
				position:       ai.x,
				isEdge:         true,
				platformHeight: float64(groundPos.y),
			}
		}
		
		// Ground is close - can continue walking
		return wallOrEdgeInfo{
			position:       float64(ai.mapMaxX),
			isEdge:         false,
			platformHeight: float64(groundPos.y),
		}
	}
}

// getGroundBelow returns the Y coordinate of the ground below the given X position
// Implements FootholdTree.cpp get_fhid_below() and Foothold.cpp ground_below() logic
func (ai *botAI) getGroundBelow(x float64) float64 {
	if ai.bot.inst == nil {
		return float64(ai.mapMaxY)
	}
	
	// Use the foothold system to find ground
	testPos := newPos(int16(x), int16(ai.y), 0)
	groundPos := ai.bot.inst.fhHist.getFinalPosition(testPos)
	
	if groundPos.foothold == 0 {
		// No foothold found - return map bottom
		return float64(ai.mapMaxY)
	}
	
	return float64(groundPos.y)
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
