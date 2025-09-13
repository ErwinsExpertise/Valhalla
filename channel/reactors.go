package channel

import (
	"time"
)

// ReactorEventType mirrors basic types we care about
const (
	reactorEventHit  = 0
	reactorEventDrop = 100
)

type rect struct {
	left, top, right, bottom int16
}

func (r rect) contains(x, y int16) bool {
	if r.right < r.left || r.bottom < r.top {
		return false
	}
	return x >= r.left && x <= r.right && y >= r.top && y <= r.bottom
}

func (r *fieldReactor) calcEventRect() (rect, bool) {
	st, ok := r.info.States[int(r.state)]
	if !ok || len(st.Events) == 0 {
		return rect{}, false
	}
	ev := st.Events[0]
	if ev.LT.X == 0 && ev.LT.Y == 0 && ev.RB.X == 0 && ev.RB.Y == 0 {
		return rect{}, false
	}
	L := int16(int(ev.LT.X) + int(r.pos.x))
	T := int16(int(ev.LT.Y) + int(r.pos.y))
	R := int16(int(ev.RB.X) + int(r.pos.x))
	B := int16(int(ev.RB.Y) + int(r.pos.y))
	return rect{left: L, top: T, right: R, bottom: B}, true
}

func (r *fieldReactor) nextStateFromTemplate() (byte, bool) {
	cur := int(r.state)
	st, ok := r.info.States[cur]
	if ok && len(st.Events) > 0 {
		ns := int(st.Events[0].State)
		if _, ok2 := r.info.States[ns]; ok2 {
			return byte(ns), true
		}
	}
	if _, ok := r.info.States[cur+1]; ok {
		return byte(cur + 1), true
	}
	return r.state, false
}

func (r *fieldReactor) isTerminal() bool {
	_, ok := r.info.States[int(r.state)+1]
	return !ok
}

func (pool *reactorPool) changeState(r *fieldReactor, next byte, frameDelay int16, cause byte) {
	r.state = next
	r.frameDelay = frameDelay
	pool.instance.send(packetMapReactorChangeState(r.spawnID, r.state, r.pos.x, r.pos.y, r.frameDelay, r.faceLeft, cause))

	// After acknowledging the change to clients, run server-side side effects for this state
	pool.processStateSideEffects(r)
}

func (pool *reactorPool) leaveAndMaybeRespawn(r *fieldReactor, animationDelayMS int) {
	pool.instance.send(packetMapReactorLeaveField(r.spawnID, r.state, r.pos.x, r.pos.y))

	if r.reactorTime > 0 {
		rt := time.Duration(r.reactorTime) * time.Second
		time.AfterFunc(rt, func() {
			// Reset and show again
			r.state = 0
			r.frameDelay = 0
			pool.instance.send(packetMapReactorEnterField(r.spawnID, r.templateID, r.state, r.pos.x, r.pos.y, r.faceLeft))
		})
	}
}

// Player-led hit/touch
func (pool *reactorPool) TriggerHit(spawnID int32, cause byte) {
	r, ok := pool.reactors[spawnID]
	if !ok {
		return
	}

	next, ok := r.nextStateFromTemplate()
	if !ok || next == r.state {
		return
	}

	pool.changeState(r, next, 0, cause)

	if r.isTerminal() {
		pool.leaveAndMaybeRespawn(r, 0)
	}
}

// Server-side drop trigger by scanning reactors (called when a drop lands or after a small delay)
func (pool *reactorPool) TryTriggerByDrop(drop fieldDrop) bool {
	// Only item drops (not mesos)
	if drop.mesos > 0 {
		return false
	}

	for _, r := range pool.reactors {
		st, has := r.info.States[int(r.state)]
		if !has || len(st.Events) == 0 {
			continue
		}
		ev := st.Events[0]
		if ev.Type != reactorEventDrop {
			continue
		}

		// Item and amount match
		if ev.ReqItemID != drop.item.ID {
			continue
		}
		if ev.ReqItemCnt > 0 && int16(ev.ReqItemCnt) != drop.item.amount {
			continue
		}

		// Rectangle check if provided
		if rr, okRect := r.calcEventRect(); okRect {
			if !rr.contains(drop.finalPos.x, drop.finalPos.y) {
				continue
			}
		}

		// Progress reactor
		next, okNext := r.nextStateFromTemplate()
		if !okNext || next == r.state {
			continue
		}

		pool.changeState(r, next, 0, 0)

		// Consume the drop
		pool.instance.dropPool.removeDrop(0, drop.ID)

		if r.isTerminal() {
			pool.leaveAndMaybeRespawn(r, 0)
		}
		return true
	}
	return false
}

// processStateSideEffects inspects the current state's event metadata and performs actions like spawning mobs.
func (pool *reactorPool) processStateSideEffects(r *fieldReactor) {
	st, ok := r.info.States[int(r.state)]
	if !ok || len(st.Events) == 0 {
		return
	}
	ev := st.Events[0]

	mobID, hasMob := ev.ExtraInts["mob"]
	if !hasMob || mobID <= 0 {
		return
	}

	count := int32(1)
	if c, ok := ev.ExtraInts["count"]; ok && c > 0 {
		count = c
	} else if a, ok := ev.ExtraInts["amount"]; ok && a > 0 {
		count = a
	}
	if count < 1 {
		count = 1
	} else if count > 15 {
		count = 15
	}

	spawnPos := r.pos
	if rr, okRect := r.calcEventRect(); okRect {
		cx := int16((int(rr.left) + int(rr.right)) / 2)
		cy := int16((int(rr.top) + int(rr.bottom)) / 2)
		spawnPos = pos{x: cx, y: cy, foothold: r.pos.foothold}
	}

	for i := int32(0); i < count; i++ {
		_ = pool.instance.lifePool.spawnMobFromID(mobID, spawnPos, false, true, true, 0)
	}
}
