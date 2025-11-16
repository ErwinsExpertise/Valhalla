package channel

import (
	"log"
	"math"
	"time"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/mpacket"
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
		log.Println("Mist: Failed to generate mist ID")
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

	log.Printf("Mist: Created mist ID=%d, owner=%d, skill=%d, level=%d, pos=(%d,%d), box=(%d,%d,%d,%d), duration=%d",
		mistID, ownerID, skillID, skillLevel, pos.x, pos.y,
		mist.box.x1, mist.box.y1, mist.box.x2, mist.box.y2, duration)

	// Send spawn packet to all players
	pool.instance.send(packetMistSpawn(mist))
	log.Printf("Mist: Sent spawn packet to all players")

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
	p.WriteBool(false) // MobMist - false for player-created poison mist, true for mob mists
	p.WriteInt32(mist.skillID)
	p.WriteByte(mist.skillLevel)
	p.WriteInt16(0) // delay
	p.WriteInt32(int32(mist.box.x1))
	p.WriteInt32(int32(mist.box.y1))
	p.WriteInt32(int32(mist.box.x2))
	p.WriteInt32(int32(mist.box.y2))
	
	log.Printf("Mist: Packet created - ID=%d, MobMist=false, Skill=%d, Level=%d, Box=(%d,%d,%d,%d)",
		mist.ID, mist.skillID, mist.skillLevel, mist.box.x1, mist.box.y1, mist.box.x2, mist.box.y2)
	
	return p
}

func packetMistRemove(mistID int32, isPoisonMist bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelAffectedAreaRemove)
	p.WriteInt32(mistID)
	log.Printf("Mist: Removal packet created for ID=%d", mistID)
	return p
}

// startPoisonMistTicker applies poison buff to mobs entering the mist area
func (server *Server) startPoisonMistTicker(inst *fieldInstance, mist *fieldMist) {
	if !mist.isPoisonMist || inst == nil {
		return
	}

	// Check every second for mobs in the mist
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		defer ticker.Stop()

		endTime := mist.createdAt.Add(time.Duration(mist.duration) * time.Second)

		for range ticker.C {
			// Check if mist has expired
			if time.Now().After(endTime) {
				return
			}

			// Check if mist still exists
			if _, exists := inst.mistPool.mists[mist.ID]; !exists {
				return
			}

			// Apply poison buff to all mobs in the mist area
			inst.dispatch <- func() {
				// Find all mobs in the mist area and apply poison buff
				for spawnID, mob := range inst.lifePool.mobs {
					if mob != nil && mob.hp > 0 && mist.isInMist(mob.pos) {
						// Check if mob is already poisoned - don't reapply
						if (mob.statBuff & skill.MobStat.Poison) != 0 {
							continue
						}
						// Apply poison mob buff
						inst.lifePool.applyMobBuff(spawnID, mist.skillID, mist.skillLevel, skill.MobStat.Poison, inst)
						log.Printf("PoisonMist: Applied poison buff to mob %d", spawnID)
					}
				}
			}
		}
	}()
}
