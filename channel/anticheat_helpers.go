package channel

import (
	"fmt"
	"strings"
	"time"
)

// getPlayerIP extracts IP address from player connection
func getPlayerIP(player *Player) string {
	if player.Conn == nil {
		return "unknown"
	}
	addr := player.Conn.String()
	// Handle both IPv4 (IP:port) and IPv6 ([IP]:port) formats
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		// For IPv6, check if there's a bracket
		if strings.HasPrefix(addr, "[") {
			if bracketEnd := strings.Index(addr, "]"); bracketEnd != -1 {
				return addr[1:bracketEnd] // Extract IP from [IP]:port
			}
		}
		// For IPv4, extract before the last colon
		return addr[:idx]
	}
	return addr
}

// Helper functions for specific violation detection scenarios

// DetectExcessiveDamage checks if damage dealt exceeds expected bounds
func (vd *ViolationDetector) DetectExcessiveDamage(player *Player, damage int32, expectedMaxDamage int32) {
	if !vd.config.CombatDetection.Enabled {
		return
	}

	// Allow some margin for calculation differences (e.g., 150% of expected)
	threshold := float64(expectedMaxDamage) * 1.5
	if float64(damage) > threshold {
		vd.RecordViolation(ViolationEvent{
			AccountID:   player.accountID,
			CharacterID: player.ID,
			IPAddress:   getPlayerIP(player),
			ViolationType: ViolationExcessiveDamage,
			Category:    CategoryCombat,
			Severity:    SeverityHigh,
			DetectionDetails: fmt.Sprintf("Damage: %d, Expected max: %d (threshold: %.0f)",
				damage, expectedMaxDamage, threshold),
			MapID:     player.mapID,
			Timestamp: time.Now(),
		})
	}
}

// DetectAttackSpeedHack checks if attacks are happening too quickly
func (vd *ViolationDetector) DetectAttackSpeedHack(player *Player, timeSinceLastAttack time.Duration, minimumDelay time.Duration) {
	if !vd.config.CombatDetection.Enabled {
		return
	}

	// Allow small margin for network latency
	if timeSinceLastAttack < minimumDelay-50*time.Millisecond {
		vd.RecordViolation(ViolationEvent{
			AccountID:   player.accountID,
			CharacterID: player.ID,
			IPAddress:   getPlayerIP(player),
			ViolationType: ViolationAttackSpeed,
			Category:    CategoryCombat,
			Severity:    SeverityMedium,
			DetectionDetails: fmt.Sprintf("Attack interval: %v, Minimum allowed: %v",
				timeSinceLastAttack, minimumDelay),
			MapID:     player.mapID,
			Timestamp: time.Now(),
		})
	}
}

// DetectInvalidSkillUse checks if a player is using a skill they don't have or can't use
func (vd *ViolationDetector) DetectInvalidSkillUse(player *Player, skillID int32, reason string) {
	if !vd.config.CombatDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidSkill,
		Category:    CategoryCombat,
		Severity:    SeverityMedium,
		DetectionDetails: fmt.Sprintf("Skill ID: %d, Reason: %s", skillID, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectSpeedHack checks if player movement speed exceeds allowed limits
func (vd *ViolationDetector) DetectSpeedHack(player *Player, calculatedSpeed, maxAllowedSpeed float64) {
	if !vd.config.MovementDetection.Enabled {
		return
	}

	// Allow 110% margin for calculation differences
	if calculatedSpeed > maxAllowedSpeed*1.1 {
		vd.RecordViolation(ViolationEvent{
			AccountID:   player.accountID,
			CharacterID: player.ID,
			IPAddress:   getPlayerIP(player),
			ViolationType: ViolationSpeedHack,
			Category:    CategoryMovement,
			Severity:    SeverityMedium,
			DetectionDetails: fmt.Sprintf("Speed: %.2f, Max allowed: %.2f",
				calculatedSpeed, maxAllowedSpeed),
			MapID:     player.mapID,
			Timestamp: time.Now(),
		})
	}
}

// DetectTeleportHack checks for invalid teleportation
func (vd *ViolationDetector) DetectTeleportHack(player *Player, oldX, oldY, newX, newY int16, reason string) {
	if !vd.config.MovementDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationTeleportHack,
		Category:    CategoryMovement,
		Severity:    SeverityHigh,
		DetectionDetails: fmt.Sprintf("Moved from (%d,%d) to (%d,%d). Reason: %s",
			oldX, oldY, newX, newY, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectInvalidPosition checks for invalid player positions
func (vd *ViolationDetector) DetectInvalidPosition(player *Player, x, y int16, reason string) {
	if !vd.config.MovementDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidPosition,
		Category:    CategoryMovement,
		Severity:    SeverityMedium,
		DetectionDetails: fmt.Sprintf("Position: (%d,%d). Reason: %s", x, y, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectInvalidEquip checks for equipping items the player doesn't own or can't use
func (vd *ViolationDetector) DetectInvalidEquip(player *Player, itemID int32, reason string) {
	if !vd.config.InventoryDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidEquip,
		Category:    CategoryInventory,
		Severity:    SeverityHigh,
		DetectionDetails: fmt.Sprintf("Item ID: %d, Reason: %s", itemID, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectInvalidItemUse checks for using items the player doesn't have
func (vd *ViolationDetector) DetectInvalidItemUse(player *Player, itemID int32, reason string) {
	if !vd.config.InventoryDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidItemUse,
		Category:    CategoryInventory,
		Severity:    SeverityMedium,
		DetectionDetails: fmt.Sprintf("Item ID: %d, Reason: %s", itemID, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectInvalidTrade checks for invalid trade operations
func (vd *ViolationDetector) DetectInvalidTrade(player *Player, reason string) {
	if !vd.config.EconomyDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidTrade,
		Category:    CategoryEconomy,
		Severity:    SeverityHigh,
		DetectionDetails: reason,
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectDuplication checks for potential item duplication
func (vd *ViolationDetector) DetectDuplication(player *Player, itemID int32, reason string) {
	if !vd.config.EconomyDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationDuplication,
		Category:    CategoryEconomy,
		Severity:    SeverityCritical,
		DetectionDetails: fmt.Sprintf("Item ID: %d, Reason: %s", itemID, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectOverflow checks for overflow/underflow exploits
func (vd *ViolationDetector) DetectOverflow(player *Player, reason string) {
	if !vd.config.EconomyDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationOverflow,
		Category:    CategoryEconomy,
		Severity:    SeverityCritical,
		DetectionDetails: reason,
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectCooldownBypass checks for cooldown bypasses
func (vd *ViolationDetector) DetectCooldownBypass(player *Player, skillID int32, cooldownRemaining time.Duration) {
	if !vd.config.SkillDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationCooldownBypass,
		Category:    CategorySkill,
		Severity:    SeverityMedium,
		DetectionDetails: fmt.Sprintf("Skill ID: %d, Cooldown remaining: %v",
			skillID, cooldownRemaining),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectUnlearnedSkill checks for using skills the player hasn't learned
func (vd *ViolationDetector) DetectUnlearnedSkill(player *Player, skillID int32) {
	if !vd.config.SkillDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationUnlearnedSkill,
		Category:    CategorySkill,
		Severity:    SeverityHigh,
		DetectionDetails: fmt.Sprintf("Skill ID: %d", skillID),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectInvalidPacketSequence checks for invalid packet sequences
func (vd *ViolationDetector) DetectInvalidPacketSequence(player *Player, reason string) {
	if !vd.config.PacketDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationInvalidSequence,
		Category:    CategoryPacket,
		Severity:    SeverityMedium,
		DetectionDetails: reason,
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}

// DetectMalformedPacket checks for malformed packets
func (vd *ViolationDetector) DetectMalformedPacket(player *Player, packetType string, reason string) {
	if !vd.config.PacketDetection.Enabled {
		return
	}

	vd.RecordViolation(ViolationEvent{
		AccountID:   player.accountID,
		CharacterID: player.ID,
		IPAddress:   getPlayerIP(player),
		ViolationType: ViolationMalformedPacket,
		Category:    CategoryPacket,
		Severity:    SeverityHigh,
		DetectionDetails: fmt.Sprintf("Packet type: %s, Reason: %s", packetType, reason),
		MapID:     player.mapID,
		Timestamp: time.Now(),
	})
}
