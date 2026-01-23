# Bot Player System

This document describes the bot player system implemented for Valhalla.

## Overview

The bot player system allows the server to spawn and control non-player characters (bots) that appear as normal players to clients. **Phase 1 and Phase 2 are complete** - bots can now move around naturally!

## What Bots Can Do (Phase 1 & 2)

✅ **Exist in maps** - Bots spawn and are visible to real players
✅ **Have appearance** - Bots have a name, level, job, and basic stats
✅ **Spawn on ground** - Proper foothold calculation prevents floating
✅ **Walk left/right** - Random directional movement at 100 pixels/second
✅ **Jump** - Occasional jumping (33% chance while walking)
✅ **Change facing** - Stance updates based on direction
✅ **Respect boundaries** - Bounce off map edges
✅ **Broadcast movement** - All players see bot movement in real-time

## What Bots Cannot Do (Yet)

❌ **Attack monsters** - No combat (Phase 3)
❌ **Use skills** - No skill system yet (Phase 3)
❌ **Interact with NPCs** - No interactions (Phase 3+)
❌ **Use portals** - Stay in spawn map
❌ **Join parties/guilds** - Excluded from social features
❌ **Gain EXP/loot** - Don't participate in progression

## Configuration

Bot support is controlled via the `enableBots` flag in channel configuration files.

### Enable Bots

In your channel config file (e.g., `config_dev.toml`, `config_channel_1.toml`):

```toml
[channel]
# ... other settings ...
enableBots = true  # Set to true to enable bot spawning with movement
```

**Default:** `false` (bots disabled)

### Disable Bots

Set `enableBots = false` or omit the setting entirely. No bot-related code will execute.

## Current Implementation (Phase 1 & 2)

### Phase 1: Minimal Bots
- Bots spawn safely in maps
- Proper foothold positioning
- Visible to real players
- Clean shutdown

### Phase 2: Basic Movement ✨ NEW!
- **Timer-driven movement**: 10 updates per second (100ms ticks)
- **Walk behavior**: Walk for 3 seconds, pause for 2 seconds
- **Random direction**: Changes randomly or on wall collision
- **Jumping**: 33% chance while walking, 2 second cooldown
- **Map boundaries**: Automatically bounce off edges
- **Foothold tracking**: Position recalculated each move

## Technical Details

### Bot Identification

Bots are identified internally by:
- **Negative character IDs**: Bot IDs are -1, -2, -3, etc. (never collide with real players)
- **`isBot` flag**: `Player.isBot` is set to `true`
- **Stub connection**: Bots use a `botConn` stub that implements `mnet.Client` with no-ops
- **AI controller**: Each bot has a `botAI` instance managing behavior

### Bot Lifecycle

1. **Spawn**: Server calls `InitializeBots()` after world connection
   - Currently spawns 1 test bot in Henesys (map 100000000)
   - Calculates proper foothold position
2. **Movement**: Background goroutine updates all bots every 100ms
   - Bot AI decides actions (walk, stop, jump)
   - Position updated with foothold calculation
   - Movement broadcast to all players
3. **Cleanup**: `RemoveAllBots()` called during channel shutdown
   - Stops movement ticker
   - Removes bots from maps and player collections

### Movement System (Phase 2)

**Bot AI Controller** (`botAI` struct):
- Tracks movement state (direction, speed, jump state)
- Decides when to walk, stop, or jump
- Calculates velocity and new position
- Generates movement packets
- Respects map boundaries

**Server Integration**:
```go
// Movement ticker (100ms = 10 updates/second)
server.botMovementTimer = time.NewTicker(100 * time.Millisecond)

// Background goroutine
go server.runBotMovementLoop()

// Each tick:
for _, bot := range server.bots {
    bot.botAI.Update()           // Decide action
    bot.botAI.PerformMovement()  // Execute & broadcast
}
```

**Movement Packet**:
- Uses standard `SendChannelPlayerMovement` opcode
- Same format as real player movement
- Includes position, velocity, foothold, stance
- Broadcast via `inst.movePlayer()`

### Database Persistence

Bots are **ephemeral** and never written to the database:
- `saver.go` checks `p.isBot` and skips persistence
- Bots disappear on server restart (by design)
- No DB schema changes required

### Code Organization

```
channel/
├── bot_conn.go      # Stub connection for bots (mnet.Client impl)
├── bot_ai.go        # Movement AI and behavior logic (Phase 2)
├── player.go        # newBotPlayer() factory, isBot flag, botAI field
├── server.go        # Bot management (spawn, remove, init, movement loop)
└── bot_test.go      # Basic bot tests
```

## Default Bot Configuration

The initial test bot has the following properties:

| Property | Value |
|----------|-------|
| Name | "TestBot" |
| Level | 1 |
| Job | Beginner (0) |
| Map | Henesys (100000000) |
| Stats | 4 STR/DEX/INT/LUK, 50 HP, 5 MP |
| Appearance | Male, default face (20000), default hair (30000) |
| Movement | Enabled (walk/jump/turn) |
| Speed | 100 pixels/second walk, -150 jump velocity |

## Troubleshooting

### Bots don't appear

1. Check `enableBots = true` in your config file
2. Check server logs for "Bots are enabled" message
3. Check server logs for "Bot movement system started"
4. Verify Henesys map (100000000) exists in NX data

### Bots appear but don't move

1. Check server logs for "Bot [name] movement enabled" message
2. Verify map boundaries are calculated (vrLimit)
3. Check movement ticker is running (should see in logs)

### Bots are floating in air

This was fixed in commit ed02893. If still happening:
1. Update to latest code
2. Verify foothold histogram is passed to `newBotPlayer()`
3. Check map has valid footholds in NX data

### Server crashes on startup

1. Ensure NX data is loaded (`Data.nx`)
2. Check for map ID 100000000 in NX data
3. Review server logs for specific error

## Safety Features

- **Hard-coded limit**: Currently only 1 bot spawns (easy to expand later)
- **Config toggle**: Single flag disables all bot logic
- **No DB writes**: Bots can't corrupt database
- **Negative IDs**: Impossible to collide with real players
- **Stub connection**: Bots can't send malformed packets
- **Map boundaries**: Can't walk off edges or through walls
- **Clean shutdown**: Movement ticker stops properly

## Performance

- Movement ticker: 100ms (10 updates/second)
- Minimal CPU usage per bot
- Packet generation only when moving
- Efficient map boundary checks
- Scales well with multiple bots

## Development Notes

When extending bot functionality:

1. Always check `server.enableBots` before bot operations
2. Exclude bots from systems that assume real players (EXP, loot, parties)
3. Test with `enableBots = false` to ensure no regressions
4. Log bot spawn/despawn events for debugging
5. Use `plr.isBot` checks where bot behavior should differ
6. Movement logic in `botAI` struct for clean separation

## Examples

### Spawning a Custom Bot

```go
// In server.go or similar
if server.enableBots {
    err := server.SpawnBot("MyBot", mapID, portalID)
    if err != nil {
        log.Printf("Failed to spawn bot: %v", err)
    }
    // Movement is automatically enabled
}
```

### Checking if Player is Bot

```go
if plr.isBot {
    // Skip EXP gain, loot, etc.
    return
}
```

### Accessing Bot AI

```go
if plr.isBot && plr.botAI != nil {
    state := plr.botAI.GetMovementState()
    log.Printf("Bot %s is %s", plr.Name, state)
}
```

### Disabling Bot Movement

```go
if bot.botAI != nil {
    bot.botAI.DisableMovement()
}
```

## Future Phases

### Phase 3: Combat-Capable Bots (PvE Only)
- Auto-attack nearby monsters
- Use limited skill set
- Target selection (nearest/lowest HP)
- Fixed/capped damage
- No EXP or drop generation

### Phase 4: Event-Compatible Bots
- Optional GM event participation
- Basic event logic (movement, positioning)
- Never win rewards by default
- Removable mid-event

### Phase 5: Advanced Behavior
- Smarter movement and targeting
- Configurable bot profiles (job, gear, behavior)
- Scripted behavior per map
- Dynamic spawn/despawn rules

## Testing

Run bot-specific tests:

```bash
go test ./channel -run TestBot
```

All tests should pass. Note that `TestBotPlayerCreation` is skipped as it requires NX data.

## License

Bot system is part of Valhalla and follows the same license.

---

**Updated**: 2026-01-22 - Phase 2 complete with movement system!
