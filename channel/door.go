package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/nx"
)

// createMysticDoor creates a mystic door for the player
// Following reference: https://github.com/sewil/OpenMG/blob/main/WvsBeta.Game/GameObjects/Door.cs#L71
func createMysticDoor(plr *Player, skillID int32, skillLevel byte) {
	// Remove existing door if player already has one
	if plr.doorMapID != 0 {
		removeMysticDoor(plr)
	}

	// Use server-side player position
	doorPos := plr.pos

	// Get skill duration from skill data
	var duration int64 = 60 // default 60 seconds
	if data, err := nx.GetPlayerSkill(skillID); err == nil {
		idx := int(skillLevel) - 1
		if idx >= 0 && idx < len(data) && data[idx].Time > 0 {
			duration = data[idx].Time
		}
	}

	// Show skill animation
	plr.inst.send(packetPlayerSkillAnimation(plr.ID, false, skillID, skillLevel))

	// Create door in source map
	createSourceDoor(plr, doorPos)

	// Create town portal in return map
	returnMapID := plr.inst.returnMapID
	if returnMapID > 0 {
		if returnField, ok := plr.inst.server.fields[returnMapID]; ok {
			if returnInst, err := returnField.getInstance(0); err == nil {
				createTownDoor(plr, returnInst, doorPos)
			}
		}
	}

	// Schedule door expiration
	go func(playerID int32, sourceMapID, townMapID int32, dur int64, server *Server) {
		time.Sleep(time.Duration(dur) * time.Second)

		// Dispatch to field goroutine for thread safety
		if field, ok := server.fields[sourceMapID]; ok {
			if inst, err := field.getInstance(0); err == nil {
				inst.dispatch <- func() {
					mysticDoorExpired(playerID, sourceMapID, townMapID, server)
				}
			}
		}
	}(plr.ID, plr.mapID, returnMapID, duration, plr.inst.server)
}

// removeMysticDoor removes a player's existing mystic door
func removeMysticDoor(plr *Player) {
	// Remove source door
	if plr.doorMapID != 0 {
		if doorField, ok := plr.inst.server.fields[plr.doorMapID]; ok {
			if doorInst, err := doorField.getInstance(plr.inst.id); err == nil {
				// Remove door visual
				doorInst.send(packetMapRemoveMysticDoor(plr.doorSpawnID, true))
				// Remove portal
				doorInst.removePortalAtIndex(plr.doorPortalIndex)
				// Remove from doors map
				delete(doorInst.mysticDoors, plr.ID)
			}
		}
	}

	// Remove town door
	if plr.townDoorMapID != 0 {
		if townField, ok := plr.inst.server.fields[plr.townDoorMapID]; ok {
			if townInst, err := townField.getInstance(0); err == nil {
				// Remove town door visual
				townInst.send(packetMapRemoveMysticDoor(plr.townDoorSpawnID, true))
				// Remove from doors map
				delete(townInst.mysticDoors, plr.ID)
			}
		}
	}

	// Send portal removal packet
	portalRemovePacket := packetMapRemovePortal()
	if plr.party != nil {
		// Send to all party members
		for _, member := range plr.party.players {
			if member != nil {
				member.Send(portalRemovePacket)
			}
		}
	} else {
		// Send to owner only
		plr.Send(portalRemovePacket)
	}

	// Clear player door state
	plr.doorMapID = 0
	plr.doorSpawnID = 0
	plr.doorPortalIndex = 0
	plr.townDoorMapID = 0
	plr.townDoorSpawnID = 0
	plr.townPortalIndex = 0
}

// createSourceDoor creates the door in the source map
func createSourceDoor(plr *Player, doorPos pos) {
	// Generate spawn ID
	doorSpawnID := plr.inst.idCounter
	plr.inst.idCounter++

	// Store on player
	plr.doorMapID = plr.mapID
	plr.doorSpawnID = doorSpawnID

	// Create portal object
	sourcePortal := portal{
		pos:         doorPos,
		name:        "tp",
		destFieldID: plr.inst.returnMapID,
		destName:    "sp",
		temporary:   true,
	}
	plr.doorPortalIndex = plr.inst.addPortal(sourcePortal)

	// Register in mystic doors map
	plr.inst.mysticDoors[plr.ID] = &mysticDoorInfo{
		ownerID:     plr.ID,
		spawnID:     doorSpawnID,
		portalIndex: plr.doorPortalIndex,
		pos:         doorPos,
		destMapID:   plr.inst.returnMapID,
	}

	// Send visual packet to entire instance
	plr.inst.send(packetMapSpawnMysticDoor(doorSpawnID, doorPos, true))

	// Send portal enable packet
	// If in party: Send to all party members in this map
	// If solo: Send to owner only
	if plr.party != nil {
		for _, member := range plr.party.players {
			if member != nil && member.mapID == plr.mapID {
				member.Send(packetMapPortal(plr.mapID, plr.inst.returnMapID, doorPos))
			}
		}
	} else {
		plr.Send(packetMapPortal(plr.mapID, plr.inst.returnMapID, doorPos))
	}

	log.Printf("[Mystic Door] Created source door: playerID=%d, mapID=%d, pos=(%d,%d), party=%v",
		plr.ID, plr.mapID, doorPos.x, doorPos.y, plr.party != nil)
}

// createTownDoor creates the door in the town map
func createTownDoor(plr *Player, townInst *fieldInstance, doorPos pos) {
	// Find available town portal
	townPortalIdx, townPortal, err := townInst.findAvailableTownPortal()
	if err != nil {
		log.Printf("[Mystic Door] ERROR: No available town portals for player %d", plr.ID)
		return
	}

	// Generate spawn ID
	townDoorSpawnID := townInst.idCounter
	townInst.idCounter++

	// Store on player
	plr.townDoorMapID = townInst.fieldID
	plr.townDoorSpawnID = townDoorSpawnID
	plr.townPortalIndex = townPortalIdx

	// Modify existing portal to point to source map
	townInst.portals[townPortalIdx].destFieldID = plr.mapID
	townInst.portals[townPortalIdx].destName = "sp"

	// Register in mystic doors map
	townInst.mysticDoors[plr.ID] = &mysticDoorInfo{
		ownerID:     plr.ID,
		spawnID:     townDoorSpawnID,
		portalIndex: townPortalIdx,
		pos:         townPortal.pos,
		destMapID:   plr.mapID,
		townPortal:  true,
	}

	// Send town door visual to entire instance
	townInst.send(packetMapSpawnMysticDoor(townInst.mysticDoors[plr.ID].spawnID, townPortal.pos, false))

	// Send portal enable packet
	// If in party: Send to all party members in town
	// If solo: Send to owner only (even if not in town yet)
	if plr.party != nil {
		for _, member := range plr.party.players {
			if member != nil && member.mapID == townInst.fieldID {
				member.Send(packetMapPortal(townInst.fieldID, plr.mapID, townPortal.pos))
			}
		}
	} else {
		plr.Send(packetMapPortal(townInst.fieldID, plr.mapID, townPortal.pos))
	}

	log.Printf("[Mystic Door] Created town door: playerID=%d, townMap=%d, sourceMap=%d, portalIdx=%d, pos=(%d,%d), party=%v",
		plr.ID, townInst.fieldID, plr.mapID, townPortalIdx, townPortal.pos.x, townPortal.pos.y, plr.party != nil)
}

// mysticDoorExpired handles door expiration
func mysticDoorExpired(playerID, sourceMapID, townMapID int32, server *Server) {
	log.Printf("[Mystic Door] Door expired for player %d", playerID)

	// Try to find the player
	var plr *Player

	// Search in source map
	if sourceField, ok := server.fields[sourceMapID]; ok {
		if sourceInst, err := sourceField.getInstance(0); err == nil {
			for _, p := range sourceInst.players {
				if p.ID == playerID {
					plr = p
					break
				}
			}
		}
	}

	// If not found, search in town map
	if plr == nil && townMapID > 0 {
		if townField, ok := server.fields[townMapID]; ok {
			if townInst, err := townField.getInstance(0); err == nil {
				for _, p := range townInst.players {
					if p.ID == playerID {
						plr = p
						break
					}
				}
			}
		}
	}

	// Remove the door
	if plr != nil {
		removeMysticDoor(plr)
	} else {
		// Player offline - still need to cleanup doors from instances
		log.Printf("[Mystic Door] Player %d offline, cleaning up doors from instances", playerID)

		// Clean source map
		if sourceField, ok := server.fields[sourceMapID]; ok {
			if sourceInst, err := sourceField.getInstance(0); err == nil {
				if doorInfo, exists := sourceInst.mysticDoors[playerID]; exists {
					sourceInst.send(packetMapRemoveMysticDoor(doorInfo.spawnID, true))
					sourceInst.removePortalAtIndex(doorInfo.portalIndex)
					delete(sourceInst.mysticDoors, playerID)
				}
			}
		}

		// Clean town map
		if townMapID > 0 {
			if townField, ok := server.fields[townMapID]; ok {
				if townInst, err := townField.getInstance(0); err == nil {
					if doorInfo, exists := townInst.mysticDoors[playerID]; exists {
						townInst.send(packetMapRemoveMysticDoor(doorInfo.spawnID, true))
						delete(townInst.mysticDoors, playerID)
					}
				}
			}
		}
	}
}
