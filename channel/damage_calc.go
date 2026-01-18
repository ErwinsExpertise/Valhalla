package channel

import (
	"log"
	"math"
	"math/rand"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/nx"
)

// DamageRange represents the min and max damage for validation
type DamageRange struct {
	Min float64
	Max float64
}

// CalcHitResult represents the result of a hit calculation
type CalcHitResult struct {
	IsCrit       bool
	IsMiss       bool
	MinDamage    float64
	MaxDamage    float64
	ExpectedDmg  float64 // For reference/logging
	ClientDamage int32   // The damage client sent
	IsValid      bool    // Whether client damage is within acceptable range
}

// Roller handles random number generation for damage calculations
type Roller struct {
	rollIndex int
	rolls     []uint32
}

// NewRoller creates a new roller with pre-generated random numbers
func NewRoller(randomizer *rand.Rand, numRolls int) *Roller {
	rolls := make([]uint32, numRolls)
	for i := 0; i < numRolls; i++ {
		rolls[i] = randomizer.Uint32()
	}
	return &Roller{
		rollIndex: 0,
		rolls:     rolls,
	}
}

// Roll returns a random value using the specified modifier
func (r *Roller) Roll(modifier float64) float64 {
	// Nil check: return a mid-range value if roller is not initialized
	if r == nil || len(r.rolls) == 0 {
		return 0.5 // Return middle of the typical 0-1 range
	}
	
	idx := r.rollIndex % len(r.rolls)
	r.rollIndex++
	roll := r.rolls[idx]
	rollValue := float64(roll%10000000) * modifier
	return rollValue
}

// ElementAmpData stores element amplification data
type ElementAmpData struct {
	Magic int
	Mana  int
}

// DamageCalculator handles all damage calculations
type DamageCalculator struct {
	player       *Player
	data         *attackData
	attackType   int
	weaponType   constant.WeaponType
	skill        *nx.PlayerSkill
	skillID      int32
	skillLevel   byte
	isRanged     bool
	masteryMod   float64
	critSkill    *nx.PlayerSkill
	critLevel    byte
	watk         int16
	attackAction constant.AttackAction
	attackOption constant.AttackOption
}

// NewDamageCalculator creates a new damage calculator
func NewDamageCalculator(plr *Player, data *attackData, attackType int) *DamageCalculator {
	calc := &DamageCalculator{
		player:       plr,
		data:         data,
		attackType:   attackType,
		isRanged:     attackType == attackRanged,
		skillID:      data.skillID,
		skillLevel:   data.skillLevel,
		attackAction: constant.AttackAction(data.action),
		attackOption: constant.AttackOption(data.option),
	}

	// Get weapon type
	weaponID := int32(0)
	for _, item := range plr.equip {
		if item.slotID == -11 { // Weapon slot
			weaponID = item.ID
			break
		}
	}
	calc.weaponType = constant.GetWeaponType(weaponID)

	// Get skill data
	if data.skillID > 0 {
		if skillData, err := nx.GetPlayerSkill(data.skillID); err == nil && len(skillData) > 0 {
			if data.skillLevel > 0 && int(data.skillLevel) <= len(skillData) {
				calc.skill = &skillData[data.skillLevel-1]
			}
		}
	}

	// Calculate mastery modifier
	calc.masteryMod = calc.GetMasteryModifier()

	// Get critical skill
	calc.critLevel, calc.critSkill = calc.GetCritSkill()

	// Calculate weapon attack
	calc.watk = calc.GetTotalWatk()

	return calc
}

// ValidateAttack validates all hits in an attack and determines critical hits
func (calc *DamageCalculator) ValidateAttack() [][]CalcHitResult {
	results := make([][]CalcHitResult, len(calc.data.attackInfo))

	for targetIdx := range calc.data.attackInfo {
		info := &calc.data.attackInfo[targetIdx]

		// Find the mob
		if calc.player.inst == nil {
			continue
		}
		mob, err := calc.player.inst.lifePool.getMobFromID(info.spawnID)
		if err != nil {
			continue
		}

		// Create roller for this target
		roller := NewRoller(calc.player.rng, constant.DamageRollsPerTarget)

		// Get element amplification
		ampData := calc.GetElementAmplification()

		// Calculate accuracy
		targetAccuracy := calc.GetTargetAccuracy(&mob)

		// Validate each hit
		results[targetIdx] = make([]CalcHitResult, len(info.damages))
		for hitIdx := range info.damages {
			results[targetIdx][hitIdx] = calc.CalculateHit(
				&mob,
				info,
				roller,
				ampData,
				targetAccuracy,
				hitIdx,
				targetIdx,
			)
		}
	}

	return results
}

// CalculateHit calculates damage range for a single hit and validates client damage
// Following the 10-step formula from the requirement
func (calc *DamageCalculator) CalculateHit(
	mob *monster,
	info *attackInfo,
	roller *Roller,
	ampData *ElementAmpData,
	targetAccuracy float64,
	hitIdx int,
	targetIdx int,
) CalcHitResult {
	result := CalcHitResult{
		ClientDamage: info.damages[hitIdx],
		IsValid:      false,
	}

	// Check for immunity
	if calc.attackType == attackMagic && mob.invincible {
		result.MinDamage = 1
		result.MaxDamage = 1
		result.IsValid = (result.ClientDamage == 1)
		return result
	}

	// Check for miss
	if calc.GetIsMiss(roller, targetAccuracy, mob) {
		result.IsMiss = true
		result.MinDamage = 0
		result.MaxDamage = 0
		result.IsValid = (result.ClientDamage == 0)
		return result
	}

	// STEP 1: Calculate min and max damage from base formula
	minDmg, maxDmg := calc.CalculateBaseDamageRange(mob, hitIdx)

	// STEP 2: Multiply by skill modifiers (including elemental)
	minDmg, maxDmg = calc.ApplySkillModifiers(minDmg, maxDmg, ampData, mob)

	// STEP 3: Calculate defense reduction
	defReduction := calc.CalculateDefenseReduction(mob, roller)
	minDmg -= defReduction
	maxDmg -= defReduction

	// Ensure damage doesn't go negative
	if minDmg < 0 {
		minDmg = 0
	}
	if maxDmg < 0 {
		maxDmg = 0
	}

	// STEP 4: Select random number from damage range (done by client, we validate)
	// We'll calculate expected damage as midpoint for reference
	baseDmg := (minDmg + maxDmg) / 2.0

	// STEP 5: Find damage multiplier (skill damage%)
	multiplier := 1.0
	if calc.skill != nil && calc.skill.Damage > 0 {
		multiplier = float64(calc.skill.Damage) / 100.0
	}

	// STEP 6: Add critical bonuses and determine if this is a crit
	result.IsCrit = calc.CheckCritical(roller)
	critMultiplier := 1.0
	if result.IsCrit && calc.critSkill != nil {
		critBonus := float64(calc.critSkill.Damage-100) / 100.0
		critMultiplier = 1.0 + critBonus
	}

	// STEP 7: Multiply with damage multiplier and crit
	totalMultiplier := multiplier * critMultiplier
	minDmg *= totalMultiplier
	maxDmg *= totalMultiplier
	baseDmg *= totalMultiplier

	// STEP 8: Clamp to [1, 99999]
	if minDmg < 1 && minDmg > 0 {
		minDmg = 1
	}
	if maxDmg > 99999 {
		maxDmg = 99999
	}
	if minDmg > 99999 {
		minDmg = 99999
	}

	// STEP 9: Apply after-modifiers (multi-target skills)
	afterMod := calc.GetAfterModifier(targetIdx, baseDmg)
	minDmg *= afterMod
	maxDmg *= afterMod
	baseDmg *= afterMod

	// STEP 10: Take integer part
	minDmg = math.Floor(minDmg)
	maxDmg = math.Floor(maxDmg)

	result.MinDamage = minDmg
	result.MaxDamage = maxDmg
	result.ExpectedDmg = baseDmg

	// Validate client damage - only care if it's OVER the acceptable range
	// If damage is under, we don't care (client may have weaker gear, buffs expired, etc.)
	tolerance := constant.DamageVarianceTolerance
	toleranceMax := maxDmg * (1.0 + tolerance)
	
	clientDmgFloat := float64(result.ClientDamage)
	result.IsValid = (clientDmgFloat <= toleranceMax)

	// Log and cap suspiciously high damage
	if !result.IsValid {
		log.Printf("Suspicious high damage from player %s (ID: %d): client=%d, max expected=%.0f (with tolerance), skill=%d",
			calc.player.Name, calc.player.ID, result.ClientDamage, toleranceMax, calc.skillID)
	}

	return result
}

// CalculateBaseDamageRange calculates min and max base damage
func (calc *DamageCalculator) CalculateBaseDamageRange(mob *monster, hitIdx int) (float64, float64) {
	// Get stat values
	str := float64(calc.player.str)
	dex := float64(calc.player.dex)
	luk := float64(calc.player.luk)
	watk := float64(calc.watk)

	// Magic damage uses different formula
	if calc.attackType == attackMagic {
		return calc.CalculateMagicDamageRange()
	}

	// Get mastery range (min = full mastery, max = full stat)
	masteryMin := calc.masteryMod
	masteryMax := 1.0

	var minStatMod, maxStatMod float64

	isSwing := calc.attackAction >= constant.AttackActionSwing1H1 && calc.attackAction <= constant.AttackActionSwing2H7

	switch calc.weaponType {
	case constant.WeaponTypeBow2:
		minStatMod = dex*masteryMin*3.4 + str
		maxStatMod = dex*masteryMax*3.4 + str

	case constant.WeaponTypeCrossbow2:
		minStatMod = dex*masteryMin*3.6 + str
		maxStatMod = dex*masteryMax*3.6 + str

	case constant.WeaponTypeAxe2H, constant.WeaponTypeBW2H:
		if isSwing {
			minStatMod = str*masteryMin*4.8 + dex
			maxStatMod = str*masteryMax*4.8 + dex
		} else {
			minStatMod = str*masteryMin*3.4 + dex
			maxStatMod = str*masteryMax*3.4 + dex
		}

	case constant.WeaponTypeSpear2, constant.WeaponTypePolearm2:
		if skill.Skill(calc.skillID) == skill.DragonRoar {
			minStatMod = str*masteryMin*4.0 + dex
			maxStatMod = str*masteryMax*4.0 + dex
		} else if isSwing != (calc.weaponType == constant.WeaponTypeSpear2) {
			minStatMod = str*masteryMin*5.0 + dex
			maxStatMod = str*masteryMax*5.0 + dex
		} else {
			minStatMod = str*masteryMin*3.0 + dex
			maxStatMod = str*masteryMax*3.0 + dex
		}

	case constant.WeaponTypeSword2H:
		minStatMod = str*masteryMin*4.6 + dex
		maxStatMod = str*masteryMax*4.6 + dex

	case constant.WeaponTypeAxe1H, constant.WeaponTypeBW1H, constant.WeaponTypeWand2, constant.WeaponTypeStaff2:
		if isSwing {
			minStatMod = str*masteryMin*4.4 + dex
			maxStatMod = str*masteryMax*4.4 + dex
		} else {
			minStatMod = str*masteryMin*3.2 + dex
			maxStatMod = str*masteryMax*3.2 + dex
		}

	case constant.WeaponTypeSword1H, constant.WeaponTypeDagger2:
		if calc.player.job/100 == 4 && calc.weaponType == constant.WeaponTypeDagger2 {
			minStatMod = luk*masteryMin*3.6 + str + dex
			maxStatMod = luk*masteryMax*3.6 + str + dex
		} else {
			minStatMod = str*masteryMin*4.0 + dex
			maxStatMod = str*masteryMax*4.0 + dex
		}

	case constant.WeaponTypeClaw2:
		if skill.Skill(calc.skillID) == skill.LuckySeven {
			// Lucky Seven uses 0.5 mastery
			minStatMod = luk * 0.5 * 5.0
			maxStatMod = luk * 5.0
		} else {
			minStatMod = luk*masteryMin*3.6 + str + dex
			maxStatMod = luk*masteryMax*3.6 + str + dex
		}

	default:
		return 0, 0
	}

	minDmg := minStatMod * watk * 0.01
	maxDmg := maxStatMod * watk * 0.01

	// Apply mob level modifier for physical attacks
	if int(calc.player.level) < int(mob.level) {
		levelPenalty := (100.0 - float64(int(mob.level)-int(calc.player.level))) / 100.0
		minDmg *= levelPenalty
		maxDmg *= levelPenalty
	}

	return minDmg, maxDmg
}

// CalculateMagicDamageRange calculates magic damage range
func (calc *DamageCalculator) CalculateMagicDamageRange() (float64, float64) {
	totalMAD := float64(math.Min(999, float64(calc.player.intt))) // Simplified
	intl := float64(calc.player.intt)
	luk := float64(calc.player.luk)

	if skill.Skill(calc.skillID) == skill.Heal {
		// Heal uses special formula
		targets := float64(len(calc.data.attackInfo) + 1)
		
		minStat := intl * 0.8 * 0.2
		maxStat := intl * 0.2
		
		minDmg := (minStat*1.5 + luk) * (targets*0.3 + 1.0) * float64(calc.skill.Prop) * 0.01 * totalMAD * 0.005 / targets
		maxDmg := (maxStat*1.5 + luk) * (targets*0.3 + 1.0) * float64(calc.skill.Prop) * 0.01 * totalMAD * 0.005 / targets
		
		return minDmg, maxDmg
	}

	// Standard magic formula
	minMAD := totalMAD * calc.masteryMod
	maxMAD := totalMAD

	minDmg := (intl*0.5 + totalMAD*0.058*totalMAD*0.058 + minMAD*3.3) * float64(calc.skill.Damage) * 0.01
	maxDmg := (intl*0.5 + totalMAD*0.058*totalMAD*0.058 + maxMAD*3.3) * float64(calc.skill.Damage) * 0.01

	return minDmg, maxDmg
}

// ApplySkillModifiers applies skill-specific modifiers including elemental
func (calc *DamageCalculator) ApplySkillModifiers(minDmg, maxDmg float64, ampData *ElementAmpData, mob *monster) (float64, float64) {
	if calc.skill == nil {
		return minDmg, maxDmg
	}

	// Apply element amplification for magic
	if calc.attackType == attackMagic {
		elemMod := float64(ampData.Magic) / 100.0
		minDmg *= elemMod
		maxDmg *= elemMod
	}

	// Apply charge element modifiers for White Knight
	// (simplified - full implementation would check active buffs)

	return minDmg, maxDmg
}

// CalculateDefenseReduction calculates defense reduction amount
func (calc *DamageCalculator) CalculateDefenseReduction(mob *monster, roller *Roller) float64 {
	// Skip defense for certain skills
	if skill.Skill(calc.skillID) == skill.Sacrifice ||
		skill.Skill(calc.skillID) == skill.Assaulter {
		return 0
	}

	var mobDef float64
	if calc.attackType == attackMagic {
		mobDef = float64(mob.mdDamage)
	} else {
		mobDef = float64(mob.pdDamage)
	}
	mobDef = math.Min(999, mobDef)

	// Defense reduces damage by 50-60% with RNG
	redMin := mobDef * 0.5
	redMax := mobDef * 0.6
	
	// Use roller for variance
	roll := roller.Roll(constant.DamageStatModifier)
	reduction := redMin + (redMax-redMin)*roll

	return reduction
}

// CheckCritical determines if a hit is critical
func (calc *DamageCalculator) CheckCritical(roller *Roller) bool {
	if !calc.isRanged || calc.critSkill == nil {
		return false
	}

	// Skills that don't crit
	if skill.Skill(calc.skillID) == skill.Blizzard {
		return false
	}

	roll := roller.Roll(constant.DamagePropModifier)
	return roll < float64(calc.critSkill.Prop)
}

// GetAfterModifier gets after-modifiers for multi-target skills
func (calc *DamageCalculator) GetAfterModifier(targetIdx int, baseDmg float64) float64 {
	if calc.skill == nil {
		return 1.0
	}

	if calc.attackOption == constant.AttackOptionSlashBlastFA {
		return constant.SlashBlastFAModifiers[targetIdx]
	}

	if calc.skillID == int32(skill.ArrowBomb) {
		if targetIdx > 0 {
			return float64(calc.skill.X) * 0.01
		}
		// First target gets 50% if it dealt damage
		if baseDmg > 0 {
			return 0.5
		}
		return 0
	}

	if calc.skillID == int32(skill.IronArrow) {
		return constant.IronArrowModifiers[targetIdx]
	}

	return 1.0
}

// GetIsMiss determines if attack misses
func (calc *DamageCalculator) GetIsMiss(roller *Roller, targetAccuracy float64, mob *monster) bool {
	roll := roller.Roll(constant.DamageStatModifier)

	var minModifier, maxModifier float64
	if calc.attackType == attackMagic {
		minModifier = 0.5
		maxModifier = 1.2
	} else {
		minModifier = 0.7
		maxModifier = 1.3
	}

	minTACC := targetAccuracy * minModifier
	randTACC := minTACC + (targetAccuracy*maxModifier-minTACC)*roll
	mobAvoid := math.Min(999, float64(mob.eva))

	return randTACC < mobAvoid
}

// GetElementAmplification calculates element amplification
func (calc *DamageCalculator) GetElementAmplification() *ElementAmpData {
	jobID := calc.player.job
	ampSkillID := int32(0)

	if jobID/10 == 21 { // FPMage
		ampSkillID = int32(skill.ElementAmplification)
	} else if jobID/10 == 22 { // ILMage
		ampSkillID = int32(skill.ILElementAmplification)
	}

	ampData := &ElementAmpData{Magic: 100, Mana: 100}
	if ampSkillID > 0 {
		if ampSkillInfo, ok := calc.player.skills[ampSkillID]; ok {
			skillData, err := nx.GetPlayerSkill(ampSkillID)
			if err == nil && len(skillData) > 0 && ampSkillInfo.Level > 0 {
				idx := int(ampSkillInfo.Level) - 1
				if idx < len(skillData) {
					ampData.Mana = int(skillData[idx].X)
					ampData.Magic = int(skillData[idx].Y)
				}
			}
		}
	}
	return ampData
}

// GetTargetAccuracy calculates accuracy against target
func (calc *DamageCalculator) GetTargetAccuracy(mob *monster) float64 {
	levelDiff := int(mob.level) - int(calc.player.level)
	if levelDiff < 0 {
		levelDiff = 0
	}

	var accuracy int
	if calc.attackType == attackMagic {
		accuracy = int(5 * (calc.player.intt/10 + calc.player.luk/10))
	} else {
		accuracy = int(calc.player.dex) // Simplified
	}

	return float64(accuracy*100) / (float64(levelDiff*10) + 255.0)
}

// GetMasteryModifier returns mastery modifier
func (calc *DamageCalculator) GetMasteryModifier() float64 {
	var mastery int
	if calc.attackType == attackMagic {
		if calc.skill != nil {
			mastery = int(calc.skill.Mastery)
		}
	} else {
		mastery = calc.GetWeaponMastery()
	}
	return (float64(mastery)*5.0 + 10.0) * 0.009000000000000001
}

// GetWeaponMastery returns weapon mastery value
func (calc *DamageCalculator) GetWeaponMastery() int {
	// Check weapon type matches attack type
	switch calc.weaponType {
	case constant.WeaponTypeBow2, constant.WeaponTypeCrossbow2, constant.WeaponTypeClaw2:
		if !calc.isRanged {
			return 0
		}
	default:
		if calc.isRanged {
			return 0
		}
	}

	var skillID int32
	switch calc.weaponType {
	case constant.WeaponTypeSword1H, constant.WeaponTypeSword2H:
		if calc.player.job/10 == 11 {
			skillID = int32(skill.SwordMastery)
		} else {
			skillID = int32(skill.PageSwordMastery)
		}
	case constant.WeaponTypeAxe1H, constant.WeaponTypeAxe2H:
		skillID = int32(skill.AxeMastery)
	case constant.WeaponTypeBW1H, constant.WeaponTypeBW2H:
		skillID = int32(skill.BwMastery)
	case constant.WeaponTypeDagger2:
		skillID = int32(skill.DaggerMastery)
	case constant.WeaponTypeSpear2:
		skillID = int32(skill.SpearMastery)
	case constant.WeaponTypePolearm2:
		skillID = int32(skill.PolearmMastery)
	case constant.WeaponTypeBow2:
		skillID = int32(skill.BowMastery)
	case constant.WeaponTypeCrossbow2:
		skillID = int32(skill.CrossbowMastery)
	case constant.WeaponTypeClaw2:
		skillID = int32(skill.ClawMastery)
	default:
		return 0
	}

	if skillID != 0 {
		if skillInfo, ok := calc.player.skills[skillID]; ok {
			if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
				if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
					return int(skillData[skillInfo.Level-1].Mastery)
				}
			}
		}
	}
	return 0
}

// GetCritSkill returns critical skill data
func (calc *DamageCalculator) GetCritSkill() (byte, *nx.PlayerSkill) {
	if !calc.isRanged {
		return 0, nil
	}

	var skillID int32
	switch calc.weaponType {
	case constant.WeaponTypeBow2, constant.WeaponTypeCrossbow2:
		skillID = int32(skill.CriticalShot)
	case constant.WeaponTypeClaw2:
		skillID = int32(skill.CriticalThrow)
	default:
		return 0, nil
	}

	if skillInfo, ok := calc.player.skills[skillID]; ok {
		if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
			if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
				return skillInfo.Level, &skillData[skillInfo.Level-1]
			}
		}
	}
	return 0, nil
}

// GetTotalWatk returns total weapon attack
func (calc *DamageCalculator) GetTotalWatk() int16 {
	watk := int16(0)
	
	for _, item := range calc.player.equip {
		if item.slotID == -11 { // Weapon slot
			watk += item.watk
		} else if item.slotID < 0 { // Other equipped items
			watk += item.watk
		}
	}
	
	// Add base STR contribution
	watk += calc.player.str / 10
	
	return int16(math.Min(float64(constant.DamageMaxPAD), float64(watk)))
}
