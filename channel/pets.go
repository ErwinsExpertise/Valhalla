package channel

import (
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/mpacket"
)

type pet struct {
	name      string
	itemID    int32
	dbID      int64
	level     byte
	closeness int16
	fullness  byte
	deadDate  int64
	spawnDate int64

	pos    pos
	stance byte
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

func packetPetMove(charID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetMove)
	p.WriteInt32(charID)
	// Encode move... https://github.com/sewil/OpenMG/blob/main/WvsBeta.Common/Objects/MovePath.cs#L143
	return p
}

func packetPetSpawn(charID int32, petData *pet) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetSpawn)
	p.WriteInt32(charID)
	p.WriteBool(true)
	p.WriteInt32(petData.itemID)
	p.WriteString(petData.name)
	p.WriteInt64(petData.dbID)
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

func packetPlayerPetUpdate(petID int64) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStatChange)
	p.WriteBool(false)
	p.WriteInt32(constant.PetID)
	p.WriteInt64(petID)

	return p
}
