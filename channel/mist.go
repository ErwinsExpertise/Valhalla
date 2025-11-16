package channel

import (
	"math"
	"time"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

// fieldMist represents a poison mist or other area effect on the field
type fieldMist struct {
	ID           int32
	ownerID      int32
	skillID      int32
	skillLevel   byte
	box          mistBox // Rectangle area of effect
	createdAt    time.Time
	duration     int64 // in seconds
	isPoisonMist bool
}

// mistBox defines the rectangular area of a mist
type mistBox struct {
	x1, y1 int16 // Top-left corner
	x2, y2 int16 // Bottom-right corner
}

// mistPool manages all mists in a field instance
type mistPool struct {
	instance *fieldInstance
	poolID   int32
	mists    map[int32]*fieldMist
}

func createNewMistPool(inst *fieldInstance) mistPool {
	return mistPool{
		instance: inst,
		mists:    make(map[int32]*fieldMist),
	}
}

func (pool *mistPool) nextID() int32 {
	for i := 0; i < 100; i++ {
		pool.poolID++
		if pool.poolID == math.MaxInt32-1 {
			pool.poolID = math.MaxInt32 / 2
		} else if pool.poolID == 0 {
			pool.poolID = 1
		}

		if _, ok := pool.mists[pool.poolID]; !ok {
			return pool.poolID
		}
	}
	return 0
}

// createMist spawns a new mist on the field
func (pool *mistPool) createMist(ownerID, skillID int32, skillLevel byte, pos pos, duration int64, isPoisonMist bool) *fieldMist {
	mistID := pool.nextID()
	if mistID == 0 {
		return nil
	}

	// Define the mist box (rectangular area)
	// For Poison Mist, typical size is about 300x200 pixels
	const mistWidth int16 = 150  // radius in x direction
	const mistHeight int16 = 100 // radius in y direction

	mist := &fieldMist{
		ID:         mistID,
		ownerID:    ownerID,
		skillID:    skillID,
		skillLevel: skillLevel,
		box: mistBox{
			x1: pos.x - mistWidth,
			y1: pos.y - mistHeight,
			x2: pos.x + mistWidth,
			y2: pos.y + mistHeight,
		},
		createdAt:    time.Now(),
		duration:     duration,
		isPoisonMist: isPoisonMist,
	}

	pool.mists[mistID] = mist

	// Send spawn packet to all players
	pool.instance.send(packetMistSpawn(mist))

	// Schedule removal after duration
	if duration > 0 {
		go func() {
			time.Sleep(time.Duration(duration) * time.Second)
			pool.instance.dispatch <- func() {
				pool.removeMist(mistID)
			}
		}()
	}

	return mist
}

// removeMist removes a mist from the field
func (pool *mistPool) removeMist(mistID int32) {
	if mist, ok := pool.mists[mistID]; ok {
		// Send removal packet
		pool.instance.send(packetMistRemove(mistID, mist.isPoisonMist))
		delete(pool.mists, mistID)
	}
}

// playerShowMists shows all active mists to a player joining the map
func (pool mistPool) playerShowMists(plr *Player) {
	for _, mist := range pool.mists {
		plr.Send(packetMistSpawn(mist))
	}
}

// isInMist checks if a position is within a mist's area
func (m *fieldMist) isInMist(p pos) bool {
	return p.x >= m.box.x1 && p.x <= m.box.x2 && p.y >= m.box.y1 && p.y <= m.box.y2
}

// Packet functions for mist

func packetMistSpawn(mist *fieldMist) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelAffectedAreaCreate)
	p.WriteInt32(mist.ID)
	p.WriteInt32(mist.ownerID)
	p.WriteInt32(mist.skillID)
	p.WriteByte(mist.skillLevel)
	p.WriteInt16(mist.box.x1)
	p.WriteInt16(mist.box.y1)
	p.WriteInt16(mist.box.x2)
	p.WriteInt16(mist.box.y2)
	p.WriteInt32(int32(mist.duration)) // Duration in seconds
	return p
}

func packetMistRemove(mistID int32, isPoisonMist bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelAffectedAreaRemove)
	p.WriteInt32(mistID)
	if isPoisonMist {
		p.WriteByte(1) // Fade out animation
	} else {
		p.WriteByte(0) // Instant removal
	}
	return p
}

// startPoisonMistTicker starts a goroutine that applies poison damage periodically to mobs in the mist
func (server *Server) startPoisonMistTicker(inst *fieldInstance, mist *fieldMist) {
	if !mist.isPoisonMist || inst == nil {
		return
	}

	// Poison ticks every 1 second
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		defer ticker.Stop()

		endTime := mist.createdAt.Add(time.Duration(mist.duration) * time.Second)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Check if mist has expired
				if time.Now().After(endTime) {
					return
				}

				// Apply poison damage to all mobs in the mist area
				inst.dispatch <- func() {
					// Get skill data for damage calculation
					skillData, err := nx.GetPlayerSkill(mist.skillID)
					if err != nil || mist.skillLevel == 0 || int(mist.skillLevel) > len(skillData) {
						return
					}

					// Poison damage is based on the X value from skill data
					poisonDamage := int32(skillData[mist.skillLevel-1].X)

					// Find all mobs in the mist area and apply poison damage
					for _, mob := range inst.lifePool.mobs {
						if mob.hp > 0 && mist.isInMist(mob.pos) {
							// Apply poison damage to mob
							inst.lifePool.mobDamaged(mob.spawnID, nil, poisonDamage)
						}
					}
				}
			}
		}
	}()
}
