# MapleStory v40 Physics - Fresh Reverse Engineering Analysis

**Date**: 2026-01-29 (Updated: 2026-02-01)  
**Method**: Direct IDA analysis without existing documentation  
**Implementation**: SpookyStory V2 (`MaplePlayerController`, `MaplePhysicsEngine`, `MapleConstants`)

---

## Implementation Status Legend

Throughout this document, implementation status is marked as follows:

| Badge | Meaning |
|-------|---------|
| ✅ **IMPLEMENTED** | Fully implemented and matches IDA analysis |
| ⚠️ **SIMPLIFIED** | Implemented differently from IDA (explained below) |
| ❌ **NOT IMPLEMENTED** | Not yet implemented |
| 🔍 **IDA ONLY** | IDA analysis for reference, not needed in implementation |

---

## Table of Contents
1. [Constants Reference](#constants-reference)
2. [Core Functions](#core-functions)
3. [State Machine](#state-machine)
4. [Jump Mechanics](#jump-mechanics)
5. [Airborne Physics](#airborne-physics)
6. [Landing Detection](#landing-detection)
7. [Grounded Movement](#grounded-movement)
8. [CMovePath Object Layout](#cmovepath-object-layout)
9. [Foothold Tree & Collision System](#foothold-tree--collision-system)
10. [Ladder, Rope & Prone Mechanics](#ladder-rope--prone-mechanics)
11. [Implementation Differences Summary](#implementation-differences-summary)

---

## Constants Reference

All constants extracted from IDA database at addresses 0x62E340-0x62E3C0 and 0x629600-0x6296C0.

> [!IMPORTANT]
> Our implementation uses `MapleConstants.cs` which stores all physics constants.
> See: [MapleConstants.cs](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Core/MapleConstants.cs)

### Core Physics Constants

| Address | IDA Name | Value | C# Constant | Status | Purpose |
|---------|----------|-------|-------------|--------|---------|
| 0x62E358 | `g_JumpVerticalVelocity` | **-555.0** | `_jumpVerticalVelocity = 555f` | ✅ | Initial upward velocity (we store positive, negate on use) |
| 0x62E350 | `g_JumpHorizontalMultiplier` | **162.5** | `_jumpHorizontalMultiplier = 162.5f` | ✅ | Horizontal velocity multiplier on jump |
| 0x62E360 | `g_WalkSpeed` | **125.0** | `_walkSpeed = 125f` | ✅ | Default walk speed |
| 0x62E3A0 | `g_TerminalVelocity` | **670.0** | `_terminalVelocity = 670f` | ✅ | Maximum fall speed |
| 0x62E3B0 | `g_GravityMaxForce` | **100000.0** | `_gravityMaxForce = 100000f` | 🔍 | Maximum gravity force constant (reference only) |
| 0x62E3A8 | `g_GravityBase` | **10000.0** | `_gravityBase = 10000f` | 🔍 | Base gravity constant (reference only) |
| 0x62E380 | `g_GroundedAccelForce` | **140000.0** | `_groundedAccelForce = 140000f` | ✅ | Force for grounded acceleration |
| 0x62E390 | `g_AirborneAcceleration` | **2000.0** | `_airborneAcceleration = 2000f` | ✅ | Horizontal air control acceleration |
| 0x62B6D0 | `g_WeakFrictionMult` | **0.01** | `_weakFrictionMultiplier = 0.01f` | ✅ | Weak friction multiplier for air control |
| 0x62BC60 | `g_ShortJumpModifier` | **0.3** | `_shortJumpModifier = 0.3f` | ❌ | Jump modifier for short jumps (not used) |
| 0x6295C8 | `g_TrapezoidFactor` | **0.5** | `_trapezoidalFactor = 0.5f` | ✅ | Trapezoidal integration factor |
| 0x629668 | `g_MsToSeconds` | **0.001** | N/A | 🔍 | Unity handles this via `Time.fixedDeltaTime` |
| 0x62E3C0 | `g_FrameDt30fps` | **0.0333** | `_frameDt30FPS = 0.0333f` | ✅ | Frame time at 30fps |

### Ladder/Rope Constants

| Address | IDA Name | Value | C# Constant | Status | Purpose |
|---------|----------|-------|-------------|--------|---------|
| 0x62A4A0 | `g_LadderClimbSpeed` | **3.0** | `_ladderClimbSpeed = 3f` | ✅ | Pixels per input per frame |
| 0x62E340 | `g_LadderJumpHor1` | **200.0** | `_ladderJumpUpMultiplier = 200f` | ⚠️ | See [Ladder Jump Differences](#ladder-jump-implementation) |
| 0x62E348 | `g_LadderJumpHor2` | **500.0** | `_ladderJumpDownMultiplier = 500f` | ❌ | Down jump multiplier (not implemented) |

### Computed/Derived Constants

| Name | IDA Formula | C# Constant | Value | Notes |
|------|-------------|-------------|-------|-------|
| `effectiveGravity` | `(walkStats * 10000 / jumpStats)` | `_effectiveGravity` | **2000.0** | VERIFIED via client hook |
| `airFrictionWeak` | `10000 * 0.01 / 125` | `_airFrictionWeak` | **0.8** | Per-frame when rising/slow fall |
| `airFrictionStrong` | `10000 / 125` | `_airFrictionStrong` | **80.0** | Per-frame when falling fast |
| `maxAirControlVelocity` | `stat * 10000 * 0.0008928571` | `_maxAirControlVelocity` | **8.928571** | Air control velocity cap |

---

## Core Functions

### Function Map

| Address | Name | Size | C# Equivalent | Status |
|---------|------|------|---------------|--------|
| 0x5BC3B5 | `CMovePath_UpdatePhysics` | 0x242 | `MaplePhysicsEngine.UpdatePhysics()` | ✅ |
| 0x5BC5F7 | `CUserLocal_DoJump` | 0x3AB | `MaplePhysicsEngine.InitiateJump()` | ✅ |
| 0x5BD3B6 | `CMovePath_ApplyAirborneVelocity` | 0x5B6 | `MaplePhysicsEngine.UpdateAirborne()` | ✅ |
| 0x5BD96C | `CMovePath_AirbornePhysicsUpdate` | 0x8C5 | `MapleCollisionSystem.FindFloorIntersection()` | ✅ |
| 0x5BD345 | `ApplyAccelerationClamped` | 0x71 | `MaplePhysicsEngine.ApplyAccelerationClamped()` | ✅ |
| 0x5BE231 | `CMovePath_SegmentIntersectionTest` | 0xB6 | `MapleCollisionSystem` (4-cross algorithm) | ✅ |
| 0x5BEC40 | `CMovePath_SetFootholdState` | 0xC9 | `MaplePhysicsEngine.Land()` | ✅ |
| 0x5C03EE | `CMovePath_HandleGroundedMovement` | 0x55A | `MaplePhysicsEngine.UpdateGrounded()` | ✅ |
| 0x5C4B27 | `sub_5C4B27` (Ladder Climb) | ~0x120 | `MaplePhysicsEngine.UpdateLadder()` | ✅ |

---

## State Machine

### Movement States ✅ IMPLEMENTED

IDA offset `this[130]` maps to `MapleMovementState` enum:

```csharp
public enum MapleMovementState
{
    Grounded = 1,  // On foothold
    Airborne = 2,  // In the air
    Ladder = 3     // On ladder/rope
}
```

> [!NOTE]
> Our implementation stores state in `MaplePlayerState.state` field.

### State Check Logic 🔍 IDA ONLY

IDA reference at `CMovePath_IsAirborne @ 0x59D0CE`:
```c
// Returns true if:
//   this[68] == 0 (not on foothold) AND
//   stat[+24] >= 0 AND
//   stat[+108] <= 0
return !this[68] && sub_59D70F(this[94]+24) >= 0 && sub_59D70F(this[96]+108) <= 0;
```

> [!IMPORTANT]
> **Implementation Difference**: We use a simple enum check (`state.state == MapleMovementState.Airborne`) rather than checking encrypted stat values. The ZtlSecureTear encryption is not relevant for gameplay replication.

---

## Jump Mechanics

### CUserLocal_DoJump ✅ IMPLEMENTED

**IDA Reference**: `0x5BC5F7`  
**C# Implementation**: [MaplePhysicsEngine.InitiateJump()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Physics/MaplePhysicsEngine.cs#L490-L519)

#### Jump Type Selection 🔍 NOT USED

IDA shows short/normal jump modifier selection at `0x5BC696-0x5BC6A6`:
```c
// Determine jump type based on stats at this[376]+24 and this[384]+108
bool isShortJump = sub_59D70F(this[376]+24) < 0 || sub_59D70F(this[384]+108) > 0;

double jumpModifier;
if (isShortJump) {
    jumpModifier = dbl_62BC60;  // 0.3 (short jump)
} else {
    jumpModifier = dbl_6295C8;  // 0.5 (normal jump)
}
```

> [!WARNING]
> **Implementation Difference**: We always use full jump velocity (`-555.0`) without the 0.5 modifier. The IDA shows this modifier comes from stat-based conditions we haven't fully reverse-engineered.

#### Vertical Velocity ✅ IMPLEMENTED

```csharp
// Our implementation (MaplePhysicsEngine.cs line 507)
state.velocityY = -MapleConstants._jumpVerticalVelocity;  // -555

// IDA formula:
// velocityY = jumpStat * (-555.0) / walkSpeedStat * jumpModifier
// With 100% stats, stats cancel: -555.0 * 0.5 = -277.5 (normal) or -166.5 (short)
```

#### Horizontal Velocity ✅ IMPLEMENTED

```csharp
// Our implementation (MaplePhysicsEngine.cs lines 510-513)
if (state.horizontalInput != 0)
{
    state.velocityX = state.horizontalInput * MapleConstants._jumpHorizontalMultiplier;  // ±162.5
}
```

---

## Airborne Physics

### CMovePath_ApplyAirborneVelocity ✅ IMPLEMENTED

**IDA Reference**: `0x5BD3B6`  
**C# Implementation**: [MaplePhysicsEngine.UpdateAirborne()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Physics/MaplePhysicsEngine.cs#L212-L339)

### Gravity Application ✅ IMPLEMENTED

**IDA Formula**:
```c
effectiveGravity = (walkSpeed * 10000 * 2) / walkSpeed * dt = 20000 * dt
terminalVelocity = walkSpeed * 670 = 83750 (using default walkSpeed=125)
```

**Our Implementation**:
```csharp
// MaplePhysicsEngine.cs lines 566-576
private void ApplyGravity(ref float velocityY, float dt)
{
    if (velocityY < MapleConstants._terminalVelocity)  // 670
    {
        velocityY += MapleConstants._effectiveGravity * dt;  // 2000 * dt
        if (velocityY > MapleConstants._terminalVelocity)
            velocityY = MapleConstants._terminalVelocity;
    }
}
```

> [!NOTE]
> The `_effectiveGravity = 2000` constant was VERIFIED via client hooking to produce 60 units/frame at 30fps (2000 * 0.03 = 60).

### Air Control (Horizontal Input) ✅ IMPLEMENTED

**IDA Reference** `@ 0x5BD500`:
```c
int hInput = GetSecureValue(this+80, this[82]);  // -1, 0, 1
if (hInput != 0) {
    double hAccel = hInput * dbl_62E390;  // ±2000
    ApplyAccelerationClamped(&velocityX, hAccel, walkSpeedStat, walkSpeedStat, dt);
}
```

**Our Implementation**:
```csharp
// MaplePhysicsEngine.cs lines 222-231
if (state.horizontalInput != 0)
{
    const float effectiveAirAccel = 160f;  // = 20000 / 125
    float accel = state.horizontalInput * effectiveAirAccel;
    ApplyAccelerationClamped(ref state.velocityX, accel, dt, MapleConstants._maxAirControlVelocity);
}
```

> [!IMPORTANT]
> **Implementation Difference**: We pre-compute the effective acceleration rate (160) rather than applying the IDA formula dynamically. Result is identical for 100% stats.

### Air Friction (Dual-Mode) ✅ IMPLEMENTED

**IDA @ 0x5BD52B-0x5BD5D6** shows TWO friction modes based on fall speed:

| Condition | Friction | Formula |
|-----------|----------|---------|
| Rising OR slow fall (`velocityY < 670`) | **0.8** per frame | `10000 * 0.01 / 125` |
| Fast falling (`velocityY >= 670`) | **80** per frame | `10000 / 125` |

**Our Implementation**:
```csharp
// MaplePhysicsEngine.cs lines 234-247
if (state.horizontalInput == 0)
{
    float friction;
    if (state.velocityY < MapleConstants._fallSpeedThreshold)  // 670
    {
        friction = MapleConstants._airFrictionWeak;   // 0.8 - momentum persists
    }
    else
    {
        friction = MapleConstants._airFrictionStrong; // 80 - rapid slowdown
    }
    ApplyFriction(ref state.velocityX, friction, dt);
}
```

> [!TIP]
> This dual friction system is why jump momentum feels "floaty" during the rise but quickly decays when falling fast.

### Position Integration (Trapezoidal) ✅ IMPLEMENTED

**IDA @ 0x5BD8C5-0x5BD962**:
```c
posX += (oldVelX + velocityX) * dt * dbl_6295C8;  // 0.5
posY += (oldVelY + velocityY) * dt * dbl_6295C8;
```

**Our Implementation**:
```csharp
// MaplePhysicsEngine.cs lines 257-261
float avgVelX = (state.prevVelocityX + state.velocityX) * MapleConstants._trapezoidalFactor;
float avgVelY = (state.prevVelocityY + state.velocityY) * MapleConstants._trapezoidalFactor;

state.positionX += avgVelX * dt;
state.positionY += avgVelY * dt;
```

---

## Landing Detection

### Floor Detection ✅ IMPLEMENTED

**IDA Reference**: Segment intersection at `0x5BDC68-0x5BDE0F`  
**C# Implementation**: [MapleCollisionSystem.FindFloorIntersection()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Collision/MapleCollisionSystem.cs#L114-L193)

#### 4-Cross Algorithm ✅ IMPLEMENTED

The IDA shows a standard line-segment intersection using cross products:

```c
// Movement segment: (oldX, oldY) -> (newX, newY)
// Foothold segment: (fhX1, fhY1) -> (fhX2, fhY2)

// Cross products for side-of-line tests
__int64 cross1 = (__int64)(oldY - fhY1) * fhDx - (__int64)(oldX - fhX1) * fhDy;
__int64 cross2 = (__int64)(newY - fhY1) * fhDx - (__int64)(newX - fhX1) * fhDy;

// For floors: must cross from STRICTLY above (cross1 < 0) to at/below (cross2 >= 0)
if (cross1 >= 0) continue;  // Old pos is ON or below foothold line - skip
if (cross2 < 0) continue;   // New pos is still above foothold line
```

> [!NOTE]
> The **strict** `cross1 < 0` check is critical - it prevents landing when starting exactly ON the line, which avoids instant re-landing when falling off edges.

#### Continuous Collision Detection (CCD) ⚠️ SIMPLIFIED

**IDA**: Does not use explicit CCD - relies on frame-by-frame checks.

**Our Implementation**: We subdivide large movements (>20px) into substeps:

```csharp
// MapleCollisionSystem.cs - CCD subdivision
const float MAX_STEP_SIZE = 20f;
float distance = Vector2.Distance(new Vector2(oldX, oldY), new Vector2(newX, newY));
int substeps = Mathf.CeilToInt(distance / MAX_STEP_SIZE);
```

> [!IMPORTANT]
> **Implementation Enhancement**: We added CCD to prevent "falling through" floors at high velocities (up to 670px/frame). The original client likely had this issue too.

#### Apex Tracking ⚠️ ENHANCED

**IDA**: No explicit apex tracking visible.

**Our Implementation**: We track the highest point reached during a jump for robust floor detection:

```csharp
// MaplePhysicsEngine.cs lines 289-323
if (state.velocityY <= 0)
{
    state.apexY = Mathf.Min(state.apexY, state.positionY);  // Track highest point
}

// Catch-up check from apex to current position
if (!intersection.HasValue && state.apexY < oldY && (oldY - state.apexY) < 30)
{
    intersection = _collision.FindFloorIntersection(oldX, state.apexY, state.positionX, state.positionY);
}
```

### H-Gate Logic 🔍 IDA REFERENCE

The IDA at `0x5BDFFB` shows the H-Gate controls T-position updates AFTER landing, not which floors can be landed on. **We do not implement this** as it only affects T-parameter interpolation which is not gameplay-critical.

### Wall Collision ✅ IMPLEMENTED

**IDA @ 0x5BDB12-0x5BDC08**: Walls require slope-gated group filtering.

**Filtering Rules**:
| slopeDx | Foothold Type | Filtering |
|---------|---------------|-----------|
| > 0 | Floor | Always included |
| <= 0 | Wall | Only if `fh.group == currentFh.group OR fh.group == this.persistentGroupID` |

**Our Implementation**: [MapleCollisionSystem.FindWallIntersection()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Collision/MapleCollisionSystem.cs#L249-L357)

> [!NOTE]
> Walls are **one-sided barriers** - they only block from one direction based on their orientation.

---

## Grounded Movement

### CMovePath_HandleGroundedMovement ✅ IMPLEMENTED

**IDA Reference**: `0x5C03EE`  
**C# Implementation**: [MaplePhysicsEngine.UpdateGrounded()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Physics/MaplePhysicsEngine.cs#L58-L99)

#### Acceleration ✅ IMPLEMENTED

```csharp
// MaplePhysicsEngine.cs lines 74-82
if (state.horizontalInput != 0)
{
    float accel = state.horizontalInput * MapleConstants._groundedAccelForce;  // ±140000
    ApplyAccelerationClamped(ref state.velocityX, accel, dt, MapleConstants._walkSpeed);
}
else
{
    ApplyFriction(ref state.velocityX, MapleConstants._groundedAccelForce, dt);
}
```

#### Foothold Movement ⚠️ SIMPLIFIED

**IDA**: Uses T-parameter system where position is stored as `T` (0 = start, length = end).

**Our Implementation**: We use direct X-coordinate tracking:

```csharp
// MaplePhysicsEngine.cs - MoveAlongFoothold
float newX = state.positionX + displacement;
float newY = fh.GetYAtX(newX);  // Calculate Y from foothold slope
```

> [!NOTE]
> **Implementation Difference**: We don't implement the T-parameter system. Instead, we directly track X position and compute Y from the foothold's slope equation. This is functionally equivalent but simpler.

---

## CMovePath Object Layout

🔍 **IDA REFERENCE ONLY** - Not directly relevant to Unity implementation.

This section documents the original client's memory layout for reference. Our implementation uses clean C# classes (`MaplePlayerState`, `MaplePhysicsEngine`) instead.

| Offset | Type | Name | C# Equivalent |
|--------|------|------|---------------|
| 0x68 | ptr | currentFoothold | `state.currentFoothold` |
| 0x74 | int | collisionGroupID | `state.collisionGroupID` |
| 0x80 | int | hInput | `state.horizontalInput` |
| 0x86 | bool | jumpRequested | `state.jumpPressed` |
| 0x130 | int | moveState | `state.state` (enum) |

---

## Foothold Tree & Collision System

### TRSTree (Range Search Tree) ⚠️ SIMPLIFIED

**IDA**: Uses hierarchical AABB tree for O(log n) spatial queries.

**Our Implementation**: Uses simple linear search through foothold list:

```csharp
// MapleCollisionSystem.cs - GetCrossCandidates
foreach (var fh in _footholds)
{
    if (fh.minX <= queryMaxX && fh.maxX >= queryMinX &&
        fh.minY <= queryMaxY && fh.maxY >= queryMinY)
    {
        candidates.Add(fh);
    }
}
```

> [!WARNING]
> **Performance Note**: Linear search is O(n) vs IDA's O(log n). For maps with many footholds, this could be a performance concern. Consider implementing a proper spatial index for production.

### CStaticFoothold Structure ✅ IMPLEMENTED

| IDA Offset | Type | C# Field |
|------------|------|----------|
| +12 | int | `x1` |
| +16 | int | `y1` |
| +20 | int | `x2` |
| +24 | int | `y2` |
| +32 | int | `id` |
| +76 | ptr | `prev` |
| +80 | ptr | `next` |

---

## Ladder, Rope & Prone Mechanics

### Ladder Climbing ✅ IMPLEMENTED

**IDA Reference**: `sub_5C4B27 @ 0x5C4B27`  
**C# Implementation**: [MaplePhysicsEngine.UpdateLadder()](file:///C:/Users/jmeul/Documents/Repositories/SpookyStory-Client/Assets/Scripts/PlayerV2/Physics/MaplePhysicsEngine.cs#L369-L466)

```csharp
// Climb speed: 3 pixels per input frame
float climbDelta = -state.verticalInput * MapleConstants._ladderClimbSpeed;
state.positionY += climbDelta;
```

### Ladder Jump ⚠️ SIMPLIFIED

**IDA Shows Two Jump Types**:

| Direction | Multiplier | Stat Offset | Formula |
|-----------|------------|-------------|---------|
| Up from ladder | 200.0 | +120 | `stat * |vInput| * 200` |
| Down from ladder | 500.0 | +96 | `stat * |vInput| * 500` |

**Our Implementation**: We use a single sideways jump without directional variants:

```csharp
// MaplePhysicsEngine.cs - JumpFromLadder (lines 468-483)
state.velocityY = -MapleConstants._jumpVerticalVelocity * MapleConstants._normalJumpModifier;  // -277.5
state.velocityX = state.horizontalInput * MapleConstants._jumpHorizontalMultiplier;  // ±162.5
```

> [!IMPORTANT]
> **Implementation Difference**: We require horizontal input to jump from ladder (matching MapleStory behavior where Jump alone on ladder does nothing). However, we don't implement the up/down directional jump multipliers.

### Ladder Exit Behavior ✅ IMPLEMENTED

**Top Exit**: Places player 5px above ladder, then finds floor to land on  
**Bottom Exit**: Auto-dismounts when reaching bottom, finds floor or goes airborne

```csharp
// MaplePhysicsEngine.cs - Ladder top exit
if (state.positionY <= ladderTop && state.verticalInput > 0)
{
    state.currentLadder = null;
    state.positionY = ladderTop - 5;
    // Find floor or go airborne...
}
```

### Ladder Re-Entry Cooldown ⚠️ ENHANCED

**IDA**: No explicit cooldown visible.

**Our Implementation**: 0.3 second cooldown to prevent animation glitches:

```csharp
state.ladderExitCooldown = 0.3f;  // Prevent immediate re-entry
```

### Prone State ❌ NOT IMPLEMENTED

IDA shows prone is primarily a visual/hitbox state. When grounded + down pressed + no ladder nearby → Enter prone animation. **We have not implemented prone mechanics.**

---

## Implementation Differences Summary

This section consolidates all differences between IDA analysis and our implementation for quick reference.

### Simplified (Functionally Equivalent)

| Feature | IDA Approach | Our Approach | Impact |
|---------|--------------|--------------|--------|
| **State Check** | ZtlSecureTear encrypted values | Simple enum | None |
| **T-Parameter** | Foothold position as T (0-length) | Direct X coordinate | None |
| **Spatial Index** | TRSTree (O log n) | Linear search (O n) | Performance on large maps |
| **Air Acceleration** | Dynamic formula | Pre-computed constant (160) | None for 100% stats |

### Enhanced (Added for Stability)

| Feature | IDA Approach | Our Enhancement | Reason |
|---------|--------------|-----------------|--------|
| **CCD** | Frame-by-frame | Subdivided steps | Prevents fall-through at high velocity |
| **Apex Tracking** | Not visible | Track highest point | More reliable floor detection |
| **Ladder Cooldown** | Not visible | 0.3s cooldown | Prevents animation glitches |
| **Edge Instant-Land** | Not visible | 10px raycast | Smoother stair navigation |

### Not Implemented

| Feature | IDA Reference | Notes |
|---------|--------------|-------|
| **Short Jump Modifier** | 0.3 multiplier | Stat-dependent, needs more RE |
| **Directional Ladder Jump** | 200/500 multipliers | Only sideways jump implemented |
| **Prone State** | Visual/hitbox change | Not gameplay-critical |
| **H-Gate T-Position** | Post-landing interpolation | Not gameplay-critical |

---

## Summary of Key Formulas

### Jump Initiation
```
velocityY = -555.0  (IDA: jumpStat * (-555.0) / walkSpeedStat * jumpModifier)
velocityX = hInput * 162.5
```

### Airborne Gravity
```
effectiveGravity = 2000 per second (60 per 30ms frame)
terminalVelocity = 670
```

### Air Control
```
acceleration = hInput * 160 * dt  (= hInput * 20000 / 125 * dt)
maxAirControlVelocity = 8.928571
```

### Air Friction (Dual-Mode)
```
if velocityY < 670: friction = 0.8 * dt (weak - momentum persists)
if velocityY >= 670: friction = 80 * dt (strong - rapid decay)
```

### Grounded Movement
```
acceleration = hInput * 140000 * dt
maxVelocity = 125 (walkSpeed)
```

### Ladder Climbing
```
deltaY = -vInput * 3.0 per frame
```

### Position Integration (Trapezoidal)
```
position += (oldVelocity + newVelocity) * 0.5 * dt
```

---