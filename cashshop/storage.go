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
	sn        int32         // Serial number from commodity (used as cash ID)
	purchased int64         // Unix timestamp of purchase
	item      channel.Item  // The actual item
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
			itemID, sn, slotNumber, amount,
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
		var creatorName sql.NullString
		var slotNumber int16
		var itemID int32
		var amount int16
		var flag int16
		var upgradeSlots byte
		var scrollLevel byte
		var str, dex, intt, luk, hp, mp int16
		var watk, matk, wdef, mdef int16
		var accuracy, avoid, hands, speed, jump int16
		var expireTime int64
		
		if err := rows.Scan(
			&itemID, &csItem.sn, &slotNumber, &amount,
			&flag, &upgradeSlots, &scrollLevel,
			&str, &dex, &intt, &luk,
			&hp, &mp, &watk, &matk,
			&wdef, &mdef, &accuracy, &avoid,
			&hands, &speed, &jump,
			&expireTime, &creatorName, &csItem.purchased,
		); err != nil {
			log.Println("Error scanning cash shop storage item:", err)
			continue
		}
		
		// Create the item using the new helper function
		var creator string
		if creatorName.Valid {
			creator = creatorName.String
		}
		
		item, ierr := channel.CreateItemFromDBValues(
			itemID, slotNumber, amount, flag, upgradeSlots, scrollLevel,
			str, dex, intt, luk, hp, mp, watk, matk, wdef, mdef,
			accuracy, avoid, hands, speed, jump, expireTime, creator,
		)
		if ierr != nil {
			log.Println("Error creating item from DB values:", ierr)
			continue
		}
		
		csItem.item = item

		if slotNumber <= 0 || slotNumber > int16(s.maxSlots) {
			continue
		}
		idx := int(slotNumber - 1)
		s.items[idx] = csItem
		if itemID != 0 {
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
		if csItem.item.GetID() == 0 || csItem.item.GetAmount() == 0 {
			continue
		}

		slotNumber := int16(i + 1)
		
		_, ierr := stmt.Exec(
			s.accountID, csItem.item.GetID(), csItem.sn, slotNumber, csItem.item.GetAmount(),
			csItem.item.GetFlag(), csItem.item.GetUpgradeSlots(), csItem.item.GetScrollLevel(),
			csItem.item.GetStr(), csItem.item.GetDex(), csItem.item.GetIntt(), csItem.item.GetLuk(),
			csItem.item.GetHP(), csItem.item.GetMP(), csItem.item.GetWatk(), csItem.item.GetMatk(),
			csItem.item.GetWdef(), csItem.item.GetMdef(), csItem.item.GetAccuracy(), csItem.item.GetAvoid(),
			csItem.item.GetHands(), csItem.item.GetSpeed(), csItem.item.GetJump(),
			csItem.item.GetExpireTime(), csItem.item.GetCreatorName(),
		)
		if ierr != nil {
			err = fmt.Errorf("failed inserting cash shop item %d (acct %d, slot %d): %w", csItem.item.GetID(), s.accountID, slotNumber, ierr)
			return
		}
	}

	if cerr := tx.Commit(); cerr != nil {
		err = fmt.Errorf("failed to commit cash shop storage save (acct %d): %w", s.accountID, cerr)
		return
	}

	return nil
}

// AddItem adds an item to cash shop storage and returns the array index (0-based) where it was added
func (s *CashShopStorage) AddItem(item channel.Item, sn int32) (int, bool) {
	for i := 0; i < int(s.maxSlots); i++ {
		if s.items[i].item.GetID() != 0 {
			continue
		}
		s.totalSlotsUsed++
		s.items[i] = CashShopItem{
			sn:   sn,
			item: item,
		}
		return i, true
	}
	return -1, false
}

// AddItemWithCashID adds an item to storage with a specific SN (used when moving from inventory back to locker)
func (s *CashShopStorage) AddItemWithCashID(item channel.Item, sn int32, cashID int64) (int, bool) {
	// cashID is actually the SN, so we just use sn parameter
	return s.AddItem(item, sn)
}

// RemoveAt removes an item at the given index
func (s *CashShopStorage) RemoveAt(idx int) (*CashShopItem, bool) {
	if idx < 0 || idx >= int(s.maxSlots) {
		return nil, false
	}
	if s.items[idx].item.GetID() == 0 {
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
		if s.items[i].item.GetID() != 0 {
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
	if s.items[idx].item.GetID() == 0 {
		return nil, false
	}
	return &s.items[idx], true
}
