package channel

import (
	"sort"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

func packetCashShopSet(plr *player, accountName string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSetCashShop)

	// CharacterDataFlag: Stats|Money|MaxSlots|Items
	p.WriteInt16(0x00BF)

	// Stats
	p.WriteInt32(plr.id)
	p.WritePaddedString(plr.name, 13)
	p.WriteByte(plr.gender)
	p.WriteByte(plr.skin)
	p.WriteInt32(plr.face)
	p.WriteInt32(plr.hair)
	p.WriteInt64(0) // Pet Cash ID

	p.WriteByte(plr.level)
	p.WriteInt16(plr.job)
	p.WriteInt16(plr.str)
	p.WriteInt16(plr.dex)
	p.WriteInt16(plr.intt)
	p.WriteInt16(plr.luk)
	p.WriteInt16(plr.hp)
	p.WriteInt16(plr.maxHP)
	p.WriteInt16(plr.mp)
	p.WriteInt16(plr.maxMP)
	p.WriteInt16(plr.ap)
	p.WriteInt16(plr.sp)
	p.WriteInt32(plr.exp)
	p.WriteInt16(plr.fame)

	p.WriteInt32(plr.mapID)
	p.WriteByte(plr.mapPos)

	// p.WriteReversedInt64(time.Now().UnixMilli())
	// p.WriteInt32(0)
	//  p.WriteInt32(0)

	// Buddy list size
	p.WriteByte(plr.buddyListSize)

	// Money
	p.WriteInt32(plr.mesos)

	if plr.equipSlotSize == 0 {
		plr.equipSlotSize = 24
	}
	if plr.useSlotSize == 0 {
		plr.useSlotSize = 24
	}
	if plr.setupSlotSize == 0 {
		plr.setupSlotSize = 24
	}
	if plr.etcSlotSize == 0 {
		plr.etcSlotSize = 24
	}
	if plr.cashSlotSize == 0 {
		plr.cashSlotSize = 24
	}

	p.WriteByte(plr.equipSlotSize)
	p.WriteByte(plr.useSlotSize)
	p.WriteByte(plr.setupSlotSize)
	p.WriteByte(plr.etcSlotSize)
	p.WriteByte(plr.cashSlotSize)

	// Equipped (normal then cash)
	for _, it := range plr.equip {
		if it.slotID < 0 && !it.cash {
			p.WriteBytes(it.inventoryBytes())
		}
	}
	p.WriteByte(0)
	for _, it := range plr.equip {
		if it.slotID < 0 && it.cash {
			p.WriteBytes(it.inventoryBytes())
		}
	}
	p.WriteByte(0)

	// Inventory tabs
	writeInv := func(items []item) {
		cp := make([]item, 0, len(items))
		for _, it := range items {
			if it.slotID > 0 {
				cp = append(cp, it)
			}
		}
		sort.Slice(cp, func(i, j int) bool { return cp[i].slotID < cp[j].slotID })
		for _, it := range cp {
			p.WriteBytes(it.inventoryBytes())
		}
		p.WriteByte(0)
	}
	writeInv(plr.equip)
	writeInv(plr.use)
	writeInv(plr.setUp)
	writeInv(plr.etc)
	writeInv(plr.cash)

	// p.WriteInt16(0) // Minigames?
	// p.WriteInt16(0) // Rings?
	// p.WriteInt16(0) // Rocks?

	p.WriteUnicodeString(accountName)

	/*
		Wishlist example
		p.WriteByte(2)
		p.WriteInt32(item-id-1)
		p.WriteInt32(item-id-2)
	*/
	p.WriteByte(0) // Wishlist

	// 120 bytes padding
	for i := 0; i < 120; i++ {
		p.WriteInt64(0)
	}

	p.WriteByte(0)

	// Stock states block: send every commodity's current state
	comms := nx.GetCommodities()
	p.WriteInt16(int16(len(comms)))
	for sn, c := range comms {
		p.WriteInt32(sn)
		p.WriteInt32(c.StockState)
	}

	return p
}

func packetCashShopUpdateAmounts(nxCredit, maplePoints int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSUpdateAmounts)
	p.WriteInt32(nxCredit)
	p.WriteInt32(maplePoints)
	return p
}

func packetCashShopShowBoughtItem(charID int32, cashItemSNHash int64, itemID int32, count int16, itemName string) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	// Cash unique ID (hash-like), then player id
	p.WriteInt64(cashItemSNHash)
	p.WriteInt32(charID)

	// 4x 0x01 bytes (unknown flags in legacy implementations)
	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}

	p.WriteInt32(itemID)

	// 4x 0x01 bytes again
	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}

	p.WriteInt16(count)
	p.WriteString(itemName)
	p.WriteInt64(0) // expiration: 0 for non-expiring
	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}
	return p
}

// packetCashShopShowBoughtQuestItem mirrors "showBoughtCSQuestItem".
func packetCashShopShowBoughtQuestItem(position byte, itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	p.WriteInt32(365) // sub-op code per reference
	p.WriteByte(0)
	p.WriteInt16(1)
	p.WriteByte(position)
	p.WriteByte(0)
	p.WriteInt32(itemID)
	return p
}

// packetCashShopShowCouponRedeemedItem mirrors "showCouponRedeemedItem".
func packetCashShopShowCouponRedeemedItem(itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	p.WriteInt16(0x3A)
	p.WriteInt32(0)
	p.WriteInt32(1)
	p.WriteInt16(1)
	p.WriteInt16(0x1A)
	p.WriteInt32(itemID)
	p.WriteInt32(0)
	return p
}

func packetCashShopSendCSItemInventory(slotType byte, it item) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	p.WriteByte(0x2F)
	// Java used: writeShort((byte)slot); write(slot)
	p.WriteInt16(int16(slotType))
	p.WriteByte(slotType)

	// Reuse existing inventory encoder (as used elsewhere in server packets)
	// Ensure 'it' is correctly populated (invID, slotID, amount, id, stats etc.)
	p.WriteBytes(it.inventoryBytes())
	return p
}

// packetCashShopWishList mirrors "sendWishList". Provide up to 10 SN values.
func packetCashShopWishList(sns []int32, update bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	if update {
		p.WriteByte(0x39)
	} else {
		p.WriteByte(0x33)
	}
	count := 10
	for i := 0; i < count; i++ {
		var v int32
		if i < len(sns) {
			v = sns[i]
		}
		p.WriteInt32(v)
	}
	return p
}

// packetCashShopWrongCoupon mirrors "wrongCouponCode".
func packetCashShopWrongCoupon() mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSAction)
	p.WriteByte(0x40)
	p.WriteByte(0x87)
	return p
}
