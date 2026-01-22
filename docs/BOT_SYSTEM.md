# Bot Player System

This document describes the bot player system implemented for Valhalla.

## Overview

The bot player system allows the server to spawn and control non-player characters (bots) that appear as normal players to clients. This is useful for:
- Solo play testing
- Filling maps during development
- GM events with simulated participants
- Demonstrations and screenshots

## Configuration

Bot support is controlled via the `enableBots` flag in channel configuration files.

### Enable Bots

In your channel config file (e.g., `config_dev.toml`, `config_channel_1.toml`):

```toml
[channel]
# ... other settings ...
enableBots = true  # Set to true to enable bot spawning
```

**Default:** `false` (bots disabled)

### Disable Bots

Set `enableBots = false` or omit the setting entirely. No bot-related code will execute.

## Current Implementation (Phase 1)

Phase 1 provides **minimal, non-interactive bots** that prove bots can exist safely in the world.

### What Bots Can Do

✅ **Exist in maps** - Bots spawn and are visible to real players
✅ **Have appearance** - Bots have a name, level, job, and basic stats
✅ **Stand still** - Bots occupy a position but don't move

### What Bots Cannot Do (Yet)

❌ **Move** - Bots are stationary
❌ **Attack** - Bots don't engage mobs or use skills  
❌ **Interact** - Bots don't use NPCs, portals, or chat
❌ **Join parties/guilds** - Bots are excluded from social features
❌ **Gain EXP/loot** - Bots don't participate in progression

## Technical Details

### Bot Identification

Bots are identified internally by:
- **Negative character IDs**: Bot IDs are -1, -2, -3, etc. (never collide with real players)
- **`isBot` flag**: `Player.isBot` is set to `true`
- **Stub connection**: Bots use a `botConn` stub that implements `mnet.Client` with no-ops

### Bot Lifecycle

1. **Spawn**: Server calls `InitializeBots()` after world connection
   - Currently spawns 1 test bot in Henesys (map 100000000)
2. **Visibility**: Bots use the same `addPlayer()` logic as real players
   - Real players see bots via `packetMapPlayerEnter`
3. **Cleanup**: `RemoveAllBots()` called during channel shutdown
   - Removes bots from maps and player collections

### Database Persistence

Bots are **ephemeral** and never written to the database:
- `saver.go` checks `p.isBot` and skips persistence
- Bots disappear on server restart (by design)
- No DB schema changes required

### Code Organization

```
channel/
├── bot_conn.go      # Stub connection for bots (mnet.Client impl)
├── player.go        # newBotPlayer() factory, isBot flag
├── server.go        # Bot management (spawn, remove, init)
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

## Future Phases

### Phase 2: Basic Movement
- Timer-driven left/right movement
- Jumping
- Facing direction changes

### Phase 3: Combat
- Auto-attack nearby monsters
- Basic skill usage
- Fixed damage (no EXP/loot generation)

### Phase 4: Event-Compatible
- Optional participation in GM events
- Removable mid-event
- No reward generation

### Phase 5: Advanced
- Smarter AI and pathfinding
- Configurable bot profiles
- Per-map behavior scripts
- Dynamic spawn/despawn rules

## Troubleshooting

### Bots don't appear

1. Check `enableBots = true` in your config file
2. Check server logs for "Bots are enabled" message
3. Verify Henesys map (100000000) exists in NX data

### Server crashes on startup

1. Ensure NX data is loaded (`Data.nx`)
2. Check for map ID 100000000 in NX data
3. Review server logs for specific error

### Bots appear but don't move

This is expected behavior in Phase 1. Movement will be added in Phase 2.

## Safety Features

- **Hard-coded limit**: Currently only 1 bot spawns (easy to expand later)
- **Config toggle**: Single flag disables all bot logic
- **No DB writes**: Bots can't corrupt database
- **Negative IDs**: Impossible to collide with real players
- **Stub connection**: Bots can't send malformed packets

## Development Notes

When extending bot functionality:

1. Always check `server.enableBots` before bot operations
2. Exclude bots from systems that assume real players (EXP, loot, parties)
3. Test with `enableBots = false` to ensure no regressions
4. Log bot spawn/despawn events for debugging
5. Use `plr.isBot` checks where bot behavior should differ

## Examples

### Spawning a Custom Bot

```go
// In server.go or similar
if server.enableBots {
    err := server.SpawnBot("MyBot", mapID, portalID)
    if err != nil {
        log.Printf("Failed to spawn bot: %v", err)
    }
}
```

### Checking if Player is Bot

```go
if plr.isBot {
    // Skip EXP gain, loot, etc.
    return
}
```

### Removing Specific Bot

```go
// Bots are tracked in server.bots slice
for i, bot := range server.bots {
    if bot.Name == "TestBot" {
        bot.inst.removePlayer(bot, false)
        server.players.RemoveFromConn(bot.Conn)
        server.bots = append(server.bots[:i], server.bots[i+1:]...)
        break
    }
}
```

## Testing

Run bot-specific tests:

```bash
go test ./channel -run TestBot
```

All tests should pass. Note that `TestBotPlayerCreation` is skipped as it requires NX data.

## License

Bot system is part of Valhalla and follows the same license.
