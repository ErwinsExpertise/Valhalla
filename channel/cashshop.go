package channel

import (
	"sort"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
)

func packetCashShopSet(plr *player, accountName string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSetCashShop)

	// Flags: Stats|Money|Equips|Consume|Install|Etc|Pet|Skills
	p.WriteInt16(0x00FF)

	// Stats (exactly like login)
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

	// No extra fields here. Buddy size comes next (like login).
	p.WriteByte(plr.buddyListSize)

	// Inventory
	p.WriteInt32(plr.mesos)

	// Equipped (negative slots), non-cash then cash
	for _, it := range sortEquipped(plr.equip, false) {
		p.WriteBytes(it.inventoryBytes())
	}
	p.WriteByte(0)

	for _, it := range sortEquipped(plr.equip, true) {
		p.WriteBytes(it.inventoryBytes())
	}
	p.WriteByte(0)

	// Per-tab writer: MaxSlots byte, then items (slotID > 0 sorted), then 0
	writeTab := func(capacity byte, items []item) {
		if capacity == 0 {
			// Defensive: most clients expect non-zero; use a sane default if DB holds 0
			capacity = 24
		}
		p.WriteByte(capacity)
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

	// Equip (inventory), Use, Setup, Etc, Cash
	writeTab(plr.equipSlotSize, plr.equip)
	writeTab(plr.useSlotSize, plr.use)
	writeTab(plr.setupSlotSize, plr.setUp)
	writeTab(plr.etcSlotSize, plr.etc)
	writeTab(plr.cashSlotSize, plr.cash)

	// Skills: count + (skillID, level) — no cooldown list in CS packet
	p.WriteInt16(int16(len(plr.skills)))
	for _, sk := range plr.skills {
		p.WriteInt32(sk.ID)
		p.WriteInt32(int32(sk.Level))
	}

	// Footer
	p.WriteBool(true)
	p.WriteString(accountName)

	// Wishlist (0)
	p.WriteInt16(0)

	// BEST items + stock states + trailing longs
	writeCSTopItemsAndStock(&p, nil, nil)

	return p
}

func sortEquipped(items []item, cash bool) []item {
	cp := make([]item, 0, len(items))
	for _, it := range items {
		if it.slotID < 0 && it.cash == cash {
			cp = append(cp, it)
		}
	}
	sort.Slice(cp, func(i, j int) bool {
		a, b := cp[i].slotID, cp[j].slotID
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		return a < b
	})
	return cp
}

// writeCSTopItemsAndStock writes:
//   - BEST items: for categories 1..8 and genders 0..1, 5 top items each.
//     Layout per entry: int category, int gender, int sn
//   - Custom stock states: ushort count + [int sn, int stockState]*
//   - 9 trailing long(0)
type bestItemKey struct {
	category byte
	gender   byte
	index    byte
}

// bestItems may provide SNs mapped by (category, gender, index 0..4). If nil, writes zeros.
// stockStates is an optional slice of (sn, state). If nil, writes zero count.
func writeCSTopItemsAndStock(p *mpacket.Packet, bestItems map[bestItemKey]int32, stockStates []struct {
	sn    int32
	state int32
}) {
	// BEST items
	for cat := byte(1); cat <= 8; cat++ {
		for gen := byte(0); gen <= 1; gen++ {
			for idx := byte(0); idx < 5; idx++ {
				p.WriteInt32(int32(cat)) // category
				p.WriteInt32(int32(gen)) // gender
				sn := int32(0)
				if bestItems != nil {
					if v, ok := bestItems[bestItemKey{category: cat, gender: gen, index: idx}]; ok {
						sn = v
					}
				}
				p.WriteInt32(sn) // SN (0 if none)
			}
		}
	}

	// Custom stock states
	if stockStates == nil {
		p.WriteInt16(0) // count
	} else {
		p.WriteInt16(int16(len(stockStates)))
		for _, s := range stockStates {
			p.WriteInt32(s.sn)
			p.WriteInt32(s.state)
		}
	}

	// 9 trailing longs
	for i := 0; i < 9; i++ {
		p.WriteInt64(0)
	}
}

// packetCashShopUpdateAmounts mirrors "sendCash" (CS_CASH): writes the player's credit balances.
// nxCredit: e.g., PayPal/PayByCash NX. maplePoints: Maple Points.
func packetCashShopUpdateAmounts(nxCredit, maplePoints int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode16(opcode.SendChannelCSUpdateAmounts)
	p.WriteInt32(nxCredit)
	p.WriteInt32(maplePoints)
	return p
}

// packetCashShopShowBoughtItem mirrors "showBoughtCSItem" (CS_OPERATION sub-op).
// This is a best-effort structure aligned to your Java reference.
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

// packetCashShopSendCSItemInventory mirrors "sendCSItemInventory" with a simplified payload.
// Note: protocol differences between versions exist; this sub-op 0x2F + an item payload
// is sufficient for many v62-like clients to select the tab and show the item in normal inventory.
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
