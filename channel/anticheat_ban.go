package channel

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Hucaru/Valhalla/common"
)

// BanType represents the type of ban
type BanType string

const (
	BanTypeTemporary BanType = "temporary"
	BanTypePermanent BanType = "permanent"
)

// BanTarget represents what entity is being banned
type BanTarget string

const (
	BanTargetCharacter BanTarget = "character"
	BanTargetAccount   BanTarget = "account"
	BanTargetIP        BanTarget = "ip"
)

// Ban represents a ban record
type Ban struct {
	ID           int64
	AccountID    *int32
	CharacterID  *int32
	IPAddress    *string
	BanType      BanType
	BanTarget    BanTarget
	Reason       string
	IssuedBy     string
	IssuedByGM   bool
	IsActive     bool
	BanStartTime time.Time
	BanEndTime   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BanService handles ban operations
type BanService struct {
	config AntiCheatConfig
}

// NewBanService creates a new ban service
func NewBanService(config AntiCheatConfig) *BanService {
	return &BanService{
		config: config,
	}
}

// IsAccountBanned checks if an account is currently banned
func (bs *BanService) IsAccountBanned(accountID int32) (bool, *Ban, error) {
	return bs.checkBan(BanTargetAccount, &accountID, nil, nil)
}

// IsCharacterBanned checks if a character is currently banned
func (bs *BanService) IsCharacterBanned(characterID int32) (bool, *Ban, error) {
	return bs.checkBan(BanTargetCharacter, nil, &characterID, nil)
}

// IsIPBanned checks if an IP address is currently banned
func (bs *BanService) IsIPBanned(ipAddress string) (bool, *Ban, error) {
	return bs.checkBan(BanTargetIP, nil, nil, &ipAddress)
}

// checkBan is a helper to check ban status
func (bs *BanService) checkBan(target BanTarget, accountID, characterID *int32, ipAddress *string) (bool, *Ban, error) {
	var query string
	var args []interface{}

	switch target {
	case BanTargetAccount:
		query = `SELECT id, accountID, characterID, ipAddress, banType, banTarget, reason, issuedBy, 
				issuedByGM, isActive, banStartTime, banEndTime, createdAt, updatedAt 
				FROM bans WHERE accountID = ? AND isActive = 1 AND (banType = 'permanent' OR banEndTime > NOW())`
		args = []interface{}{*accountID}
	case BanTargetCharacter:
		query = `SELECT id, accountID, characterID, ipAddress, banType, banTarget, reason, issuedBy, 
				issuedByGM, isActive, banStartTime, banEndTime, createdAt, updatedAt 
				FROM bans WHERE characterID = ? AND isActive = 1 AND (banType = 'permanent' OR banEndTime > NOW())`
		args = []interface{}{*characterID}
	case BanTargetIP:
		query = `SELECT id, accountID, characterID, ipAddress, banType, banTarget, reason, issuedBy, 
				issuedByGM, isActive, banStartTime, banEndTime, createdAt, updatedAt 
				FROM bans WHERE ipAddress = ? AND isActive = 1 AND (banType = 'permanent' OR banEndTime > NOW())`
		args = []interface{}{*ipAddress}
	default:
		return false, nil, fmt.Errorf("invalid ban target: %s", target)
	}

	var ban Ban
	var accountIDVal, characterIDVal sql.NullInt32
	var ipAddressVal, issuedByVal sql.NullString
	var banEndTimeVal sql.NullTime

	err := common.DB.QueryRow(query, args...).Scan(
		&ban.ID, &accountIDVal, &characterIDVal, &ipAddressVal,
		&ban.BanType, &ban.BanTarget, &ban.Reason, &issuedByVal,
		&ban.IssuedByGM, &ban.IsActive, &ban.BanStartTime, &banEndTimeVal,
		&ban.CreatedAt, &ban.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, fmt.Errorf("error checking ban: %w", err)
	}

	// Convert nullable fields
	if accountIDVal.Valid {
		val := int32(accountIDVal.Int32)
		ban.AccountID = &val
	}
	if characterIDVal.Valid {
		val := int32(characterIDVal.Int32)
		ban.CharacterID = &val
	}
	if ipAddressVal.Valid {
		ban.IPAddress = &ipAddressVal.String
	}
	if issuedByVal.Valid {
		ban.IssuedBy = issuedByVal.String
	}
	if banEndTimeVal.Valid {
		ban.BanEndTime = &banEndTimeVal.Time
	}

	return true, &ban, nil
}

// IssueBan issues a new ban
func (bs *BanService) IssueBan(accountID *int32, characterID *int32, ipAddress *string, 
	banType BanType, banTarget BanTarget, reason string, issuedBy string, issuedByGM bool) error {
	
	var banEndTime *time.Time
	if banType == BanTypeTemporary {
		endTime := time.Now().Add(bs.config.DefaultTempBanDuration)
		banEndTime = &endTime
	}

	// Insert ban record
	query := `INSERT INTO bans (accountID, characterID, ipAddress, banType, banTarget, reason, 
			issuedBy, issuedByGM, isActive, banStartTime, banEndTime) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, NOW(), ?)`
	
	_, err := common.DB.Exec(query, accountID, characterID, ipAddress, banType, banTarget, 
		reason, issuedBy, issuedByGM, banEndTime)
	if err != nil {
		return fmt.Errorf("failed to issue ban: %w", err)
	}

	// Update escalation counter if applicable
	if accountID != nil && banType == BanTypeTemporary {
		if issuedByGM && !bs.config.GMBansIncrementCounter {
			// GM bans don't count towards escalation
			log.Printf("GM ban issued for account %d, not incrementing escalation counter", *accountID)
		} else {
			err = bs.incrementTempBanCount(*accountID)
			if err != nil {
				log.Printf("Error incrementing temp ban count: %v", err)
			}

			// Check if escalation to permanent ban is needed
			err = bs.checkEscalation(*accountID)
			if err != nil {
				log.Printf("Error checking escalation: %v", err)
			}
		}
	}

	log.Printf("Ban issued: type=%s target=%s accountID=%v characterID=%v ipAddress=%v reason=%s issuedBy=%s",
		banType, banTarget, accountID, characterID, ipAddress, reason, issuedBy)

	return nil
}

// incrementTempBanCount increments the temporary ban count for an account
func (bs *BanService) incrementTempBanCount(accountID int32) error {
	query := `INSERT INTO ban_escalation (accountID, tempBanCount, lastTempBanTime) 
			VALUES (?, 1, NOW()) 
			ON DUPLICATE KEY UPDATE 
			tempBanCount = tempBanCount + 1, 
			lastTempBanTime = NOW()`
	
	_, err := common.DB.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to increment temp ban count: %w", err)
	}
	return nil
}

// checkEscalation checks if an account should be permanently banned based on escalation rules
func (bs *BanService) checkEscalation(accountID int32) error {
	var tempBanCount int
	var permanentBanIssued bool

	query := `SELECT tempBanCount, permanentBanIssued FROM ban_escalation WHERE accountID = ?`
	err := common.DB.QueryRow(query, accountID).Scan(&tempBanCount, &permanentBanIssued)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No escalation record yet
		}
		return fmt.Errorf("failed to check escalation: %w", err)
	}

	// If already permanently banned, don't escalate again
	if permanentBanIssued {
		return nil
	}

	// Check if threshold reached
	if tempBanCount >= bs.config.TempBansBeforePermanent {
		log.Printf("Account %d has reached escalation threshold (%d temp bans), issuing permanent ban",
			accountID, tempBanCount)

		// Issue permanent ban
		err = bs.IssueBan(&accountID, nil, nil, BanTypePermanent, BanTargetAccount,
			fmt.Sprintf("Automatic escalation after %d temporary bans", tempBanCount),
			"SYSTEM", false)
		if err != nil {
			return fmt.Errorf("failed to issue escalation ban: %w", err)
		}

		// Mark permanent ban issued
		_, err = common.DB.Exec(`UPDATE ban_escalation SET permanentBanIssued = 1 WHERE accountID = ?`, accountID)
		if err != nil {
			log.Printf("Failed to mark permanent ban issued: %v", err)
		}
	}

	return nil
}

// Unban removes an active ban
func (bs *BanService) Unban(banID int64, unbannedBy string) error {
	query := `UPDATE bans SET isActive = 0, updatedAt = NOW() WHERE id = ?`
	result, err := common.DB.Exec(query, banID)
	if err != nil {
		return fmt.Errorf("failed to unban: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no ban found with ID %d", banID)
	}

	log.Printf("Ban %d removed by %s", banID, unbannedBy)
	return nil
}

// GetBanHistory retrieves ban history for an account or character
func (bs *BanService) GetBanHistory(accountID *int32, characterID *int32, limit int) ([]Ban, error) {
	var query string
	var args []interface{}

	if accountID != nil {
		query = `SELECT id, accountID, characterID, ipAddress, banType, banTarget, reason, issuedBy, 
				issuedByGM, isActive, banStartTime, banEndTime, createdAt, updatedAt 
				FROM bans WHERE accountID = ? ORDER BY createdAt DESC LIMIT ?`
		args = []interface{}{*accountID, limit}
	} else if characterID != nil {
		query = `SELECT id, accountID, characterID, ipAddress, banType, banTarget, reason, issuedBy, 
				issuedByGM, isActive, banStartTime, banEndTime, createdAt, updatedAt 
				FROM bans WHERE characterID = ? ORDER BY createdAt DESC LIMIT ?`
		args = []interface{}{*characterID, limit}
	} else {
		return nil, fmt.Errorf("must provide either accountID or characterID")
	}

	rows, err := common.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get ban history: %w", err)
	}
	defer rows.Close()

	var bans []Ban
	for rows.Next() {
		var ban Ban
		var accountIDVal, characterIDVal sql.NullInt32
		var ipAddressVal, issuedByVal sql.NullString
		var banEndTimeVal sql.NullTime

		err := rows.Scan(&ban.ID, &accountIDVal, &characterIDVal, &ipAddressVal,
			&ban.BanType, &ban.BanTarget, &ban.Reason, &issuedByVal,
			&ban.IssuedByGM, &ban.IsActive, &ban.BanStartTime, &banEndTimeVal,
			&ban.CreatedAt, &ban.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning ban record: %v", err)
			continue
		}

		// Convert nullable fields
		if accountIDVal.Valid {
			val := int32(accountIDVal.Int32)
			ban.AccountID = &val
		}
		if characterIDVal.Valid {
			val := int32(characterIDVal.Int32)
			ban.CharacterID = &val
		}
		if ipAddressVal.Valid {
			ban.IPAddress = &ipAddressVal.String
		}
		if issuedByVal.Valid {
			ban.IssuedBy = issuedByVal.String
		}
		if banEndTimeVal.Valid {
			ban.BanEndTime = &banEndTimeVal.Time
		}

		bans = append(bans, ban)
	}

	return bans, nil
}

// ExpireOldBans marks old temporary bans as inactive
func (bs *BanService) ExpireOldBans() error {
	query := `UPDATE bans SET isActive = 0 
			WHERE isActive = 1 AND banType = 'temporary' AND banEndTime < NOW()`
	
	result, err := common.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to expire old bans: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows > 0 {
		log.Printf("Expired %d old temporary bans", rows)
	}

	return nil
}
