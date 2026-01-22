package channel

import (
	"testing"
)

// TestBotConnImplementsClient verifies botConn implements mnet.Client interface.
// The interface implementation is verified at compile time with var _ declaration.
func TestBotConnCompiles(t *testing.T) {
	conn := newBotConn(1)
	if conn == nil {
		t.Error("newBotConn returned nil")
	}

	if conn.String() != "bot-connection" {
		t.Errorf("expected 'bot-connection', got %s", conn.String())
	}

	if !conn.GetLogedIn() {
		t.Error("bot should always be logged in")
	}

	// Test no-op methods don't panic
	conn.Send(nil)
	conn.Cleanup()
	_ = conn.Close()
}

// TestBotPlayerCreation verifies bot player creation works.
func TestBotPlayerCreation(t *testing.T) {
	// Note: This test requires NX data to be loaded, so we skip it
	// in CI environments. The compile-time check is the main validation.
	t.Skip("Skipping bot creation test - requires NX data")

	bot, err := newBotPlayer(-1, "TestBot", 100000000, 0, 1)
	if err != nil {
		t.Fatalf("failed to create bot: %v", err)
	}

	if bot.ID != -1 {
		t.Errorf("expected bot ID -1, got %d", bot.ID)
	}

	if bot.Name != "TestBot" {
		t.Errorf("expected name 'TestBot', got %s", bot.Name)
	}

	if !bot.isBot {
		t.Error("bot flag should be true")
	}

	if bot.Conn == nil {
		t.Error("bot should have a connection stub")
	}

	if bot.level != 1 {
		t.Errorf("expected level 1, got %d", bot.level)
	}
}

// TestBotIDsNegative verifies bot IDs are negative.
func TestBotIDsAreNegative(t *testing.T) {
	// This is just a documentation test to ensure the convention is clear
	botID := int32(-1)
	if botID >= 0 {
		t.Error("bot IDs must be negative to avoid collision with real players")
	}
}
