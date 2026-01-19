package anticheat

import (
"database/sql"
"fmt"
"sync"
"time"
)

// Simple in-memory violation tracking with rolling windows
type AntiCheat struct {
violations map[string][]time.Time // "accountID:type" -> timestamps
failedAuth map[string][]time.Time // "user:X" or "ip:X" or "hwid:X" -> timestamps
mu         sync.RWMutex
db         *sql.DB
}

func New(db *sql.DB) *AntiCheat {
ac := &AntiCheat{
violations: make(map[string][]time.Time),
failedAuth: make(map[string][]time.Time),
db:         db,
}
go ac.cleanup()
return ac
}

// Track a violation - returns true if threshold exceeded and player should be banned
func (ac *AntiCheat) Track(accountID int32, violationType string, threshold int, window time.Duration) bool {
ac.mu.Lock()
defer ac.mu.Unlock()

key := fmt.Sprintf("%d:%s", accountID, violationType)
now := time.Now()
cutoff := now.Add(-window)

// Filter out old violations and add new one
timestamps := []time.Time{now}
for _, t := range ac.violations[key] {
if t.After(cutoff) {
timestamps = append(timestamps, t)
}
}
ac.violations[key] = timestamps

return len(timestamps) >= threshold
}

// Track failed auth attempt - returns true if should ban
func (ac *AntiCheat) TrackFailedAuth(identifier string) bool {
ac.mu.Lock()
defer ac.mu.Unlock()

now := time.Now()
cutoff := now.Add(-30 * time.Minute)

timestamps := []time.Time{now}
for _, t := range ac.failedAuth[identifier] {
if t.After(cutoff) {
timestamps = append(timestamps, t)
}
}
ac.failedAuth[identifier] = timestamps

return len(timestamps) >= 10 // 10 attempts in 30min
}

// Clear auth attempts on successful login
func (ac *AntiCheat) ClearAuth(identifiers ...string) {
ac.mu.Lock()
defer ac.mu.Unlock()
for _, id := range identifiers {
delete(ac.failedAuth, id)
}
}

// Background cleanup every 5 minutes
func (ac *AntiCheat) cleanup() {
ticker := time.NewTicker(5 * time.Minute)
for range ticker.C {
ac.mu.Lock()
cutoff := time.Now().Add(-1 * time.Hour)

for k, timestamps := range ac.violations {
var keep []time.Time
for _, t := range timestamps {
if t.After(cutoff) {
keep = append(keep, t)
}
}
if len(keep) > 0 {
ac.violations[k] = keep
} else {
delete(ac.violations, k)
}
}

for k, timestamps := range ac.failedAuth {
var keep []time.Time
for _, t := range timestamps {
if t.After(cutoff) {
keep = append(keep, t)
}
}
if len(keep) > 0 {
ac.failedAuth[k] = keep
} else {
delete(ac.failedAuth, k)
}
}
ac.mu.Unlock()
}
}

// IssueBan creates a temporary ban (hours=0 means permanent)
func (ac *AntiCheat) IssueBan(accountID int32, hours int, reason string, ip, hwid string) error {
var banEnd interface{}
if hours > 0 {
banEnd = time.Now().Add(time.Duration(hours) * time.Hour)
}

// Insert ban record
_, err := ac.db.Exec(`INSERT INTO bans (accountID, reason, banEnd, ip, hwid) VALUES (?, ?, ?, ?, ?)`,
accountID, reason, banEnd, ip, hwid)
if err != nil {
return err
}

// Set isBanned flag
_, err = ac.db.Exec(`UPDATE accounts SET isBanned = 1 WHERE accountID = ?`, accountID)

// Track temp bans for escalation
if hours > 0 {
count, _ := ac.incrementTempBans(accountID)
// Auto-escalate to permanent after 3 temp bans
if count >= 3 {
ac.IssueBan(accountID, 0, "Escalated: 3+ temporary bans", ip, hwid)
}
}

return err
}

// IsBanned checks if account/IP/HWID is banned
func (ac *AntiCheat) IsBanned(accountID int32, ip, hwid string) (bool, string, error) {
var reason string
err := ac.db.QueryRow(`
SELECT reason FROM bans 
WHERE (accountID = ? OR ip = ? OR (hwid = ? AND hwid != '')) 
AND (banEnd IS NULL OR banEnd > NOW())
LIMIT 1`, accountID, ip, hwid).Scan(&reason)

if err == sql.ErrNoRows {
return false, "", nil
}
if err != nil {
return false, "", err
}
return true, reason, nil
}

// Unban removes all bans for an account
func (ac *AntiCheat) Unban(accountID int32) error {
_, err := ac.db.Exec(`DELETE FROM bans WHERE accountID = ?`, accountID)
if err != nil {
return err
}

_, err = ac.db.Exec(`UPDATE accounts SET isBanned = 0 WHERE accountID = ?`, accountID)
return err
}

// GetBanHistory returns recent ban records
func (ac *AntiCheat) GetBanHistory(accountID int32, limit int) ([]string, error) {
rows, err := ac.db.Query(`
SELECT reason, banEnd, createdAt FROM bans 
WHERE accountID = ? 
ORDER BY createdAt DESC LIMIT ?`, accountID, limit)
if err != nil {
return nil, err
}
defer rows.Close()

var history []string
for rows.Next() {
var reason string
var banEnd sql.NullTime
var createdAt time.Time
if err := rows.Scan(&reason, &banEnd, &createdAt); err != nil {
continue
}

durStr := "permanent"
if banEnd.Valid {
durStr = banEnd.Time.Format("2006-01-02 15:04")
}
history = append(history, fmt.Sprintf("%s: %s (until %s)", 
createdAt.Format("2006-01-02 15:04"), reason, durStr))
}
return history, nil
}

// Internal: track temp ban count for escalation
func (ac *AntiCheat) incrementTempBans(accountID int32) (int, error) {
_, err := ac.db.Exec(`
INSERT INTO ban_escalation (accountID, count) VALUES (?, 1)
ON DUPLICATE KEY UPDATE count = count + 1`, accountID)
if err != nil {
return 0, err
}

var count int
err = ac.db.QueryRow(`SELECT count FROM ban_escalation WHERE accountID = ?`, accountID).Scan(&count)
return count, err
}
