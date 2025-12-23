package cashshop

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Hucaru/Valhalla/channel"
	"github.com/Hucaru/Valhalla/common"
)

// Cash shop storage capacity bounds
const (
	cashShopStorageMinSlots byte = 50
	cashShopStorageMaxSlots byte = 255
)

// CashShopStorage represents account-wide cash shop storage
type CashShopStorage struct {
	accountID      int32
	maxSlots       byte
	totalSlotsUsed byte
	items          []CashShopItem
}

// CashShopItem represents an item in the cash shop storage
type CashShopItem struct {
	dbID      int64
	itemID    int32
	sn        int32 // serial number from commodity
	slotID    int16
	amount    int16
	purchased int64 // unix timestamp
	
	// Item properties
	flag         int16
	upgradeSlots byte
	scrollLevel  byte
	str          int16
	dex          int16
	intt         int16
	luk          int16
	hp           int16
	mp           int16
	watk         int16
	matk         int16
	wdef         int16
	mdef         int16
	accuracy     int16
	avoid        int16
	hands        int16
	speed        int16
	jump         int16
	expireTime   int64
	creatorName  string
	invID        byte
}

func clampByte(v, min, max byte) byte {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// NewCashShopStorage creates a new cash shop storage instance
func NewCashShopStorage(accountID int32) *CashShopStorage {
	return &CashShopStorage{
		accountID: accountID,
		maxSlots:  cashShopStorageMinSlots,
		items:     make([]CashShopItem, cashShopStorageMinSlots),
	}
}

func (s *CashShopStorage) ensureCapacity() {
	if s.items == nil || byte(len(s.items)) != s.maxSlots {
		newArr := make([]CashShopItem, s.maxSlots)
		if s.items != nil {
			copy(newArr, s.items)
		}
		s.items = newArr
	}
}

// Load loads cash shop storage from database
func (s *CashShopStorage) Load() error {
	var slots sql.NullInt64
	if err := common.DB.QueryRow(
		"SELECT slots FROM account_cashshop_storage WHERE accountID=?",
		s.accountID,
	).Scan(&slots); err != nil {
		if err == sql.ErrNoRows {
			// Initialize storage for this account
			if _, ierr := common.DB.Exec(
				"INSERT INTO account_cashshop_storage(accountID, slots) VALUES(?,?)",
				s.accountID, cashShopStorageMinSlots,
			); ierr != nil {
				return fmt.Errorf("couldn't initialize cash shop storage for account %d: %w", s.accountID, ierr)
			}
			s.maxSlots = cashShopStorageMinSlots
			s.ensureCapacity()
			s.totalSlotsUsed = 0
			return nil
		}
		return fmt.Errorf("failed to load cash shop storage header for account %d: %w", s.accountID, err)
	}

	if slots.Valid {
		s.maxSlots = clampByte(byte(slots.Int64), cashShopStorageMinSlots, cashShopStorageMaxSlots)
	} else {
		s.maxSlots = cashShopStorageMinSlots
	}

	s.ensureCapacity()
	s.totalSlotsUsed = 0

	rows, qerr := common.DB.Query(`
		SELECT 
			id, itemID, sn, slotNumber, amount,
			flag, upgradeSlots, level, str, dex, intt, luk, hp, mp,
			watk, matk, wdef, mdef, accuracy, avoid, hands, speed, jump,
			expireTime, creatorName, UNIX_TIMESTAMP(purchaseDate)
		FROM account_cashshop_storage_items
		WHERE accountID=?
		ORDER BY slotNumber ASC`, s.accountID)
	if qerr != nil {
		return fmt.Errorf("failed to load cash shop storage items for account %d: %w", s.accountID, qerr)
	}
	defer rows.Close()

	for rows.Next() {
		var csItem CashShopItem
		var creator sql.NullString
		var slotNumber int16
		
		if err := rows.Scan(
			&csItem.dbID, &csItem.itemID, &csItem.sn, &slotNumber, &csItem.amount,
			&csItem.flag, &csItem.upgradeSlots, &csItem.scrollLevel,
			&csItem.str, &csItem.dex, &csItem.intt, &csItem.luk,
			&csItem.hp, &csItem.mp, &csItem.watk, &csItem.matk,
			&csItem.wdef, &csItem.mdef, &csItem.accuracy, &csItem.avoid,
			&csItem.hands, &csItem.speed, &csItem.jump,
			&csItem.expireTime, &creator, &csItem.purchased,
		); err != nil {
			log.Println("Error scanning cash shop storage item:", err)
			continue
		}
		
		if creator.Valid {
			csItem.creatorName = creator.String
		}
		
		csItem.slotID = slotNumber

		if slotNumber <= 0 || int(slotNumber) > len(s.items) {
			continue
		}
		idx := int(slotNumber - 1)
		s.items[idx] = csItem
		if csItem.itemID != 0 {
			s.totalSlotsUsed++
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error while reading cash shop storage items for account %d: %w", s.accountID, err)
	}

	return nil
}

// Save saves cash shop storage to database
func (s *CashShopStorage) Save() (err error) {
	tx, terr := common.DB.Begin()
	if terr != nil {
		return fmt.Errorf("couldn't open transaction to save cash shop storage (acct %d): %w", s.accountID, terr)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, uerr := tx.Exec(
		"UPDATE account_cashshop_storage SET slots=? WHERE accountID=?",
		s.maxSlots, s.accountID,
	); uerr != nil {
		err = fmt.Errorf("failed to update cash shop storage header (acct %d): %w", s.accountID, uerr)
		return
	}

	if _, derr := tx.Exec(
		"DELETE FROM account_cashshop_storage_items WHERE accountID=?",
		s.accountID,
	); derr != nil {
		err = fmt.Errorf("failed to clear cash shop storage items (acct %d): %w", s.accountID, derr)
		return
	}

	const ins = `
		INSERT INTO account_cashshop_storage_items(
			accountID, itemID, sn, slotNumber, amount, flag, upgradeSlots, level,
			str, dex, intt, luk, hp, mp, watk, matk, wdef, mdef, accuracy, avoid, hands,
			speed, jump, expireTime, creatorName
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	stmt, perr := tx.Prepare(ins)
	if perr != nil {
		err = fmt.Errorf("failed to prepare item insert (acct %d): %w", s.accountID, perr)
		return
	}
	defer stmt.Close()

	for i := range s.items {
		csItem := s.items[i]
		if csItem.itemID == 0 || csItem.amount == 0 {
			continue
		}

		slotNumber := int16(i + 1)
		if _, ierr := stmt.Exec(
			s.accountID, csItem.itemID, csItem.sn, slotNumber, csItem.amount,
			csItem.flag, csItem.upgradeSlots, csItem.scrollLevel,
			csItem.str, csItem.dex, csItem.intt, csItem.luk,
			csItem.hp, csItem.mp, csItem.watk, csItem.matk,
			csItem.wdef, csItem.mdef, csItem.accuracy, csItem.avoid,
			csItem.hands, csItem.speed, csItem.jump,
			csItem.expireTime, csItem.creatorName,
		); ierr != nil {
			err = fmt.Errorf("failed inserting cash shop item %d (acct %d, slot %d): %w", csItem.itemID, s.accountID, slotNumber, ierr)
			return
		}
	}

	if cerr := tx.Commit(); cerr != nil {
		err = fmt.Errorf("failed to commit cash shop storage save (acct %d): %w", s.accountID, cerr)
		return
	}

	return nil
}

// AddItem adds an item to cash shop storage and returns the slot it was added to
func (s *CashShopStorage) AddItem(item channel.Item, sn int32) (int, bool) {
	data := item.ExportData()
	for i := 0; i < int(s.maxSlots); i++ {
		if s.items[i].itemID != 0 {
			continue
		}
		s.totalSlotsUsed++
		slotID := int16(i + 1)
		s.items[i] = CashShopItem{
			itemID:       data.ItemID,
			sn:           sn,
			slotID:       slotID,
			amount:       data.Amount,
			flag:         data.Flag,
			upgradeSlots: data.UpgradeSlots,
			scrollLevel:  data.ScrollLevel,
			str:          data.Str,
			dex:          data.Dex,
			intt:         data.Intt,
			luk:          data.Luk,
			hp:           data.HP,
			mp:           data.MP,
			watk:         data.Watk,
			matk:         data.Matk,
			wdef:         data.Wdef,
			mdef:         data.Mdef,
			accuracy:     data.Accuracy,
			avoid:        data.Avoid,
			hands:        data.Hands,
			speed:        data.Speed,
			jump:         data.Jump,
			expireTime:   data.ExpireTime,
			creatorName:  data.CreatorName,
			invID:        data.InvID,
		}
		return i, true
	}
	return -1, false
}

// RemoveAt removes an item at the given index
func (s *CashShopStorage) RemoveAt(idx int) (*CashShopItem, bool) {
	if idx < 0 || idx >= int(s.maxSlots) {
		return nil, false
	}
	if s.items[idx].itemID == 0 {
		return nil, false
	}

	item := s.items[idx]
	s.items[idx] = CashShopItem{}
	if s.totalSlotsUsed > 0 {
		s.totalSlotsUsed--
	}
	return &item, true
}

// GetAllItems returns all items in storage
func (s *CashShopStorage) GetAllItems() []CashShopItem {
	out := make([]CashShopItem, 0, s.totalSlotsUsed)
	for i := range s.items {
		if s.items[i].itemID != 0 {
			out = append(out, s.items[i])
		}
	}
	return out
}

// SlotsAvailable returns true if there are slots available
func (s *CashShopStorage) SlotsAvailable() bool {
	return s.totalSlotsUsed < s.maxSlots
}

// GetItemBySlot returns the item at the given slot (1-indexed)
func (s *CashShopStorage) GetItemBySlot(slot int16) (*CashShopItem, bool) {
	if slot <= 0 || int(slot) > len(s.items) {
		return nil, false
	}
	idx := int(slot - 1)
	if s.items[idx].itemID == 0 {
		return nil, false
	}
	return &s.items[idx], true
}

// ToItem converts a CashShopItem back to a channel.Item for giving to player
func (csItem *CashShopItem) ToItem() (channel.Item, error) {
	// Create base item
	item, err := channel.CreateItemFromID(csItem.itemID, csItem.amount)
	if err != nil {
		return item, err
	}
	
	// Restore all the stored properties
	data := item.ExportData()
	data.Flag = csItem.flag
	data.UpgradeSlots = csItem.upgradeSlots
	data.ScrollLevel = csItem.scrollLevel
	data.Str = csItem.str
	data.Dex = csItem.dex
	data.Intt = csItem.intt
	data.Luk = csItem.luk
	data.HP = csItem.hp
	data.MP = csItem.mp
	data.Watk = csItem.watk
	data.Matk = csItem.matk
	data.Wdef = csItem.wdef
	data.Mdef = csItem.mdef
	data.Accuracy = csItem.accuracy
	data.Avoid = csItem.avoid
	data.Hands = csItem.hands
	data.Speed = csItem.speed
	data.Jump = csItem.jump
	data.ExpireTime = csItem.expireTime
	data.CreatorName = csItem.creatorName
	
	// Use the helper to rebuild item from data
	return channel.ItemFromExportedData(data)
}
