package anticheat

import "time"

// Config contains all anti-cheat and ban system configuration
type Config struct {
	Enabled                    bool                       `toml:"enabled" mapstructure:"enabled"`
	DefaultTempBanDuration     time.Duration              `toml:"defaultTempBanDuration" mapstructure:"defaultTempBanDuration"`
	TempBansBeforePermanent    int                        `toml:"tempBansBeforePermanent" mapstructure:"tempBansBeforePermanent"`
	IPBanMode                  string                     `toml:"ipBanMode" mapstructure:"ipBanMode"` // "never", "permanent_only", "always"
	GMBansIncrementCounter     bool                       `toml:"gmBansIncrementCounter" mapstructure:"gmBansIncrementCounter"`
	ViolationDetection         ViolationDetectionConfig   `toml:"violationDetection" mapstructure:"violationDetection"`
	CombatDetection            CombatDetectionConfig      `toml:"combatDetection" mapstructure:"combatDetection"`
	MovementDetection          MovementDetectionConfig    `toml:"movementDetection" mapstructure:"movementDetection"`
	InventoryDetection         InventoryDetectionConfig   `toml:"inventoryDetection" mapstructure:"inventoryDetection"`
	EconomyDetection           EconomyDetectionConfig     `toml:"economyDetection" mapstructure:"economyDetection"`
	SkillDetection             SkillDetectionConfig       `toml:"skillDetection" mapstructure:"skillDetection"`
	PacketDetection            PacketDetectionConfig      `toml:"packetDetection" mapstructure:"packetDetection"`
}

// ViolationDetectionConfig contains general violation detection settings
type ViolationDetectionConfig struct {
	RollingWindowDuration time.Duration `toml:"rollingWindowDuration" mapstructure:"rollingWindowDuration"`
	CleanupInterval       time.Duration `toml:"cleanupInterval" mapstructure:"cleanupInterval"`
}

// CombatDetectionConfig contains combat-related detection settings
type CombatDetectionConfig struct {
	Enabled                  bool          `toml:"enabled" mapstructure:"enabled"`
	ExcessiveDamageThreshold int           `toml:"excessiveDamageThreshold" mapstructure:"excessiveDamageThreshold"`
	ExcessiveDamageWindow    time.Duration `toml:"excessiveDamageWindow" mapstructure:"excessiveDamageWindow"`
	ExcessiveDamageBanType   string        `toml:"excessiveDamageBanType" mapstructure:"excessiveDamageBanType"` // "temporary", "permanent"
	
	AttackSpeedThreshold     int           `toml:"attackSpeedThreshold" mapstructure:"attackSpeedThreshold"`
	AttackSpeedWindow        time.Duration `toml:"attackSpeedWindow" mapstructure:"attackSpeedWindow"`
	AttackSpeedBanType       string        `toml:"attackSpeedBanType" mapstructure:"attackSpeedBanType"`
	
	InvalidSkillThreshold    int           `toml:"invalidSkillThreshold" mapstructure:"invalidSkillThreshold"`
	InvalidSkillWindow       time.Duration `toml:"invalidSkillWindow" mapstructure:"invalidSkillWindow"`
	InvalidSkillBanType      string        `toml:"invalidSkillBanType" mapstructure:"invalidSkillBanType"`
}

// MovementDetectionConfig contains movement-related detection settings
type MovementDetectionConfig struct {
	Enabled                    bool          `toml:"enabled" mapstructure:"enabled"`
	SpeedHackThreshold         int           `toml:"speedHackThreshold" mapstructure:"speedHackThreshold"`
	SpeedHackWindow            time.Duration `toml:"speedHackWindow" mapstructure:"speedHackWindow"`
	SpeedHackBanType           string        `toml:"speedHackBanType" mapstructure:"speedHackBanType"`
	
	TeleportHackThreshold      int           `toml:"teleportHackThreshold" mapstructure:"teleportHackThreshold"`
	TeleportHackWindow         time.Duration `toml:"teleportHackWindow" mapstructure:"teleportHackWindow"`
	TeleportHackBanType        string        `toml:"teleportHackBanType" mapstructure:"teleportHackBanType"`
	
	InvalidPositionThreshold   int           `toml:"invalidPositionThreshold" mapstructure:"invalidPositionThreshold"`
	InvalidPositionWindow      time.Duration `toml:"invalidPositionWindow" mapstructure:"invalidPositionWindow"`
	InvalidPositionBanType     string        `toml:"invalidPositionBanType" mapstructure:"invalidPositionBanType"`
}

// InventoryDetectionConfig contains inventory-related detection settings
type InventoryDetectionConfig struct {
	Enabled                     bool          `toml:"enabled" mapstructure:"enabled"`
	InvalidEquipThreshold       int           `toml:"invalidEquipThreshold" mapstructure:"invalidEquipThreshold"`
	InvalidEquipWindow          time.Duration `toml:"invalidEquipWindow" mapstructure:"invalidEquipWindow"`
	InvalidEquipBanType         string        `toml:"invalidEquipBanType" mapstructure:"invalidEquipBanType"`
	
	InvalidItemUseThreshold     int           `toml:"invalidItemUseThreshold" mapstructure:"invalidItemUseThreshold"`
	InvalidItemUseWindow        time.Duration `toml:"invalidItemUseWindow" mapstructure:"invalidItemUseWindow"`
	InvalidItemUseBanType       string        `toml:"invalidItemUseBanType" mapstructure:"invalidItemUseBanType"`
}

// EconomyDetectionConfig contains economy-related detection settings
type EconomyDetectionConfig struct {
	Enabled                   bool          `toml:"enabled" mapstructure:"enabled"`
	InvalidTradeThreshold     int           `toml:"invalidTradeThreshold" mapstructure:"invalidTradeThreshold"`
	InvalidTradeWindow        time.Duration `toml:"invalidTradeWindow" mapstructure:"invalidTradeWindow"`
	InvalidTradeBanType       string        `toml:"invalidTradeBanType" mapstructure:"invalidTradeBanType"`
	
	DuplicationThreshold      int           `toml:"duplicationThreshold" mapstructure:"duplicationThreshold"`
	DuplicationWindow         time.Duration `toml:"duplicationWindow" mapstructure:"duplicationWindow"`
	DuplicationBanType        string        `toml:"duplicationBanType" mapstructure:"duplicationBanType"`
	
	OverflowThreshold         int           `toml:"overflowThreshold" mapstructure:"overflowThreshold"`
	OverflowWindow            time.Duration `toml:"overflowWindow" mapstructure:"overflowWindow"`
	OverflowBanType           string        `toml:"overflowBanType" mapstructure:"overflowBanType"`
}

// SkillDetectionConfig contains skill-related detection settings
type SkillDetectionConfig struct {
	Enabled                   bool          `toml:"enabled" mapstructure:"enabled"`
	CooldownBypassThreshold   int           `toml:"cooldownBypassThreshold" mapstructure:"cooldownBypassThreshold"`
	CooldownBypassWindow      time.Duration `toml:"cooldownBypassWindow" mapstructure:"cooldownBypassWindow"`
	CooldownBypassBanType     string        `toml:"cooldownBypassBanType" mapstructure:"cooldownBypassBanType"`
	
	UnlearnedSkillThreshold   int           `toml:"unlearnedSkillThreshold" mapstructure:"unlearnedSkillThreshold"`
	UnlearnedSkillWindow      time.Duration `toml:"unlearnedSkillWindow" mapstructure:"unlearnedSkillWindow"`
	UnlearnedSkillBanType     string        `toml:"unlearnedSkillBanType" mapstructure:"unlearnedSkillBanType"`
}

// PacketDetectionConfig contains packet-related detection settings
type PacketDetectionConfig struct {
	Enabled                    bool          `toml:"enabled" mapstructure:"enabled"`
	InvalidSequenceThreshold   int           `toml:"invalidSequenceThreshold" mapstructure:"invalidSequenceThreshold"`
	InvalidSequenceWindow      time.Duration `toml:"invalidSequenceWindow" mapstructure:"invalidSequenceWindow"`
	InvalidSequenceBanType     string        `toml:"invalidSequenceBanType" mapstructure:"invalidSequenceBanType"`
	
	MalformedPacketThreshold   int           `toml:"malformedPacketThreshold" mapstructure:"malformedPacketThreshold"`
	MalformedPacketWindow      time.Duration `toml:"malformedPacketWindow" mapstructure:"malformedPacketWindow"`
	MalformedPacketBanType     string        `toml:"malformedPacketBanType" mapstructure:"malformedPacketBanType"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		Enabled:                 true,
		DefaultTempBanDuration:  7 * 24 * time.Hour, // 7 days
		TempBansBeforePermanent: 3,
		IPBanMode:               "permanent_only",
		GMBansIncrementCounter:  false,
		ViolationDetection: ViolationDetectionConfig{
			RollingWindowDuration: 5 * time.Minute,
			CleanupInterval:       10 * time.Minute,
		},
		CombatDetection: CombatDetectionConfig{
			Enabled:                  true,
			ExcessiveDamageThreshold: 5,
			ExcessiveDamageWindow:    5 * time.Minute,
			ExcessiveDamageBanType:   "temporary",
			AttackSpeedThreshold:     5,
			AttackSpeedWindow:        5 * time.Minute,
			AttackSpeedBanType:       "temporary",
			InvalidSkillThreshold:    3,
			InvalidSkillWindow:       5 * time.Minute,
			InvalidSkillBanType:      "temporary",
		},
		MovementDetection: MovementDetectionConfig{
			Enabled:                  true,
			SpeedHackThreshold:       5,
			SpeedHackWindow:          5 * time.Minute,
			SpeedHackBanType:         "temporary",
			TeleportHackThreshold:    3,
			TeleportHackWindow:       5 * time.Minute,
			TeleportHackBanType:      "temporary",
			InvalidPositionThreshold: 5,
			InvalidPositionWindow:    5 * time.Minute,
			InvalidPositionBanType:   "temporary",
		},
		InventoryDetection: InventoryDetectionConfig{
			Enabled:                 true,
			InvalidEquipThreshold:   3,
			InvalidEquipWindow:      5 * time.Minute,
			InvalidEquipBanType:     "temporary",
			InvalidItemUseThreshold: 5,
			InvalidItemUseWindow:    5 * time.Minute,
			InvalidItemUseBanType:   "temporary",
		},
		EconomyDetection: EconomyDetectionConfig{
			Enabled:              true,
			InvalidTradeThreshold: 3,
			InvalidTradeWindow:   5 * time.Minute,
			InvalidTradeBanType:  "permanent",
			DuplicationThreshold: 1,
			DuplicationWindow:    5 * time.Minute,
			DuplicationBanType:   "permanent",
			OverflowThreshold:    3,
			OverflowWindow:       5 * time.Minute,
			OverflowBanType:      "permanent",
		},
		SkillDetection: SkillDetectionConfig{
			Enabled:                 true,
			CooldownBypassThreshold: 5,
			CooldownBypassWindow:    5 * time.Minute,
			CooldownBypassBanType:   "temporary",
			UnlearnedSkillThreshold: 3,
			UnlearnedSkillWindow:    5 * time.Minute,
			UnlearnedSkillBanType:   "temporary",
		},
		PacketDetection: PacketDetectionConfig{
			Enabled:                  true,
			InvalidSequenceThreshold: 10,
			InvalidSequenceWindow:    5 * time.Minute,
			InvalidSequenceBanType:   "temporary",
			MalformedPacketThreshold: 5,
			MalformedPacketWindow:    5 * time.Minute,
			MalformedPacketBanType:   "temporary",
		},
	}
}
