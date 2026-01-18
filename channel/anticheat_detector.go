package channel

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Hucaru/Valhalla/common"
)

// ViolationType represents a specific type of violation
type ViolationType string

const (
	// Combat violations
	ViolationExcessiveDamage ViolationType = "excessive_damage"
	ViolationAttackSpeed     ViolationType = "attack_speed_hack"
	ViolationInvalidSkill    ViolationType = "invalid_skill_use"

	// Movement violations
	ViolationSpeedHack       ViolationType = "speed_hack"
	ViolationTeleportHack    ViolationType = "teleport_hack"
	ViolationInvalidPosition ViolationType = "invalid_position"

	// Inventory violations
	ViolationInvalidEquip   ViolationType = "invalid_equip"
	ViolationInvalidItemUse ViolationType = "invalid_item_use"

	// Economy violations
	ViolationInvalidTrade ViolationType = "invalid_trade"
	ViolationDuplication  ViolationType = "item_duplication"
	ViolationOverflow     ViolationType = "overflow_exploit"

	// Skill violations
	ViolationCooldownBypass ViolationType = "cooldown_bypass"
	ViolationUnlearnedSkill ViolationType = "unlearned_skill"

	// Packet violations
	ViolationInvalidSequence ViolationType = "invalid_packet_sequence"
	ViolationMalformedPacket ViolationType = "malformed_packet"
)

// ViolationCategory represents broad categories of violations
type ViolationCategory string

const (
	CategoryCombat    ViolationCategory = "combat"
	CategoryMovement  ViolationCategory = "movement"
	CategoryInventory ViolationCategory = "inventory"
	CategoryEconomy   ViolationCategory = "economy"
	CategorySkill     ViolationCategory = "skill"
	CategoryPacket    ViolationCategory = "packet"
)

// ViolationSeverity represents how severe a violation is
type ViolationSeverity string

const (
	SeverityLow      ViolationSeverity = "low"
	SeverityMedium   ViolationSeverity = "medium"
	SeverityHigh     ViolationSeverity = "high"
	SeverityCritical ViolationSeverity = "critical"
)

// ViolationEvent represents a single violation event
type ViolationEvent struct {
	AccountID        int32
	CharacterID      int32
	IPAddress        string
	ViolationType    ViolationType
	Category         ViolationCategory
	Severity         ViolationSeverity
	DetectionDetails string
	MapID            int32
	Timestamp        time.Time
}

// ViolationCounter tracks violations in a rolling window
type ViolationCounter struct {
	AccountID     int32
	CharacterID   int32
	ViolationType ViolationType
	Count         int
	WindowStart   time.Time
	LastViolation time.Time
}

// ViolationDetector handles violation detection and tracking
type ViolationDetector struct {
	config      AntiCheatConfig
	banService  *BanService
	mu          sync.RWMutex
	counters    map[string]*ViolationCounter // key: "accountID:characterID:violationType"
	cleanupTick *time.Ticker
}

// NewViolationDetector creates a new violation detector
func NewViolationDetector(config AntiCheatConfig, banService *BanService) *ViolationDetector {
	vd := &ViolationDetector{
		config:     config,
		banService: banService,
		counters:   make(map[string]*ViolationCounter),
	}

	// Start cleanup goroutine
	vd.cleanupTick = time.NewTicker(config.ViolationDetection.CleanupInterval)
	go vd.cleanupExpiredCounters()

	return vd
}

// RecordViolation records a violation and checks if action should be taken
func (vd *ViolationDetector) RecordViolation(event ViolationEvent) error {
	if !vd.config.Enabled {
		return nil
	}

	// Log the violation to database
	err := vd.logViolation(event)
	if err != nil {
		log.Printf("Error logging violation: %v", err)
	}

	// Update rolling counter
	counter := vd.updateCounter(event)

	// Check if threshold exceeded
	threshold, window, banType := vd.getThresholdConfig(event.ViolationType)
	if counter.Count >= threshold {
		// Check if violation occurred within the window
		if time.Since(counter.WindowStart) <= window {
			log.Printf("Violation threshold exceeded for account %d character %d: %s (count: %d, threshold: %d)",
				event.AccountID, event.CharacterID, event.ViolationType, counter.Count, threshold)

			// Take action
			err = vd.takeAction(event, counter, banType)
			if err != nil {
				log.Printf("Error taking action for violation: %v", err)
				return err
			}

			// Reset counter after action
			vd.resetCounter(event.AccountID, event.CharacterID, event.ViolationType)
		}
	}

	return nil
}

// logViolation logs a violation to the database
func (vd *ViolationDetector) logViolation(event ViolationEvent) error {
	query := `INSERT INTO violation_logs 
			(accountID, characterID, ipAddress, violationType, violationCategory, severity, 
			detectionDetails, mapID, timestamp) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := common.DB.Exec(query, event.AccountID, event.CharacterID, event.IPAddress,
		event.ViolationType, event.Category, event.Severity, event.DetectionDetails,
		event.MapID, event.Timestamp)

	if err != nil {
		return fmt.Errorf("failed to log violation: %w", err)
	}

	return nil
}

// updateCounter updates the rolling window counter for a violation
func (vd *ViolationDetector) updateCounter(event ViolationEvent) *ViolationCounter {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	key := fmt.Sprintf("%d:%d:%s", event.AccountID, event.CharacterID, event.ViolationType)
	counter, exists := vd.counters[key]

	now := time.Now()
	windowDuration := vd.config.ViolationDetection.RollingWindowDuration

	if !exists {
		// Create new counter
		counter = &ViolationCounter{
			AccountID:     event.AccountID,
			CharacterID:   event.CharacterID,
			ViolationType: event.ViolationType,
			Count:         1,
			WindowStart:   now,
			LastViolation: now,
		}
		vd.counters[key] = counter
	} else {
		// Check if we're still within the rolling window
		if now.Sub(counter.WindowStart) > windowDuration {
			// Window expired, reset counter
			counter.Count = 1
			counter.WindowStart = now
		} else {
			// Within window, increment
			counter.Count++
		}
		counter.LastViolation = now
	}

	// Also update database counter
	vd.updateDatabaseCounter(counter)

	return counter
}

// updateDatabaseCounter persists counter to database
func (vd *ViolationDetector) updateDatabaseCounter(counter *ViolationCounter) {
	query := `INSERT INTO violation_counters 
			(accountID, characterID, violationType, count, windowStart, lastViolation) 
			VALUES (?, ?, ?, ?, ?, ?) 
			ON DUPLICATE KEY UPDATE 
			count = VALUES(count), 
			windowStart = VALUES(windowStart), 
			lastViolation = VALUES(lastViolation)`

	_, err := common.DB.Exec(query, counter.AccountID, counter.CharacterID, counter.ViolationType,
		counter.Count, counter.WindowStart, counter.LastViolation)
	if err != nil {
		log.Printf("Error updating database counter: %v", err)
	}
}

// resetCounter resets a violation counter
func (vd *ViolationDetector) resetCounter(accountID, characterID int32, violationType ViolationType) {
	vd.mu.Lock()
	defer vd.mu.Unlock()

	key := fmt.Sprintf("%d:%d:%s", accountID, characterID, violationType)
	delete(vd.counters, key)

	// Also delete from database
	query := `DELETE FROM violation_counters 
			WHERE accountID = ? AND characterID = ? AND violationType = ?`
	_, err := common.DB.Exec(query, accountID, characterID, violationType)
	if err != nil {
		log.Printf("Error deleting database counter: %v", err)
	}
}

// getThresholdConfig returns the threshold and window for a violation type
func (vd *ViolationDetector) getThresholdConfig(violationType ViolationType) (threshold int, window time.Duration, banType string) {
	switch violationType {
	// Combat violations
	case ViolationExcessiveDamage:
		return vd.config.CombatDetection.ExcessiveDamageThreshold,
			vd.config.CombatDetection.ExcessiveDamageWindow,
			vd.config.CombatDetection.ExcessiveDamageBanType
	case ViolationAttackSpeed:
		return vd.config.CombatDetection.AttackSpeedThreshold,
			vd.config.CombatDetection.AttackSpeedWindow,
			vd.config.CombatDetection.AttackSpeedBanType
	case ViolationInvalidSkill:
		return vd.config.CombatDetection.InvalidSkillThreshold,
			vd.config.CombatDetection.InvalidSkillWindow,
			vd.config.CombatDetection.InvalidSkillBanType

	// Movement violations
	case ViolationSpeedHack:
		return vd.config.MovementDetection.SpeedHackThreshold,
			vd.config.MovementDetection.SpeedHackWindow,
			vd.config.MovementDetection.SpeedHackBanType
	case ViolationTeleportHack:
		return vd.config.MovementDetection.TeleportHackThreshold,
			vd.config.MovementDetection.TeleportHackWindow,
			vd.config.MovementDetection.TeleportHackBanType
	case ViolationInvalidPosition:
		return vd.config.MovementDetection.InvalidPositionThreshold,
			vd.config.MovementDetection.InvalidPositionWindow,
			vd.config.MovementDetection.InvalidPositionBanType

	// Inventory violations
	case ViolationInvalidEquip:
		return vd.config.InventoryDetection.InvalidEquipThreshold,
			vd.config.InventoryDetection.InvalidEquipWindow,
			vd.config.InventoryDetection.InvalidEquipBanType
	case ViolationInvalidItemUse:
		return vd.config.InventoryDetection.InvalidItemUseThreshold,
			vd.config.InventoryDetection.InvalidItemUseWindow,
			vd.config.InventoryDetection.InvalidItemUseBanType

	// Economy violations
	case ViolationInvalidTrade:
		return vd.config.EconomyDetection.InvalidTradeThreshold,
			vd.config.EconomyDetection.InvalidTradeWindow,
			vd.config.EconomyDetection.InvalidTradeBanType
	case ViolationDuplication:
		return vd.config.EconomyDetection.DuplicationThreshold,
			vd.config.EconomyDetection.DuplicationWindow,
			vd.config.EconomyDetection.DuplicationBanType
	case ViolationOverflow:
		return vd.config.EconomyDetection.OverflowThreshold,
			vd.config.EconomyDetection.OverflowWindow,
			vd.config.EconomyDetection.OverflowBanType

	// Skill violations
	case ViolationCooldownBypass:
		return vd.config.SkillDetection.CooldownBypassThreshold,
			vd.config.SkillDetection.CooldownBypassWindow,
			vd.config.SkillDetection.CooldownBypassBanType
	case ViolationUnlearnedSkill:
		return vd.config.SkillDetection.UnlearnedSkillThreshold,
			vd.config.SkillDetection.UnlearnedSkillWindow,
			vd.config.SkillDetection.UnlearnedSkillBanType

	// Packet violations
	case ViolationInvalidSequence:
		return vd.config.PacketDetection.InvalidSequenceThreshold,
			vd.config.PacketDetection.InvalidSequenceWindow,
			vd.config.PacketDetection.InvalidSequenceBanType
	case ViolationMalformedPacket:
		return vd.config.PacketDetection.MalformedPacketThreshold,
			vd.config.PacketDetection.MalformedPacketWindow,
			vd.config.PacketDetection.MalformedPacketBanType

	default:
		// Default to medium threshold
		return 5, 5 * time.Minute, "temporary"
	}
}

// takeAction takes appropriate action when a violation threshold is exceeded
func (vd *ViolationDetector) takeAction(event ViolationEvent, counter *ViolationCounter, banTypeStr string) error {
	var banType BanType
	if banTypeStr == "permanent" {
		banType = BanTypePermanent
	} else {
		banType = BanTypeTemporary
	}

	reason := fmt.Sprintf("Anti-cheat: %s violation detected (%d occurrences within window)",
		event.ViolationType, counter.Count)

	// Apply IP ban based on config
	var ipAddress *string
	if vd.config.IPBanMode == "always" ||
		(vd.config.IPBanMode == "permanent_only" && banType == BanTypePermanent) {
		ipAddress = &event.IPAddress
	}

	// Issue the ban
	err := vd.banService.IssueBan(&event.AccountID, &event.CharacterID, ipAddress,
		banType, BanTargetAccount, reason, "ANTICHEAT", false)
	if err != nil {
		return fmt.Errorf("failed to issue ban: %w", err)
	}

	// Update violation log with action taken
	updateQuery := `UPDATE violation_logs 
				SET actionTaken = ? 
				WHERE accountID = ? AND characterID = ? AND violationType = ? 
				ORDER BY timestamp DESC LIMIT ?`
	_, err = common.DB.Exec(updateQuery, fmt.Sprintf("Ban issued: %s", banType),
		event.AccountID, event.CharacterID, event.ViolationType, counter.Count)
	if err != nil {
		log.Printf("Error updating violation logs with action: %v", err)
	}

	return nil
}

// cleanupExpiredCounters periodically removes expired counters
func (vd *ViolationDetector) cleanupExpiredCounters() {
	for range vd.cleanupTick.C {
		vd.mu.Lock()
		now := time.Now()
		windowDuration := vd.config.ViolationDetection.RollingWindowDuration

		for key, counter := range vd.counters {
			if now.Sub(counter.LastViolation) > windowDuration*2 {
				delete(vd.counters, key)
			}
		}
		vd.mu.Unlock()

		// Also cleanup database
		query := `DELETE FROM violation_counters 
				WHERE lastViolation < DATE_SUB(NOW(), INTERVAL ? SECOND)`
		_, err := common.DB.Exec(query, int(windowDuration.Seconds()*2))
		if err != nil {
			log.Printf("Error cleaning up database counters: %v", err)
		}
	}
}

// Stop stops the violation detector
func (vd *ViolationDetector) Stop() {
	if vd.cleanupTick != nil {
		vd.cleanupTick.Stop()
	}
}

// GetViolationHistory retrieves violation history for a character
func (vd *ViolationDetector) GetViolationHistory(characterID int32, limit int) ([]ViolationEvent, error) {
	query := `SELECT accountID, characterID, ipAddress, violationType, violationCategory, 
			severity, detectionDetails, mapID, timestamp 
			FROM violation_logs 
			WHERE characterID = ? 
			ORDER BY timestamp DESC 
			LIMIT ?`

	rows, err := common.DB.Query(query, characterID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get violation history: %w", err)
	}
	defer rows.Close()

	var events []ViolationEvent
	for rows.Next() {
		var event ViolationEvent
		var ipAddress sql.NullString

		err := rows.Scan(&event.AccountID, &event.CharacterID, &ipAddress,
			&event.ViolationType, &event.Category, &event.Severity,
			&event.DetectionDetails, &event.MapID, &event.Timestamp)
		if err != nil {
			log.Printf("Error scanning violation event: %v", err)
			continue
		}

		if ipAddress.Valid {
			event.IPAddress = ipAddress.String
		}

		events = append(events, event)
	}

	return events, nil
}
