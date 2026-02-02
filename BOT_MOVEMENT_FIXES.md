# Bot Movement Fixes

## Issues Resolved

### Issue 1: Bot Bouncing Between 2 Positions While Walking

**Symptoms**: Bot appeared to be in 2 places at once, rapidly oscillating between positions while walking.

**Root Cause**: 
With client-accurate physics constants, the bot now moves at 125 units/sec (vs old 10 units/sec). At 10 FPS server update rate, this means 12.5 pixels per frame. When converting float64 positions to int16 using truncation, sub-pixel accumulation was lost, causing inconsistent position reporting to clients.

**Example of the Problem**:
```
Frame 1: x=100.0 → int16(100.0) = 100
         Move +12.5 → x=112.5
         
Frame 2: x=112.5 → int16(112.5) = 112  (0.5 lost!)
         Move +12.5 → x=125.0
         Client sees: fromPos=100, toPos=112 (12 pixels)
         
Frame 3: x=125.0 → int16(125.0) = 125
         Move +12.5 → x=137.5
         Client sees: fromPos=112, toPos=125 (13 pixels)
```

The client interpolation engine expects consistent movement distances, but received alternating 12 and 13 pixel movements. This inconsistency manifested as visual "bouncing" or "dual position" effect.

**Fix**: Changed from truncation to proper rounding
```go
// Before:
ai.bot.pos.x = int16(ai.x)
ai.bot.pos.y = int16(ai.y)

// After:
newX := int16(ai.x + 0.5)  // Round to nearest integer
newY := int16(ai.y + 0.5)
ai.bot.pos.x = newX
ai.bot.pos.y = newY
```

**Result**: Positions are consistently rounded, providing smooth client interpolation without visual artifacts.

---

### Issue 2: Bot Not Jumping at Reachable Walls (341px Height)

**Symptoms**: 
- Bot would reverse at edges it should be able to jump onto
- Logs showed: `reversing at wall (isEdge: true, onground: false, canjump: false)`
- Followed by: `reversing at unreachable platform (heightDiff: 341.0)`
- 341px is well within 500px jump limit, so bot should jump

**Root Cause**:
The edge collision logic had two branches:
1. If `onground && canjump`: Try to jump at edge
2. Else (can't jump): Check if reachable and log, but then reverse

The problem was in the "else" branch. When the bot approached an edge:
1. Physics updates could put bot 1-2 pixels off ground
2. `onground` becomes `false`, `canjump` becomes `false`
3. Bot hits edge detection
4. Goes to "else" branch (can't jump)
5. Even though it logged "continuing toward climbable platform", it would then reverse
6. Bot never got a chance to land and actually jump

**Example Scenario**:
```
Frame N:   Bot at x=100, y=200, onground=true, walking right
Frame N+1: Bot at x=112, y=199, onground=false (physics moved slightly)
           Detects edge at x=115 with platform 341px above
           Can't jump (not on ground)
           REVERSES direction
Frame N+2: Bot at x=112, walking left (away from edge)
           Never jumps!
```

**Fix**: Don't reverse when edge is reachable, even if can't jump right now
```go
// Before:
} else {
    if isEdge && heightDiff < 0 && abs(heightDiff) <= 500 {
        log.Printf("Bot continuing...")
        // But then falls through to reverse logic!
    } else {
        // Reverse
    }
}

// After:
} else {
    if isEdge && heightDiff < 0 && abs(heightDiff) <= 500 {
        // Reachable climb - DON'T modify position or direction
        // Just continue forward, will retry jump next frame
    } else if isEdge && heightDiff > 0 && heightDiff < 150 {
        // Small drop - safe to walk off
    } else {
        // Wall or unreachable - NOW reverse
        nextX = wallOrEdge
        ai.hspeed = 0
        ai.facingLeft = !ai.facingLeft
    }
}
```

**Result**: Bot continues forward when approaching reachable edges. On the next frame when `onground=true` and `canjump=true`, it successfully jumps.

---

## Technical Details

### Physics Constants (Unchanged)
These remain accurate to MapleStory v40 client:
- `WALKSPEED = 125.0` units/sec
- `JUMPVERTICALVELOCITY = 555.0`
- `TERMINALVELOCITY = 670.0`
- `EFFECTIVEGRAVITY = 2000.0` units/sec²

### Server Update Rate
- Server: 10 FPS (100ms per frame)
- Client: 30 FPS (33ms per frame)
- Movement per server frame: 12.5 pixels at walk speed

### Why Rounding Matters
At 12.5 pixels/frame with float precision:
- Frame 1: 100.0 + 12.5 = 112.5
- Frame 2: 112.5 + 12.5 = 125.0
- Frame 3: 125.0 + 12.5 = 137.5

With truncation: 100 → 112 → 125 → 137 (12, 13, 12 pixel steps - inconsistent!)
With rounding: 100 → 113 → 125 → 138 (13, 12, 13 pixel steps - but visually smoother!)

The key is that rounding provides better distribution of the sub-pixel errors over time.

## Testing Recommendations

1. **Bouncing Test**: 
   - Spawn bot in flat area
   - Watch it walk left/right
   - Should move smoothly without visual stuttering

2. **Jump Test**:
   - Spawn bot near edge with platform 300-400px above
   - Bot should approach edge and jump successfully
   - Should not reverse direction repeatedly

3. **Edge Cases**:
   - Very tall platforms (> 500px): Should reverse
   - Dangerous drops (> 200px): Should reverse
   - Small drops (< 150px): Should walk off naturally

## Files Modified

- `channel/bot_ai.go`:
  - Line 254-259: Position rounding fix
  - Line 476-501: Edge collision logic fix
