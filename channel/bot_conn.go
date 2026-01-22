package channel

import (
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

// botConn is a stub connection for bot players.
// It implements mnet.Client interface but does nothing for network operations.
// This allows bots to be treated identically to real players throughout the codebase.
type botConn struct {
	accountID  int32
	worldID    byte
	channelID  byte
	adminLevel int
}

// newBotConn creates a new bot connection stub.
func newBotConn(channelID byte) *botConn {
	return &botConn{
		accountID:  -1, // Bots have no real account
		worldID:    0,
		channelID:  channelID,
		adminLevel: 0,
	}
}

// MConn interface implementation
func (bc *botConn) String() string {
	return "bot-connection"
}

func (bc *botConn) Send(p mpacket.Packet) {
	// Bots don't send packets - no-op
}

func (bc *botConn) Cleanup() {
	// Nothing to cleanup for bots
}

func (bc *botConn) Close() error {
	// Nothing to close for bots
	return nil
}

// Client interface implementation
func (bc *botConn) GetLogedIn() bool {
	return true // Bots are always "logged in"
}

func (bc *botConn) SetLogedIn(bool) {
	// No-op for bots
}

func (bc *botConn) GetAccountID() int32 {
	return bc.accountID
}

func (bc *botConn) SetAccountID(id int32) {
	bc.accountID = id
}

func (bc *botConn) GetGender() byte {
	return 0 // Male
}

func (bc *botConn) SetGender(byte) {
	// No-op for bots
}

func (bc *botConn) GetWorldID() byte {
	return bc.worldID
}

func (bc *botConn) SetWorldID(id byte) {
	bc.worldID = id
}

func (bc *botConn) GetChannelID() byte {
	return bc.channelID
}

func (bc *botConn) SetChannelID(id byte) {
	bc.channelID = id
}

func (bc *botConn) GetAdminLevel() int {
	return bc.adminLevel
}

func (bc *botConn) SetAdminLevel(level int) {
	bc.adminLevel = level
}

func (bc *botConn) GetHWID() string {
	return "bot-hwid"
}

func (bc *botConn) SetHWID(string) {
	// No-op for bots
}

func (bc *botConn) GetCashShopStorage() interface{} {
	return nil
}

func (bc *botConn) SetCashShopStorage(interface{}) {
	// No-op for bots
}

// Verify at compile time that botConn implements mnet.Client
var _ mnet.Client = (*botConn)(nil)
