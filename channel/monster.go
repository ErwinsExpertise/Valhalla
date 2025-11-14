package channel

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

type monster struct {
	controller, summoner   *Player
	id                     int32
	spawnID                int32
	pos                    pos
	faceLeft               bool
	hp, mp                 int32
	maxHP, maxMP           int32
	hpRecovery, mpRecovery int32
	level                  int32
	exp                    int32
	maDamage               int32
	mdDamage               int32
	paDamage               int32
	pdDamage               int32
	summonType             int8 // -2: fade in spawn animation, -1: no spawn animation, 0: balrog summon effect?
	summonOption           int32
	boss                   bool
	undead                 bool
	elemAttr               int32
	invincible             bool
	speed                  int32
	eva                    int32
	acc                    int32
	link                   int32
	flySpeed               int32
	noRegen                int32
	skills                 map[byte]byte
	revives                []int32
	stance                 byte
	poison                 bool

	lastAttackTime int64
	lastSkillTime  int64
	skillTimes     map[byte]int64

	skillID    byte
	skillLevel byte
	statBuff   int32
	
	// Player-inflicted debuffs
	debuffs map[int32]*mobDebuff // key: stat mask from MobStat
	debuffExpireTimers map[int32]*time.Timer

	dmgTaken map[*Player]int32

	dropsItems bool
	dropsMesos bool

	hpBgColour byte
	hpFgColour byte

	spawnInterval int64
	timeToSpawn   time.Time

	lastStatusUpdate int64
	lastHeal         int64
	lastTimeAttacked int64
}

func createMonsterFromData(spawnID int32, life nx.Life, m nx.Mob, dropsItems, dropsMesos bool) monster {
	return monster{
		id:                 life.ID,
		spawnID:            spawnID,
		pos:                newPos(life.X, life.Y, life.Foothold),
		faceLeft:           life.FaceLeft,
		hp:                 m.HP,
		mp:                 m.MP,
		maxHP:              m.MaxHP,
		maxMP:              m.MaxMP,
		exp:                int32(m.Exp),
		revives:            m.Revives,
		summonType:         constant.MobSummonTypeRegen,
		boss:               m.Boss > 0,
		hpBgColour:         byte(m.HPTagBGColor),
		hpFgColour:         byte(m.HPTagColor),
		spawnInterval:      life.MobTime,
		dmgTaken:           make(map[*Player]int32),
		skills:             nx.GetMobSkills(life.ID),
		skillTimes:         make(map[byte]int64),
		debuffs:            make(map[int32]*mobDebuff),
		debuffExpireTimers: make(map[int32]*time.Timer),
		poison:             false,
		lastHeal:           time.Now().Unix(),
		lastSkillTime:      0,
	}
}

func createMonsterFromID(spawnID, id int32, p pos, controller *Player, dropsItems, dropsMesos bool, summoner int32) (monster, error) {
	m, err := nx.GetMob(id)
	if err != nil {
		return monster{}, fmt.Errorf("Unknown mob ID: %v", id)
	}

	// If this isn't working with regards to position make the foothold equal to Player? nearest to pos?
	mob := createMonsterFromData(spawnID, nx.Life{ID: id, Foothold: p.foothold, X: p.x, Y: p.y, FaceLeft: true}, m, dropsItems, dropsMesos)
	mob.summoner = controller

	return mob, nil
}

func (m *monster) setController(controller *Player, follow bool) {
	if controller == nil {
		return
	}
	m.controller = controller
	controller.Send(packetMobControl(*m, follow))
}

func (m *monster) removeController() {
	if m.controller != nil {
		m.controller.Send(packetMobEndControl(*m))
		m.controller = nil
	}
}

func (m *monster) acknowledgeController(moveID int16, movData movementFrag, allowedToUseSkill bool, skill, level byte) {
	m.pos.x = movData.x
	m.pos.y = movData.y
	m.pos.foothold = movData.foothold
	m.stance = movData.stance
	m.faceLeft = m.stance%2 == 1

	if m.controller == nil {
		return
	}

	// Clamp MP to int16 range to avoid overflow
	mp16 := int16(math.MaxInt16)
	if m.mp < int32(math.MaxInt16) {
		mp16 = int16(m.mp)
	}

	m.controller.Send(packetMobControlAcknowledge(m.spawnID, moveID, allowedToUseSkill, mp16, skill, level))
}

func (m monster) hasHPBar() (bool, int32, int32, int32, byte, byte) {
	return (m.boss && m.hpBgColour > 0), m.id, m.hp, m.maxHP, m.hpFgColour, m.hpBgColour
}

func (m *monster) getMobSkill(delay int16, skillLevel, skillID byte) (byte, byte, nx.MobSkill) {
	// If sealed, cannot use skills
	if (m.statBuff & skill.MobStat.SealSkill) > 0 {
		return 0, 0, nx.MobSkill{}
	}

	levels, err := nx.GetMobSkill(skillID)
	if err != nil {
		m.skillID = 0
		return 0, 0, nx.MobSkill{}
	}

	if skillLevel == 0 || int(skillLevel) > len(levels) {
		return 0, 0, nx.MobSkill{}
	}
	skillData := levels[skillLevel-1]

	m.mp -= skillData.MpCon
	if m.mp < 0 {
		m.mp = 0
	}

	return skillID, skillLevel, skillData
}

func (m *monster) giveDamage(damager *Player, dmg ...int32) {
	for _, v := range dmg {
		if v > m.hp {
			v = m.hp
		}
		m.hp -= v

		if damager != nil {
			m.dmgTaken[damager] += v
		}
	}
	// Always update lastTimeAttacked
	m.lastTimeAttacked = time.Now().Unix()
}

func (m monster) displayBytes() []byte {
	p := mpacket.NewPacket()

	p.WriteInt32(m.spawnID)
	p.WriteByte(0x00) // control status?
	p.WriteInt32(m.id)

	p.WriteInt32(0) // some kind of status?

	p.WriteInt16(m.pos.x)
	p.WriteInt16(m.pos.y)

	var bitfield byte
	if m.summoner != nil {
		bitfield = 0x08
	} else {
		bitfield = 0x02
	}

	if m.faceLeft {
		bitfield |= 0x01
	} else {
		bitfield |= 0x04
	}

	if m.stance%2 == 1 {
		bitfield |= 0x01
	}

	if m.flySpeed > 0 {
		bitfield |= 0x04
	}

	p.WriteByte(bitfield)        // 0x08 - a summon, 0x04 - flying, 0x02 - ???, 0x01 - faces left
	p.WriteInt16(m.pos.foothold) // foothold to oscillate around
	p.WriteInt16(m.pos.foothold) // spawn foothold
	p.WriteInt8(m.summonType)

	if m.summonType == constant.MobSummonTypeRevive || m.summonType >= 0 {
		p.WriteInt32(m.summonOption) // when -3 used to link mob to a death using spawnID
	}

	p.WriteInt32(0) // encode mob status
	return p
}

func (m monster) String() string {
	sid := strconv.Itoa(int(m.spawnID))
	mid := strconv.Itoa(int(m.id))

	hp := strconv.Itoa(int(m.hp))
	mhp := strconv.Itoa(int(m.maxHP))

	mp := strconv.Itoa(int(m.mp))
	mmp := strconv.Itoa(int(m.maxMP))

	return sid + "(" + mid + ") " + hp + "/" + mhp + " " + mp + "/" + mmp + " (" + m.pos.String() + ")"
}

func (m *monster) update(t time.Time) {
	checkTime := t.Unix()
	m.lastStatusUpdate = checkTime

	if m.hp <= 0 {
		return
	}

	if m.poison {
		// Handle poison (TODO: scale by poison level)
		m.hp -= 10
		if m.hp < 0 {
			m.hp = 0
		}
	}

	// Periodic regen
	if (checkTime - m.lastHeal) > 30 {
		regenhp, regenmp := m.calculateHeal()
		m.healMob(regenhp, regenmp)
		m.lastHeal = checkTime
	}
}

func (mob *monster) chooseNextSkill() (byte, byte) {
	var chosenID, chosenLevel byte
	if (mob.statBuff&skill.MobStat.SealSkill) > 0 || (time.Now().Unix()-mob.lastSkillTime) < 10 {
		return 0, 0
	}

	candidates := make([]byte, 0, len(mob.skills))
	for id, lvl := range mob.skills {
		levels, err := nx.GetMobSkill(id)
		if err != nil {
			continue
		}
		if lvl == 0 || int(lvl) > len(levels) {
			continue
		}
		skillData := levels[lvl-1]

		// MP check
		if mob.mp < skillData.MpCon {
			continue
		}
		// Cooldown check
		if last, ok := mob.skillTimes[id]; ok {
			if last+skillData.Interval > time.Now().Unix() {
				continue
			}
		}

		// Skip buffs already active
		if mob.statBuff > 0 {
			alreadySet := false
			switch id {
			case skill.Mob.WeaponAttackUp, skill.Mob.WeaponAttackUpAoe:
				alreadySet = (mob.statBuff & skill.MobStat.PowerUp) > 0
			case skill.Mob.MagicAttackUp, skill.Mob.MagicAttackUpAoe:
				alreadySet = (mob.statBuff & skill.MobStat.MagicUp) > 0
			case skill.Mob.WeaponDefenceUp, skill.Mob.WeaponDefenceUpAoe:
				alreadySet = (mob.statBuff & skill.MobStat.PowerGuardUp) > 0
			case skill.Mob.MagicDefenceUp, skill.Mob.MagicDefenceUpAoe:
				alreadySet = (mob.statBuff & skill.MobStat.MagicGuardUp) > 0
			case skill.Mob.WeaponImmunity:
				alreadySet = (mob.statBuff & skill.MobStat.PhysicalImmune) > 0
			case skill.Mob.MagicImmunity:
				alreadySet = (mob.statBuff & skill.MobStat.MagicImmune) > 0
			case skill.Mob.McSpeedUp:
				alreadySet = (mob.statBuff & skill.MobStat.Speed) > 0
			default:
			}
			if alreadySet {
				continue
			}
		}

		candidates = append(candidates, id)
	}

	if len(candidates) > 0 {
		chosenID = candidates[rand.Intn(len(candidates))]
		chosenLevel = mob.skills[chosenID]
	}

	if chosenLevel == 0 {
		chosenID = 0
	}

	return chosenID, chosenLevel
}

func (m *monster) healMob(hp, mp int32) {
	if hp > 0 && m.hp < m.maxHP {
		newHP := m.hp + hp
		if newHP > m.maxHP {
			newHP = m.maxHP
		} else if newHP < 0 {
			newHP = 0
		}
		m.hp = newHP
	}

	if mp > 0 && m.mp < m.maxMP {
		newMP := m.mp + mp
		if newMP > m.maxMP {
			newMP = m.maxMP
		} else if newMP < 0 {
			newMP = 0
		}
		m.mp = newMP
	}
}

func (m monster) calculateHeal() (hp int32, mp int32) {
	// Base regen: 1% per tick
	hp = m.maxHP / 100
	mp = m.maxMP / 100

	// Always allow MP regen (mobs need MP to attack).
	// Only allow HP regen if not recently attacked (60s grace).
	if time.Now().Unix()-m.lastTimeAttacked < 60 {
		return 0, mp
	}
	return hp, mp
}

func packetMobControl(m monster, chase bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelControlMob)
	if chase {
		p.WriteByte(0x02) // 2 chase, 1 no chase, 0 no control
	} else {
		p.WriteByte(0x01)
	}

	p.Append(m.displayBytes())

	return p
}

func packetMobControlAcknowledge(mobID int32, moveID int16, allowedToUseSkill bool, mp int16, skill byte, level byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelControlMobAck)
	p.WriteInt32(mobID)
	p.WriteInt16(moveID)
	p.WriteBool(allowedToUseSkill)
	p.WriteInt16(mp) // Protocol appears to expect 16-bit here; value is clamped at call site
	p.WriteByte(skill)
	p.WriteByte(level)

	return p
}

func packetMobEndControl(m monster) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelControlMob)
	p.WriteByte(0)
	p.WriteInt32(m.spawnID)

	return p
}

func packetMobShowHpChange(spawnID int32, dmg int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelMobDamage)
	p.WriteInt32(spawnID)
	p.WriteByte(0)
	p.WriteInt32(dmg)

	return p
}

// mobDebuff represents a debuff applied to a mob by a player skill
type mobDebuff struct {
	skillID   int32
	value     int16
	duration  int16 // in seconds
	expiresAt int64 // unix milliseconds
}

// applyDebuff applies a debuff to the mob from a player skill
func (m *monster) applyDebuff(skillID int32, skillLevel byte, statMask int32, inst *fieldInstance) {
	if m.debuffs == nil {
		m.debuffs = make(map[int32]*mobDebuff)
	}
	if m.debuffExpireTimers == nil {
		m.debuffExpireTimers = make(map[int32]*time.Timer)
	}

	// Get skill data
	skillData, err := nx.GetPlayerSkill(skillID)
	if err != nil || skillLevel == 0 || int(skillLevel) > len(skillData) {
		return
	}
	
	si := skillData[skillLevel-1]
	
	// Calculate debuff value based on skill type
	var value int16
	switch skill.Skill(skillID) {
	case skill.Threaten, skill.ArmorCrash, skill.PowerCrash, skill.MagicCrash:
		// These reduce mob stats by X%
		value = int16(si.X)
	case skill.Slow, skill.ILSlow:
		// Slow reduces speed by X%
		value = int16(si.X)
	case skill.Seal, skill.ILSeal:
		// Seal prevents skill use
		value = 1
	case skill.ShadowWeb:
		// Shadow Web immobilizes
		value = int16(si.X)
	case skill.Doom:
		// Doom
		value = int16(si.X)
	default:
		value = 1
	}

	duration := int16(si.Time)
	if duration <= 0 {
		duration = 30 // default 30 seconds
	}

	expiresAt := time.Now().Add(time.Duration(duration) * time.Second).UnixMilli()

	// Store the debuff
	m.debuffs[statMask] = &mobDebuff{
		skillID:   skillID,
		value:     value,
		duration:  duration,
		expiresAt: expiresAt,
	}

	// Update the statBuff mask
	m.statBuff |= statMask

	// Cancel existing timer if any
	if timer, ok := m.debuffExpireTimers[statMask]; ok && timer != nil {
		timer.Stop()
	}

	// Set expiration timer
	if inst != nil && inst.dispatch != nil {
		m.debuffExpireTimers[statMask] = time.AfterFunc(time.Duration(duration)*time.Second, func() {
			inst.dispatch <- func() {
				m.removeDebuff(statMask, inst)
			}
		})
	}

	// Send packet to show the debuff
	inst.send(packetMobStatSet(m.spawnID, statMask, value, skillID, duration))
}

// removeDebuff removes a debuff from the mob
func (m *monster) removeDebuff(statMask int32, inst *fieldInstance) {
	if m.debuffs == nil {
		return
	}

	// Remove from debuffs map
	delete(m.debuffs, statMask)

	// Update statBuff mask
	m.statBuff &^= statMask

	// Cancel timer
	if timer, ok := m.debuffExpireTimers[statMask]; ok && timer != nil {
		timer.Stop()
		delete(m.debuffExpireTimers, statMask)
	}

	// Send packet to remove the debuff
	if inst != nil {
		inst.send(packetMobStatReset(m.spawnID, statMask))
	}
}

// packetMobStatSet sends a packet to apply a debuff to a mob
// Following the OpenMG encoding pattern exactly
func packetMobStatSet(spawnID int32, statMask int32, value int16, skillID int32, duration int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelMobStatSet)
	p.WriteInt32(spawnID)
	
	// Write the stat mask (this gets written first, then we write data for each bit)
	p.WriteUint32(uint32(statMask))
	
	// Convert duration from seconds to deciseconds (100ms units)
	durationDeciseconds := duration * 10
	
	// For each bit set in the mask, write data IN ORDER
	// Order matches OpenMG's Encode method
	if (statMask & skill.MobStat.PhysicalDamage) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.PhysicalDefense) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.MagicDamage) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.MagicDefense) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Accurrency) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Evasion) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Speed) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Stun) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Freeze) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Poison) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Seal) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Darkness) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.PowerUp) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.MagicUp) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.PowerGuardUp) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.MagicGuardUp) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.PhysicalImmune) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.MagicImmune) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Doom) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Web) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.HardSkin) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Ambush) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Venom) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.Blind) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	if (statMask & skill.MobStat.SealSkill) != 0 {
		p.WriteInt16(value)
		p.WriteInt32(skillID)
		p.WriteInt16(durationDeciseconds)
	}
	
	// Write delay at the end (in milliseconds, typically 0)
	p.WriteInt16(0)

	return p
}

// packetMobStatReset sends a packet to remove a debuff from a mob
func packetMobStatReset(spawnID int32, statMask int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelMobStatReset)
	p.WriteInt32(spawnID)
	p.WriteInt32(statMask)
	p.WriteByte(1) // reset flag

	return p
}
