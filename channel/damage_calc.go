package channel

import (
	"math"
	"math/rand"

	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/nx"
)

// WeaponType represents the type of weapon being used
type WeaponType int

const (
	WeaponTypeNone     WeaponType = 0
	WeaponTypeSword1H  WeaponType = 30
	WeaponTypeAxe1H    WeaponType = 31
	WeaponTypeBW1H     WeaponType = 32
	WeaponTypeDagger   WeaponType = 33
	WeaponTypeWand     WeaponType = 37
	WeaponTypeStaff    WeaponType = 38
	WeaponTypeSword2H  WeaponType = 40
	WeaponTypeAxe2H    WeaponType = 41
	WeaponTypeBW2H     WeaponType = 42
	WeaponTypeSpear    WeaponType = 43
	WeaponTypePolearm  WeaponType = 44
	WeaponTypeBow      WeaponType = 45
	WeaponTypeCrossbow WeaponType = 46
	WeaponTypeClaw     WeaponType = 47
)

// AttackAction represents different attack animations
type AttackAction int

const (
	AttackActionSwing1H1 AttackAction = 0x05
	AttackActionSwing1H2 AttackAction = 0x06
	AttackActionSwing1H3 AttackAction = 0x07
	AttackActionSwing1H4 AttackAction = 0x08

	AttackActionSwing2H1 AttackAction = 0x09
	AttackActionSwing2H2 AttackAction = 0x0A
	AttackActionSwing2H3 AttackAction = 0x0B
	AttackActionSwing2H4 AttackAction = 0x0C
	AttackActionSwing2H5 AttackAction = 0x0D
	AttackActionSwing2H6 AttackAction = 0x0E
	AttackActionSwing2H7 AttackAction = 0x0F

	AttackActionStab1 AttackAction = 0x10
	AttackActionStab2 AttackAction = 0x11
	AttackActionStab3 AttackAction = 0x12
	AttackActionStab4 AttackAction = 0x13
	AttackActionStab5 AttackAction = 0x14
	AttackActionStab6 AttackAction = 0x15

	AttackActionBullet1 AttackAction = 0x16
	AttackActionBullet2 AttackAction = 0x17
	AttackActionBullet3 AttackAction = 0x18
	AttackActionBullet4 AttackAction = 0x19
	AttackActionBullet5 AttackAction = 0x1A
	AttackActionBullet6 AttackAction = 0x1B

	AttackActionProne AttackAction = 0x20
	AttackActionHeal  AttackAction = 0x28
	AttackActionUnk35 AttackAction = 0x35
)

// AttackOption represents flags for special attack properties
type AttackOption byte

const (
	AttackOptionNormal          AttackOption = 0
	AttackOptionSlashBlastFA    AttackOption = 1
	AttackOptionMortalBlowProp  AttackOption = 4
	AttackOptionShadowPartner   AttackOption = 8
	AttackOptionMortalBlowMelee AttackOption = 16
)

// CalcConstants holds calculation constants and modifiers
type CalcConstants struct{}

const (
	MaxHits    = 15
	MaxTargets = 15
	MaxPAD     = 999
)

var (
	// SlashBlastFAModifiers defines damage modifiers for Slash Blast Final Attack
	SlashBlastFAModifiers = [MaxTargets]float64{
		0.666667,
		0.222222,
		0.074074,
		0.024691,
		0.008229999999999,
		0.002743,
		0.000914,
		0.000305,
		0.000102,
		0.000033,
		0.000011,
		0.000004,
		0.000001,
		0.0,
		0.0,
	}

	// IronArrowModifiers defines damage modifiers for Iron Arrow skill
	IronArrowModifiers = [MaxTargets]float64{
		1.0,
		0.8,
		0.64,
		0.512,
		0.4096,
		0.32768,
		0.262144,
		0.209715,
		0.167772,
		0.134218,
		0.107374,
		0.085899,
		0.068719,
		0.054976,
		0.04398,
	}
)

// GetWeaponType extracts weapon type from item ID
func GetWeaponType(weaponID int32) WeaponType {
	if weaponID/1000000 != 1 {
		return WeaponTypeNone
	}
	weaponType := (weaponID / 10000) % 100
	if weaponType < 30 {
		return WeaponTypeNone
	}
	if weaponType > 33 && (weaponType <= 36 || weaponType > 38 && (weaponType <= 39 || weaponType > 47)) {
		return WeaponTypeNone
	}
	return WeaponType(weaponType)
}

// Roller handles random number generation for damage calculations
type Roller struct {
	rollIndex int
	rolls     []uint32
}

const (
	StatModifier = 0.000000100000010000001
	PropModifier = 0.0000100000010000001
)

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

// CalcHit represents a single hit calculation
type CalcHit struct {
	IsCrit  bool
	IsMiss  bool
	Damage  float64
	roller  *Roller
	info    *attackInfo
	data    *attackData
	mob     *monster
	player  *Player
	skillID int32
	str     int16
	dex     int16
	intl    int16
	luk     int16
	watk    int16

	critSkill        *nx.PlayerSkill
	critLevel        byte
	masteryModifier  float64
	targetAccuracy   float64
	weaponType       WeaponType
	attackAction     AttackAction
	skill            *nx.PlayerSkill
	isRanged         bool
	hitIdx           int
	attackType       int
	ampData          *ElementAmpData
	calcDamageOption byte
	hits             []*CalcHit // Reference to all hits for shadow partner
}

// CalcTargetAttack represents damage calculation for a single target
type CalcTargetAttack struct {
	Hits        []*CalcHit
	Mob         *monster
	Info        *attackInfo
	Character   *Player
	TotalDamage float64
}

// CalcDamage represents the complete damage calculation for an attack
type CalcDamage struct {
	skill         *nx.PlayerSkill
	Data          *attackData
	weaponType    WeaponType
	isRanged      bool
	character     *Player
	attackOption  AttackOption
	AttackAction  AttackAction
	TargetAttacks []*CalcTargetAttack
}

// IsJob checks if the character job matches the given job ID
func IsJob(charJobID int16, jobID int16) bool {
	if jobID%100 != 0 {
		return jobID/10 == charJobID/10 && charJobID%10 >= jobID%10
	}
	return jobID/100 == charJobID/100
}

// GetElementAmplification calculates element amplification data
func (c *CalcTargetAttack) GetElementAmplification(data *attackData) *ElementAmpData {
	jobID := c.Character.job
	ampSkillID := int32(0)
	skillID := data.skillID

	// Fire/Poison Mage
	if IsJob(jobID, 210) { // FPMage
		ampSkillID = int32(skill.ElementAmplification)
	} else if IsJob(jobID, 220) { // ILMage
		ampSkillID = int32(skill.ILElementAmplification)
	}

	ampData := &ElementAmpData{Magic: 100, Mana: 100}
	if ampSkillID > 0 {
		if ampSkillInfo, ok := c.Character.skills[ampSkillID]; ok {
			skillData, err := nx.GetPlayerSkill(ampSkillID)
			if err == nil && len(skillData) > 0 && ampSkillInfo.Level > 0 {
				idx := int(ampSkillInfo.Level) - 1
				if idx < len(skillData) {
					skipMP := false
					// Check if we should skip MP consumption based on skill ID
					if skillID > int32(skill.ElementComposition) {
						// Various checks for specific skills
						if skillID < int32(skill.ColdBeam) || skillID > int32(skill.ThunderBolt) {
							skipMP = true
						}
					}
					if !skipMP {
						ampData.Mana = int(skillData[idx].X)
					}
					ampData.Magic = int(skillData[idx].Y)
				}
			}
		}
	}
	return ampData
}

// GetTargetAccuracy calculates accuracy against the target
func (c *CalcTargetAttack) GetTargetAccuracy(data *attackData) float64 {
	levelDiff := int(c.Mob.level) - int(c.Character.level)
	if levelDiff < 0 {
		levelDiff = 0
	}

	var accuracy int
	if data.attackType == attackMagic {
		accuracy = 5 * (int(c.Character.intt)/10 + int(c.Character.luk)/10)
	} else {
		// Total ACC from stats and equipment
		accuracy = int(c.Character.dex) // Simplified - should include equipment bonuses
	}
	return float64(accuracy*100) / (float64(levelDiff*10) + 255.0)
}

// NewCalcDamage creates a new damage calculation
func NewCalcDamage(chr *Player, data *attackData, attackType int) *CalcDamage {
	calc := &CalcDamage{
		Data:          data,
		character:     chr,
		isRanged:      attackType == attackRanged,
		AttackAction:  AttackAction(data.action),
		TargetAttacks: make([]*CalcTargetAttack, len(data.attackInfo)),
	}

	// Get weapon type
	weaponID := int32(0)
	for _, item := range chr.equip {
		if item.slotID == -11 { // Weapon slot
			weaponID = item.ID
			break
		}
	}
	calc.weaponType = GetWeaponType(weaponID)

	// Get skill data
	skillID := data.skillID
	if skillID > 0 {
		if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
			if data.skillLevel > 0 && int(data.skillLevel) <= len(skillData) {
				// Store reference to the specific skill level data
				calc.skill = &skillData[data.skillLevel-1]
			}
		}
	}

	masteryModifier := calc.GetMasteryModifier()
	critLevel, critSkill := calc.GetCritSkill()
	calc.attackOption = AttackOption(data.option)

	// Calculate weapon attack
	watk := int16(math.Min(float64(MaxPAD), float64(chr.getTotalWatk())))

	// Handle mortal blow option
	if (calc.attackOption & (AttackOptionMortalBlowProp | AttackOptionMortalBlowMelee)) != 0 {
		chr.rng.Uint32() // Consume one random number
	}

	// Calculate damage for each target
	for targetIdx := range calc.TargetAttacks {
		info := &data.attackInfo[targetIdx]
		
		// Find the mob in the field
		if chr.inst == nil {
			continue
		}
		mob, err := chr.inst.lifePool.getMobFromID(info.spawnID)
		if err != nil {
			continue
		}

		roller := NewRoller(chr.rng, 7)
		targetDmg := &CalcTargetAttack{
			Character: chr,
			Mob:       &mob,
			Info:      info,
			Hits:      make([]*CalcHit, len(info.damages)),
		}

		ampData := targetDmg.GetElementAmplification(data)
		targetAccuracy := targetDmg.GetTargetAccuracy(data)

		// Calculate each hit
		for hitIdx := range targetDmg.Hits {
			hit := &CalcHit{
				player:           chr,
				roller:           roller,
				info:             info,
				data:             data,
				mob:              &mob,
				critSkill:        critSkill,
				critLevel:        critLevel,
				masteryModifier:  masteryModifier,
				targetAccuracy:   targetAccuracy,
				weaponType:       calc.weaponType,
				attackAction:     calc.AttackAction,
				skill:            calc.skill,
				isRanged:         calc.isRanged,
				hitIdx:           hitIdx,
				skillID:          skillID,
				str:              chr.str,
				dex:              chr.dex,
				intl:             chr.intt,
				luk:              chr.luk,
				watk:             watk,
				attackType:       attackType,
				ampData:          ampData,
				calcDamageOption: info.calcDamageStatIndex,
				hits:             targetDmg.Hits,
			}

			hit.Calculate()
			targetDmg.Hits[hitIdx] = hit
			targetDmg.TotalDamage += hit.Damage
		}

		calc.TargetAttacks[targetIdx] = targetDmg
		if calc.ApplyAfterModifiers(targetDmg, targetIdx) {
			break
		}
	}

	return calc
}

// GetMasteryModifier returns the mastery modifier for the attack
func (c *CalcDamage) GetMasteryModifier() float64 {
	var mastery int
	if c.Data.attackType == attackMagic {
		if c.skill != nil {
			mastery = int(c.skill.Mastery)
		}
	} else {
		mastery = c.GetWeaponMastery()
	}
	return (float64(mastery)*5.0 + 10.0) * 0.009000000000000001
}

// ApplyAfterModifiers applies skill-specific damage modifiers after calculation
func (c *CalcDamage) ApplyAfterModifiers(targetDmg *CalcTargetAttack, targetIdx int) bool {
	if c.skill == nil {
		return false
	}

	var modifier float64
	brk := false

	if c.attackOption == AttackOptionSlashBlastFA {
		modifier = SlashBlastFAModifiers[targetIdx]
	} else if c.Data.skillID == int32(skill.ArrowBomb) {
		if targetIdx > 0 {
			modifier = float64(c.skill.X) * 0.01
		} else {
			modifier = 0.5
			if targetDmg.TotalDamage == 0 {
				brk = true
			}
		}
	} else if c.Data.skillID == int32(skill.IronArrow) {
		modifier = IronArrowModifiers[targetIdx]
	} else {
		return false
	}

	for _, hit := range targetDmg.Hits {
		hit.Damage *= modifier
	}
	return brk
}

// GetCritSkill returns critical hit skill data
func (c *CalcDamage) GetCritSkill() (byte, *nx.PlayerSkill) {
	if !c.isRanged {
		return 0, nil
	}

	var skillID int32
	switch c.weaponType {
	case WeaponTypeBow, WeaponTypeCrossbow:
		skillID = int32(skill.CriticalShot)
	case WeaponTypeClaw:
		skillID = int32(skill.CriticalThrow)
	default:
		return 0, nil
	}

	if skillInfo, ok := c.character.skills[skillID]; ok {
		if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
			if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
				return skillInfo.Level, &skillData[skillInfo.Level-1]
			}
		}
	}
	return 0, nil
}

// GetWeaponMastery returns weapon mastery value
func (c *CalcDamage) GetWeaponMastery() int {
	// Check if weapon type matches attack type
	switch c.weaponType {
	case WeaponTypeBow, WeaponTypeCrossbow, WeaponTypeClaw:
		if !c.isRanged {
			return 0
		}
	default:
		if c.isRanged {
			return 0
		}
	}

	var skillID int32
	switch c.weaponType {
	case WeaponTypeSword1H, WeaponTypeSword2H:
		if c.character.job/10 == 11 { // Fighter
			skillID = int32(skill.SwordMastery)
		} else {
			skillID = int32(skill.PageSwordMastery) // Page
		}
	case WeaponTypeAxe1H, WeaponTypeAxe2H:
		skillID = int32(skill.AxeMastery)
	case WeaponTypeBW1H, WeaponTypeBW2H:
		skillID = int32(skill.BwMastery)
	case WeaponTypeDagger:
		skillID = int32(skill.DaggerMastery)
	case WeaponTypeSpear:
		skillID = int32(skill.SpearMastery)
	case WeaponTypePolearm:
		skillID = int32(skill.PolearmMastery)
	case WeaponTypeBow:
		skillID = int32(skill.BowMastery)
	case WeaponTypeCrossbow:
		skillID = int32(skill.CrossbowMastery)
	case WeaponTypeClaw:
		skillID = int32(skill.ClawMastery)
	default:
		return 0
	}

	if skillID != 0 {
		if skillInfo, ok := c.character.skills[skillID]; ok {
			if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
				if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
					return int(skillData[skillInfo.Level-1].Mastery)
				}
			}
		}
	}
	return 0
}

// Calculate performs the complete damage calculation for a hit
func (h *CalcHit) Calculate() {
	// Check for immunity
	if h.attackType == attackMagic && h.mob.invincible {
		h.Damage = 1
		return
	}

	// Handle summon attacks
	if h.attackType == attackSummon {
		h.CalcSummonDamage()
		return
	}

	// Handle meso explosion
	if skill.Skill(h.skillID) == skill.MesoExplosion {
		h.CalcMesoExplosion()
		return
	}

	// Check for miss
	if h.GetIsMiss() {
		return
	}

	// Apply base damage calculation
	h.ApplyBaseDamage()

	// Apply mob level modifier for physical attacks
	if h.attackType != attackMagic {
		h.ApplyMobLevelModifier()
	}

	// Apply elemental modifiers
	h.ApplySpecialElementModifiers(h.Damage * float64(h.ampData.Magic) * 0.01)

	// Apply charge element modifiers
	if h.attackType != attackMagic {
		h.ApplyChargeElementModifiers()
	}

	// Apply mob defense reduction
	h.ApplyMobDefReduction()

	if h.attackType != attackMagic {
		baseDmg := int(h.Damage)
		h.ApplySkillDamage()
		h.ApplyComboAttack()
		h.ApplyCrit(baseDmg)
		
		// Apply mob power guard
		// (mob buffs not implemented yet)

		h.ApplyShadowPartner()
	}

	// Clamp damage to valid range
	maxDamage := 999999.0 // Config max damage
	h.Damage = math.Min(maxDamage, math.Max(1, h.Damage))
}

// CalcMesoExplosion calculates meso explosion damage
func (h *CalcHit) CalcMesoExplosion() {
	if h.hitIdx >= len(h.info.mesoDropIDs) {
		h.Damage = 0
		return
	}

	// Get meso amount from drop (simplified)
	mesosUsed := float64(100) // Placeholder - need to look up actual drop

	var mesoModifier float64
	if mesosUsed <= 1000 {
		mesoModifier = (mesosUsed*0.82 + 28.0) * 0.0001886792452830189
	} else {
		mesoModifier = mesosUsed / (mesosUsed + 5250)
	}

	skillDmgX := float64(h.skill.X)
	roll := h.roller.Roll(0.0000000500000050000005)
	h.Damage = (50 * skillDmgX) * (roll + 0.5) * mesoModifier
}

// CalcSummonDamage calculates summon damage
func (h *CalcHit) CalcSummonDamage() {
	summonID := h.data.summonType
	
	// Check if this is a dragon summon
	if summonID == int32(skill.SummonDragon) {
		h.CalcDragonDamage()
		return
	}

	// Regular summon damage
	summonSkillInfo, ok := h.player.skills[summonID]
	if !ok {
		return
	}

	skillData, err := nx.GetPlayerSkill(summonID)
	if err != nil || len(skillData) == 0 || summonSkillInfo.Level == 0 {
		return
	}

	idx := int(summonSkillInfo.Level) - 1
	if idx >= len(skillData) {
		return
	}

	summonSkill := &skillData[idx]
	roll := h.roller.Roll(0.00000003000000300000031)
	h.Damage = (float64(h.dex)*(roll+0.7)*2.5 + float64(h.str)) * float64(summonSkill.Damage) * 0.01
}

// CalcDragonDamage calculates dragon summon damage
func (h *CalcHit) CalcDragonDamage() {
	dragonSkillID := int32(skill.SummonDragon)
	dragonSkillInfo, ok := h.player.skills[dragonSkillID]
	if !ok {
		return
	}

	skillData, err := nx.GetPlayerSkill(dragonSkillID)
	if err != nil || len(skillData) == 0 || dragonSkillInfo.Level == 0 {
		return
	}

	idx := int(dragonSkillInfo.Level) - 1
	if idx >= len(skillData) {
		return
	}

	dragonSkill := &skillData[idx]

	magic := float64(h.player.intt) // Simplified - should include equipment MAD
	totalInt := float64(h.player.intt)
	
	// Use mob's FS (Final Stat) for mastery modifier
	statModifier := (float64(h.mob.pdDamage)*5.0 + 10.0) * 0.009000000000000001
	rolledStat := h.RollStat(magic, statModifier, 1.0)

	h.Damage = (magic*0.058*magic*0.058 + totalInt*0.5 + rolledStat*3.3) * float64(dragonSkill.Damage) * 0.01
}

// GetIsMiss calculates if the attack misses
func (h *CalcHit) GetIsMiss() bool {
	roll := h.roller.Roll(StatModifier)

	var minModifier, maxModifier float64
	if h.attackType == attackMagic {
		minModifier = 0.5
		maxModifier = 1.2
	} else {
		minModifier = 0.7
		maxModifier = 1.3
	}

	minTACC := h.targetAccuracy * minModifier
	randTACC := minTACC
	maxTACC := h.targetAccuracy * maxModifier

	randTACC += (maxTACC - randTACC) * roll
	mobAvoid := math.Min(999, float64(h.mob.eva))

	h.IsMiss = randTACC < mobAvoid
	return h.IsMiss
}

// ApplyBaseDamage calculates the base damage
func (h *CalcHit) ApplyBaseDamage() {
	// Check special cases
	if h.BowMeleeBaseDmg() || h.ClawMeleeBaseDmg() || h.ProneBaseDmg() {
		return
	}

	if h.attackType == attackMagic {
		h.CalcMagicBaseDamage()
		return
	}

	var statModifier float64
	isSwing := h.attackAction >= AttackActionSwing1H1 && h.attackAction <= AttackActionSwing2H7

	switch h.weaponType {
	case WeaponTypeBow:
		dex := h.RollStat(float64(h.dex), h.masteryModifier, 1.0)
		statModifier = dex*3.4 + float64(h.str)

	case WeaponTypeCrossbow:
		dex := h.RollStat(float64(h.dex), h.masteryModifier, 1.0)
		statModifier = dex*3.6 + float64(h.str)

	case WeaponTypeAxe2H, WeaponTypeBW2H:
		str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
		if isSwing {
			statModifier = str*4.8 + float64(h.dex)
		} else {
			statModifier = str*3.4 + float64(h.dex)
		}

	case WeaponTypeSpear, WeaponTypePolearm:
		str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
		if skill.Skill(h.skillID) == skill.DragonRoar {
			statModifier = str*4.0 + float64(h.dex)
		} else if isSwing != (h.weaponType == WeaponTypeSpear) {
			statModifier = str*5.0 + float64(h.dex)
		} else {
			statModifier = str*3.0 + float64(h.dex)
		}

	case WeaponTypeSword2H:
		str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
		statModifier = str*4.6 + float64(h.dex)

	case WeaponTypeAxe1H, WeaponTypeBW1H, WeaponTypeWand, WeaponTypeStaff:
		str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
		if isSwing {
			statModifier = str*4.4 + float64(h.dex)
		} else {
			statModifier = str*3.2 + float64(h.dex)
		}

	case WeaponTypeSword1H, WeaponTypeDagger:
		if h.player.job/100 == 4 && h.weaponType == WeaponTypeDagger {
			luk := h.RollStat(float64(h.luk), h.masteryModifier, 1.0)
			secondary := float64(h.str + h.dex)
			statModifier = h.CalcStatModifier(luk, secondary)
		} else {
			str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
			statModifier = str*4.0 + float64(h.dex)
		}

	case WeaponTypeClaw:
		if skill.Skill(h.skillID) == skill.LuckySeven {
			luk := h.RollStat(float64(h.luk), 0.5, 1.0)
			statModifier = luk * 5.0
		} else if skill.Skill(h.skillID) == skill.ShadowMeso {
			// Shadow Meso special calculation
			moneyCon := float64(h.skill.X) * 0.5
			stat := h.RollStat(moneyCon, h.masteryModifier, 3.0)
			h.Damage = 10.0 * math.Floor(stat)

			propRoll := h.roller.Roll(PropModifier)
			if float64(h.skill.Prop) > propRoll {
				h.IsCrit = true
				bonusDmg := 100 + h.skill.X
				h.Damage *= float64(bonusDmg) * 0.01
			}
			return
		} else {
			luk := h.RollStat(float64(h.luk), h.masteryModifier, 1.0)
			secondary := float64(h.dex + h.str)
			statModifier = h.CalcStatModifier(luk, secondary)
		}

	default:
		return
	}

	h.Damage = statModifier * float64(h.watk) * 0.01
}

// CalcStatModifier calculates stat modifier for primary/secondary stats
func (h *CalcHit) CalcStatModifier(primary, secondary float64) float64 {
	return secondary + primary*3.6
}

// CalcMagicBaseDamage calculates base magic damage
func (h *CalcHit) CalcMagicBaseDamage() {
	totalMAD := int16(math.Min(999, float64(h.player.intt))) // Simplified - should include equipment

	if skill.Skill(h.skillID) == skill.Heal {
		targets := len(h.data.attackInfo) + 1
		rolledStat := h.RollStat(float64(h.intl), 0.8, 0.2)
		h.Damage = (rolledStat*1.5 + float64(h.luk)) *
			(float64(targets)*0.3 + 1.0) *
			(float64(h.skill.Prop) * 0.01) *
			float64(totalMAD) *
			0.005 /
			float64(targets)
	} else {
		rolledStat := h.RollStat(float64(totalMAD), h.masteryModifier, 1.0)
		h.Damage = (float64(h.intl)*0.5 + float64(totalMAD)*0.058*float64(totalMAD)*0.058 + rolledStat*3.3) *
			float64(h.skill.Damage) *
			0.01
	}
}

// BowMeleeBaseDmg calculates bow melee damage
func (h *CalcHit) BowMeleeBaseDmg() bool {
	if h.weaponType != WeaponTypeBow && h.weaponType != WeaponTypeCrossbow {
		return false
	}
	if (AttackActionBullet1 <= h.attackAction && h.attackAction <= AttackActionBullet6) || h.attackAction == AttackActionUnk35 {
		return false
	}

	dex := h.RollStat(float64(h.dex), h.masteryModifier, 1.0)
	if skill.Skill(h.skillID) != skill.CBPowerKnockback && skill.Skill(h.skillID) != skill.PowerKnockback {
		h.Damage = (float64(h.str) + dex) * float64(h.watk) * 0.005
	} else {
		h.Damage = (dex*3.4 + float64(h.str)) * float64(h.watk) * 0.005
	}
	return true
}

// ClawMeleeBaseDmg calculates claw melee damage
func (h *CalcHit) ClawMeleeBaseDmg() bool {
	if h.weaponType != WeaponTypeClaw || (AttackActionBullet1 <= h.attackAction && h.attackAction <= AttackActionBullet6) || h.attackAction == AttackActionUnk35 {
		return false
	}

	luk := h.RollStat(float64(h.luk), h.masteryModifier, 1.0)
	h.Damage = (float64(h.str) + float64(h.dex) + luk) * float64(h.watk) * 0.006666666666666667
	return true
}

// ProneBaseDmg calculates prone attack damage
func (h *CalcHit) ProneBaseDmg() bool {
	if h.attackAction != AttackActionProne {
		return false
	}
	str := h.RollStat(float64(h.str), h.masteryModifier, 1.0)
	h.Damage = (float64(h.dex) + str) * float64(h.watk) * 0.005
	return true
}

// ApplyMobLevelModifier applies level difference modifier
func (h *CalcHit) ApplyMobLevelModifier() {
	if int(h.player.level) < int(h.mob.level) {
		diff := int(h.mob.level) - int(h.player.level)
		h.Damage = (100.0 - float64(diff)) * h.Damage * 0.01
	}
}

// ElementModifier represents element resistance/weakness
type ElementModifier int

const (
	ElementModifierNormal     ElementModifier = 0
	ElementModifierNullify    ElementModifier = 1
	ElementModifierHalf       ElementModifier = 2
	ElementModifierOneAndHalf ElementModifier = 3
)

// ApplySpecialElementModifiers applies elemental damage modifiers
func (h *CalcHit) ApplySpecialElementModifiers(dmg float64) {
	if h.skill == nil {
		return
	}

	var newDmg float64
	if skill.Skill(h.skillID) == skill.ElementComposition || skill.Skill(h.skillID) == skill.ILElementComposition {
		// Element composition splits damage between two elements
		halfDmg := dmg * 0.5
		
		// Determine which elements to use
		var elements []int
		if skill.Skill(h.skillID) == skill.ElementComposition {
			elements = []int{2, 3} // Fire, Poison
		} else {
			elements = []int{1, 4} // Ice, Lightning
		}

		total := 0.0
		for _, elem := range elements {
			// Get mob's resistance to this element
			// Simplified - actual implementation would read from mob.elemAttr
			_ = elem // Mark as used
			mobModifier := ElementModifierNormal
			total += h.ApplyMobElemModifier(halfDmg, mobModifier, 1.0)
		}
		newDmg = total
	} else {
		modifier := 1.0
		// Special skills with element-specific modifiers
		if skill.Skill(h.skillID) == skill.Inferno || skill.Skill(h.skillID) == skill.Blizzard {
			skillLevel := h.data.skillLevel
			modifier = float64(20+skillLevel) * 0.0099999998
		}

		// Get element from skill and check mob resistance
		// Simplified - actual implementation would use skill.ElementFlags
		elemModifier := ElementModifierNormal
		newDmg = h.ApplyMobElemModifier(dmg, elemModifier, modifier)
	}
	h.Damage = newDmg
}

// ApplyMobElemModifier applies mob element modifier to damage
func (h *CalcHit) ApplyMobElemModifier(dmg float64, modifier ElementModifier, extraModifier float64) float64 {
	newDmg := dmg

	if modifier == ElementModifierNullify {
		return 0.0
	} else if modifier == ElementModifierHalf {
		return (1.0 - extraModifier*0.5) * dmg
	} else if modifier == ElementModifierOneAndHalf {
		newDmg = (extraModifier*0.5 + 1.0) * dmg
		if dmg >= newDmg {
			newDmg = dmg
		}
		maxDamage := 999999.0 // Config max damage
		newDmg = math.Min(maxDamage, newDmg)
	}
	return newDmg
}

// ApplyChargeElementModifiers applies charge skill element modifiers
func (h *CalcHit) ApplyChargeElementModifiers() {
	if h.player.job/10 != 12 { // Not White Knight
		return
	}

	if h.player.buffs == nil {
		return
	}

	// Check for active charge buff
	// Simplified - actual implementation would check specific buff
	var chargeSkillID int32
	hasCharge := false
	
	// Look for any charge skill buff
	chargeSkills := []skill.Skill{
		skill.SwordFireCharge, skill.SwordIceCharge, skill.SwordLitCharge,
		skill.BwFireCharge, skill.BwIceCharge, skill.BwLitCharge,
	}
	
	for _, chargeSkill := range chargeSkills {
		if _, ok := h.player.buffs.activeSkillLevels[int32(chargeSkill)]; ok {
			chargeSkillID = int32(chargeSkill)
			hasCharge = true
			break
		}
	}

	if !hasCharge {
		return
	}

	chargeSkillInfo, ok := h.player.skills[chargeSkillID]
	if !ok {
		return
	}

	skillData, err := nx.GetPlayerSkill(chargeSkillID)
	if err != nil || len(skillData) == 0 || chargeSkillInfo.Level == 0 {
		return
	}

	idx := int(chargeSkillInfo.Level) - 1
	if idx >= len(skillData) {
		return
	}

	chargeSkill := &skillData[idx]

	specialModifier := float64(chargeSkill.Z) * 0.0099999998
	damageModifier := float64(chargeSkill.Damage) * 0.0099999998
	
	// Get element from charge skill and check mob resistance
	elemModifier := ElementModifierNormal // Simplified
	dmg := damageModifier * h.Damage
	h.Damage = h.ApplyMobElemModifier(dmg, elemModifier, specialModifier)
}

// ApplyMobDefReduction applies mob defense reduction
func (h *CalcHit) ApplyMobDefReduction() {
	if skill.Skill(h.skillID) == skill.Sacrifice ||
		skill.Skill(h.skillID) == skill.Assaulter {
		return
	}

	var mobDef float64
	if h.attackType == attackMagic {
		mobDef = float64(h.mob.mdDamage)
	} else {
		mobDef = float64(h.mob.pdDamage)
	}
	mobDef = math.Min(999, mobDef)

	redMin := mobDef * 0.5
	redMax := mobDef * 0.6
	reduction := redMin
	statRoll := h.roller.Roll(StatModifier)
	reduction += (redMax - reduction) * statRoll
	h.Damage -= reduction
}

// ApplySkillDamage applies skill damage multiplier
func (h *CalcHit) ApplySkillDamage() {
	if h.skill == nil || h.skill.Damage <= 0 {
		return
	}
	h.Damage *= float64(h.skill.Damage) * 0.01
}

// ApplyComboAttack applies combo attack bonus
func (h *CalcHit) ApplyComboAttack() {
	comboSkillID := int32(skill.ComboAttack)
	comboSkillInfo, hasCombo := h.player.skills[comboSkillID]
	if !hasCombo {
		return
	}

	// Get current combo orbs from player buffs
	currOrbs := 0
	if h.player.buffs != nil {
		// Try to get combo count from buffs - simplified for now
		// In actual implementation, this would read from buff state
		currOrbs = 0 // Placeholder - needs proper buff tracking
	}

	if currOrbs <= 0 {
		return
	}

	skillData, err := nx.GetPlayerSkill(comboSkillID)
	if err != nil || len(skillData) == 0 || comboSkillInfo.Level == 0 {
		return
	}

	idx := int(comboSkillInfo.Level) - 1
	if idx >= len(skillData) {
		return
	}

	comboDmg := skillData[idx].Damage
	var modifier float64

	if currOrbs == 1 {
		modifier = float64(comboDmg)
	} else if h.skillID >= int32(skill.SwordPanic) && h.skillID <= int32(skill.AxeComa) {
		// Special combo skills with different scaling
		switch currOrbs {
		case 2:
			modifier = float64(24*comboSkillInfo.Level-24)/29 + float64(comboDmg) + 6
		case 3:
			modifier = float64(int(comboSkillInfo.Level)<<6-64)/29 + float64(comboDmg) + 16
		case 4:
			modifier = float64(120*comboSkillInfo.Level-120)/29 + float64(comboDmg) + 30
		case 5:
			modifier = float64(184*comboSkillInfo.Level-184)/29 + float64(comboDmg) + 46
		default:
			return
		}
	} else {
		// Regular combo attack bonus
		var dmgBonus float64
		switch currOrbs {
		case 2:
			dmgBonus = float64(5*comboSkillInfo.Level - 5)
		case 3:
			dmgBonus = float64(10*comboSkillInfo.Level - 10)
		case 4:
			dmgBonus = float64(15*comboSkillInfo.Level - 15)
		case 5:
			dmgBonus = float64(20*comboSkillInfo.Level - 20)
		default:
			return
		}
		modifier = float64(comboDmg) + dmgBonus/29
	}
	h.Damage *= modifier * 0.01
}

// ApplyCrit applies critical hit bonus
func (h *CalcHit) ApplyCrit(baseDmg int) {
	// Some skills don't crit
	if skill.Skill(h.skillID) == skill.Blizzard ||
		skill.Skill(h.skillID) == skill.ShadowMeso {
		return
	}

	if h.critSkill == nil || h.critSkill.Prop <= 0 || h.critSkill.Damage <= 0 || h.critLevel <= 0 {
		return
	}

	roll := h.roller.Roll(PropModifier)
	if roll < float64(h.critSkill.Prop) {
		h.IsCrit = true
		critBonus := h.critSkill.Damage - 100
		h.Damage += float64(critBonus) * 0.01 * float64(baseDmg)
	}
}

// ApplyShadowPartner applies shadow partner damage
func (h *CalcHit) ApplyShadowPartner() {
	if !h.isRanged {
		return
	}

	// Check if shadow partner buff is active
	spSkillID := int32(skill.ShadowPartner)
	if h.player.buffs == nil {
		return
	}
	if _, hasSP := h.player.buffs.activeSkillLevels[spSkillID]; !hasSP {
		return
	}

	spSkillInfo, hasSkill := h.player.skills[spSkillID]
	if !hasSkill {
		return
	}

	skillData, err := nx.GetPlayerSkill(spSkillID)
	if err != nil || len(skillData) == 0 || spSkillInfo.Level == 0 {
		return
	}

	idx := int(spSkillInfo.Level) - 1
	if idx >= len(skillData) {
		return
	}

	spSkill := &skillData[idx]

	// Shadow partner creates a clone of hits
	spHits := len(h.info.damages) / 2
	if h.hitIdx < spHits {
		return
	}

	var dmgModifier int
	if h.skill != nil {
		dmgModifier = int(spSkill.Y)
	} else {
		dmgModifier = int(spSkill.X)
	}

	// Clone damage from the first half of hits
	clonedHit := h.hits[h.hitIdx-spHits]
	h.Damage = float64(dmgModifier) * clonedHit.Damage / 100.0
	h.IsCrit = clonedHit.IsCrit
}

// RollStat rolls a stat with mastery and modifier
func (h *CalcHit) RollStat(stat float64, masteryModifier float64, statModifier float64) float64 {
	statRoll := h.roller.Roll(StatModifier)
	modifiedStat := stat * statModifier
	masteryStat := stat * masteryModifier
	
	if modifiedStat == masteryStat {
		return modifiedStat
	}
	
	least := math.Min(modifiedStat, masteryStat)
	diff := math.Abs(modifiedStat - masteryStat)
	modifiedStat = least + diff*statRoll
	return modifiedStat
}

// getTotalWatk returns total weapon attack (simplified)
func (p *Player) getTotalWatk() int16 {
	// Simplified - should include equipment bonuses
	return 50 // Placeholder
}
