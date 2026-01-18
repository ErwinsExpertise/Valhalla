# Anti-Cheat & Ban System

This document describes the comprehensive anti-cheat and ban system implemented in Valhalla.

## Overview

The anti-cheat system provides:

- **Automated cheat detection** across multiple categories (combat, movement, inventory, economy, skills, packets)
- **Configurable ban enforcement** with temporary and permanent bans
- **Escalation system** that automatically upgrades repeat offenders from temporary to permanent bans
- **Rolling window violation tracking** to prevent false positives from single events
- **GM commands** for manual ban management
- **Comprehensive audit logging** of all violations and enforcement actions
- **IP banning** support with configurable policies

## Architecture

The system consists of several key components:

### 1. Ban Service (`anticheat_ban.go`)

Handles all ban operations:
- Checking if accounts/characters/IPs are banned
- Issuing temporary and permanent bans
- Escalation from temporary to permanent bans
- Ban history tracking
- Automatic expiration of old bans

### 2. Violation Detector (`anticheat_detector.go`)

Tracks violations in rolling time windows:
- Records violation events to database
- Maintains in-memory counters for active violation windows
- Evaluates thresholds and triggers enforcement actions
- Cleans up expired violation data

### 3. Detection Helpers (`anticheat_helpers.go`)

Provides convenient methods for specific violation types:
- `DetectExcessiveDamage()` - Combat damage validation
- `DetectAttackSpeedHack()` - Attack speed validation
- `DetectSpeedHack()` - Movement speed validation
- `DetectTeleportHack()` - Teleport validation
- `DetectInvalidEquip()` - Equipment validation
- `DetectInvalidTrade()` - Trade validation
- `DetectDuplication()` - Item duplication detection
- And many more...

### 4. Configuration (`anticheat_config.go`)

Comprehensive configuration structure allowing fine-tuning of:
- Ban durations and escalation thresholds
- Violation detection windows and thresholds
- Enable/disable individual detection categories
- IP ban policies

## Database Schema

The system uses four main tables:

### `bans`
Stores all ban records (active and historical)
- Supports account, character, and IP bans
- Tracks ban type (temporary/permanent)
- Records who issued the ban (GM or system)
- Stores ban reason and timestamps

### `ban_escalation`
Tracks temporary ban counts for escalation
- Counts temporary bans per account
- Tracks when permanent ban was issued
- Used to automatically escalate repeat offenders

### `violation_logs`
Comprehensive log of all violations
- Records every violation with full context
- Includes detection details and severity
- Tracks which map violation occurred in
- Logs what action was taken (if any)

### `violation_counters`
Maintains rolling window violation counts
- Tracks violation count within current window
- Records window start time and last violation time
- Used for threshold evaluation

## Configuration

See `config_anticheat_example.toml` for a complete configuration example.

Key configuration options:

```toml
[anticheat]
enabled = true
defaultTempBanDuration = "168h"  # 7 days
tempBansBeforePermanent = 3      # Escalate after 3 temp bans
ipBanMode = "permanent_only"     # When to apply IP bans
gmBansIncrementCounter = false   # Whether GM bans count for escalation
```

Each detection category has its own configuration:

```toml
[anticheat.combatDetection]
enabled = true
excessiveDamageThreshold = 5     # Violations within window
excessiveDamageWindow = "5m"     # Rolling window duration
excessiveDamageBanType = "temporary"
```

## GM Commands

### `/ban <player> [duration|perm] [reason...]`

Ban a player with optional duration and reason.

Examples:
```
/ban PlayerName perm Repeated hacking
/ban PlayerName 24 Speedhacking (24 hour ban)
/ban PlayerName Using unauthorized tools
```

### `/unban <player>`

Remove all active bans for a player.

Example:
```
/unban PlayerName
```

### `/banhistory <player>`

View a player's complete ban history.

Example:
```
/banhistory PlayerName
```

### `/violations <player> [limit]`

View recent violations for a player.

Examples:
```
/violations PlayerName
/violations PlayerName 20
```

## Detection Categories

### Combat Detection
- **Excessive Damage**: Detects damage exceeding calculated maximum (150% margin)
- **Attack Speed**: Monitors attack intervals for speed hacks
- **Invalid Skills**: Checks if player uses skills they don't have

### Movement Detection
- **Speed Hacks**: Monitors movement speed (110% margin)
- **Teleport Hacks**: Detects invalid teleportation
- **Invalid Positions**: Checks for impossible positions (walls, out of bounds)

### Inventory Detection
- **Invalid Equip**: Detects equipping items player doesn't own or can't use
- **Invalid Item Use**: Checks item usage validity

### Economy Detection
- **Invalid Trades**: Monitors trade operations for exploits
- **Item Duplication**: Detects potential duplication attempts (very strict)
- **Overflow Exploits**: Checks for meso/item count overflow/underflow

### Skill Detection
- **Cooldown Bypass**: Monitors skill cooldown violations
- **Unlearned Skills**: Detects use of skills player hasn't learned

### Packet Detection
- **Invalid Sequences**: Checks for impossible packet sequences
- **Malformed Packets**: Detects corrupted or tampered packets

## Rolling Window Logic

The system uses rolling time windows to prevent false positives:

1. When a violation occurs, it's added to a counter for that violation type
2. The counter tracks violations within a configurable time window (default 5 minutes)
3. Only when the threshold is exceeded within the window is action taken
4. After action, the counter resets for that violation type

This prevents single network glitches or calculation differences from triggering bans.

## Escalation System

The escalation system automatically upgrades repeat offenders:

1. First offense: Temporary ban (default 7 days)
2. Second offense: Another temporary ban
3. Third offense: Another temporary ban
4. Fourth offense: **Automatic permanent ban**

The threshold is configurable via `tempBansBeforePermanent`.

GM-issued bans can optionally be excluded from escalation counting via `gmBansIncrementCounter`.

## Ban Checking

Bans are checked at multiple points:

1. **During login**: Via `CheckPlayerBan()` in server initialization
2. **During character selection**: Before allowing into game
3. **Periodically**: Old temporary bans are automatically expired

The system checks in this order:
1. Account ban
2. Character ban  
3. IP ban

If any check fails, the connection is rejected with ban details.

## IP Banning

IP bans can be configured with three modes:

- **never**: Never apply IP bans (default for dev/testing)
- **permanent_only**: Only apply IP bans for permanent account bans
- **always**: Apply IP bans for all bans (strictest)

IP banning is useful for preventing ban evasion but should be used carefully to avoid blocking legitimate users behind shared IPs (schools, internet cafes, etc.).

## Logging & Auditing

All enforcement actions are logged:

- **Violation logs**: Every detected violation with full context
- **Ban records**: Complete ban history with issuer and reason
- **Action tracking**: What action was taken for each violation

Logs can be queried via:
- GM commands (`/violations`, `/banhistory`)
- Direct database queries
- Future admin panel (if implemented)

## Integration Points

To integrate anti-cheat checks into game logic:

1. **Get detector instance**: `server.violationDetector`
2. **Call appropriate detection method**:

```go
// Example: Check damage
if server.violationDetector != nil {
    server.violationDetector.DetectExcessiveDamage(
        player,
        damageDealt,
        calculatedMaxDamage,
    )
}

// Example: Check movement
if server.violationDetector != nil {
    server.violationDetector.DetectSpeedHack(
        player,
        calculatedSpeed,
        maxAllowedSpeed,
    )
}
```

3. **The system handles**:
   - Logging the violation
   - Updating counters
   - Checking thresholds
   - Issuing bans if needed

## Performance Considerations

The system is designed for minimal performance impact:

- **In-memory counters**: Active violation tracking uses memory, not constant DB queries
- **Batch cleanup**: Expired data cleaned up periodically, not per-violation
- **Indexed queries**: Database tables are properly indexed
- **Configurable windows**: Adjust window sizes based on server load

## Future Enhancements

Potential improvements:

- [ ] Web admin panel for ban management
- [ ] Whitelist system for trusted players
- [ ] Machine learning integration for pattern detection
- [ ] Integration with Discord for ban notifications
- [ ] Custom violation types via plugins
- [ ] Ban appeal system
- [ ] Temporary mute/restrictions (in addition to bans)

## Testing Recommendations

Before enabling in production:

1. **Enable logging only first**: Set all `banType` to `"none"` initially
2. **Monitor violation logs**: Check for false positives
3. **Tune thresholds**: Adjust based on legitimate player behavior
4. **Test GM commands**: Verify ban/unban functionality
5. **Test escalation**: Manually trigger escalation scenarios
6. **Load test**: Ensure performance under high player counts

## Support

For issues, questions, or feature requests related to the anti-cheat system:
- Check the main project README
- Review the configuration examples
- Examine the code comments
- Open an issue on GitHub

## License

This anti-cheat system is part of the Valhalla project and shares the same license.
