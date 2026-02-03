package channel

import (
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
)

type coconutType byte

const (
	coconutSpawn   coconutType = 0
	coconutHit     coconutType = 1
	coconutBreak   coconutType = 2
	coconutDestroy coconutType = 3
)

func packetSnowballState(state byte, x1 int16, hp1 byte, x2 int16, hp2 byte, snowballDmg, snowmanDmg1, snowmanDmg2 int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSnowballState)
	p.WriteByte(state)
	p.WriteInt16(x1)
	p.WriteByte(hp1)
	p.WriteInt16(x2)
	p.WriteByte(hp2)
	p.WriteInt16(snowballDmg)
	p.WriteInt16(snowmanDmg1)
	p.WriteInt16(snowmanDmg2)
	return p
}

func packetSnowballHit(hitType byte, damage, delay int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSnowballHit)
	p.WriteByte(hitType)
	p.WriteInt16(damage)
	p.WriteInt16(delay)
	return p
}

func packetCoconutAttack(coconutID, delay int16, kind coconutType) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCoconutAttack)
	p.WriteInt16(coconutID)
	p.WriteInt16(delay)
	p.WriteByte(byte(kind))
	return p
}

func packetCoconutScore(maple, story int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCoconutScore)
	p.WriteInt16(maple)
	p.WriteInt16(story)
	return p
}

func packetFieldSpecificData(team byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelFieldSpecificData)
	p.WriteByte(team)
	return p
}
