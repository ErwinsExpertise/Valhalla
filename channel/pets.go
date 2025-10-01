package channel

import (
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
)

type pet struct {
	name            string
	itemID          int32
	sn              int32
	itemDBID        int64
	level           byte
	closeness       int16
	fullness        byte
	deadDate        int64
	spawnDate       int64
	lastInteraction int64

	pos    pos
	stance byte

	spawned bool
}

func newPet(itemID, sn int32, dbID int64) *pet {
	return &pet{
		name:            "",
		itemID:          itemID,
		sn:              sn,
		itemDBID:        dbID,
		stance:          0,
		level:           0,
		closeness:       0,
		fullness:        100,
		deadDate:        (time.Now().UnixMilli()*10000 + 116444592000000000 + (time.Hour.Milliseconds() * 24 * 90)),
		spawnDate:       0,
		lastInteraction: 0,
	}
}

func savePet(item *Item) error {
	if item.petData == nil {
		return nil
	}

	p := item.petData

	_, err := common.DB.Exec(`
		INSERT INTO pets (
			itemID, name, sn, level, closeness, fullness,
			deadDate, spawnDate, lastInteraction
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			level = VALUES(level),
			closeness = VALUES(closeness),
			fullness = VALUES(fullness),
			deadDate = VALUES(deadDate),
			spawnDate = VALUES(spawnDate),
			lastInteraction = VALUES(lastInteraction)
	`, item.dbID,
		p.name,
		p.sn,
		p.level,
		p.closeness,
		p.fullness,
		p.deadDate,
		p.spawnDate,
		p.lastInteraction,
	)
	return err
}

func (p *pet) updateMovement(frag movementFrag) {
	p.pos.x = frag.x
	p.pos.y = frag.y
	p.pos.foothold = frag.foothold
	p.stance = frag.stance
}

const (
	petRemoveNone   byte = 0
	petRemoveHungry byte = 1
	petRemoveExpire byte = 2

	petSpawn      byte = 1
	petConnect    byte = 2
	petShowRemote byte = 4
	petChangeMap  byte = 8

	petResetPos    = petSpawn | petConnect | petChangeMap
	petResetHunger = petSpawn | petConnect
	petResetStat   = petSpawn | petConnect
)

func packetPetAction(charID int32, op, action byte, text string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetAction)
	p.WriteInt32(charID)
	p.WriteByte(op)
	p.WriteByte(action)
	p.WriteString(text)
	return p
}

func packetPetNameChange(charID int32, name string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetNameChange)
	p.WriteInt32(charID)
	p.WriteString(name)
	return p
}

func packetPetInteraction(charID int32, interactionId byte, inc, food bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetInteraction)
	p.WriteInt32(charID)
	p.WriteBool(food)
	if !food {
		p.WriteByte(interactionId)
	}
	p.WriteBool(inc)

	return p
}

func packetPetMove(charID int32, move []byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetMove)
	p.WriteInt32(charID)
	p.WriteBytes(move)
	return p
}

func packetPetSpawn(charID int32, petData *pet) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetSpawn)
	p.WriteInt32(charID)
	p.WriteBool(true)
	p.WriteInt32(petData.itemID)
	p.WriteString(petData.name)
	p.WriteUint64(uint64(petData.sn))
	p.WriteInt16(petData.pos.x)
	p.WriteInt16(petData.pos.y)
	p.WriteByte(petData.stance)
	p.WriteInt16(petData.pos.foothold)
	return p
}

func packetPetRemove(charID int32, reason byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetSpawn)
	p.WriteInt32(charID)
	p.WriteBool(false)
	p.WriteByte(reason)

	return p
}
