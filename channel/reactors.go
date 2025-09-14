package channel

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type reactorTableEntry map[string]interface{}

var reactorTable map[string]map[string]reactorTableEntry

// PopulateReactorTable from json file
func PopulateReactorTable(reactorJSON string) error {
	f, err := os.Open(reactorJSON)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &reactorTable)
}

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

// Minimal helpers to pull primitives from generic JSON
func getInt(e reactorTableEntry, key string, def int) int {
	if v, ok := e[key]; ok && v != nil {
		switch t := v.(type) {
		case float64:
			return int(t)
		case int:
			return t
		case int32:
			return int(t)
		case int64:
			return int(t)
		}
	}
	return def
}
func getString(e reactorTableEntry, key, def string) string {
	if v, ok := e[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

// entriesForState returns all actions for a reactor name whose "state" matches.
// Uses the action name as the exact key. No cross-group fallback.
func entriesForState(name string, state byte) []reactorTableEntry {
	if reactorTable == nil {
		return nil
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	group, ok := reactorTable[name]
	if !ok || len(group) == 0 {
		return nil
	}

	// Collect numeric subkeys in order, then filter by state
	type kv struct {
		n int
		k string
	}
	keys := make([]kv, 0, len(group))
	for k := range group {
		if n, err := strconv.Atoi(k); err == nil {
			keys = append(keys, kv{n: n, k: k})
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].n < keys[j].n })

	out := make([]reactorTableEntry, 0, len(keys))
	for _, it := range keys {
		e := group[it.k]
		if getInt(e, "state", -1) == int(state) {
			out = append(out, e)
		}
	}
	return out
}

func (pool *reactorPool) broadcastRedMessage(msg string) {
	if msg == "" {
		return
	}
	p := packetMessageRedText(msg)
	pool.instance.send(p)
}

// processStateSideEffects inspects either NX metadata or JSON (reactorTable) and performs actions (e.g., spawn mobs).
func (pool *reactorPool) processStateSideEffects(r *fieldReactor) {
	// 1) Handle NX-defined extras (existing logic)
	st, ok := r.info.States[int(r.state)]
	if ok && len(st.Events) > 0 {
		ev := st.Events[0]

		mobID, hasMob := ev.ExtraInts["mob"]
		if hasMob && mobID > 0 {
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
	}

	// 2) JSON-based actions: spawn ALL entries for this state (e.g., boss spawns many parts)
	entries := entriesForState(r.name, r.state)
	if len(entries) == 0 {
		return
	}

	// Send first non-empty message as a red message
	for _, e := range entries {
		if msg := getString(e, "message", ""); msg != "" {
			pool.broadcastRedMessage(msg)
			break
		}
	}

	for _, e := range entries {
		switch getInt(e, "type", -1) {
		case 1: // Mob summon
			mobID := getInt(e, "0", 0)
			if mobID <= 0 {
				continue
			}
			// Count defaults to key "2" if present (many entries use 1)
			count := getInt(e, "2", 1)
			if count < 1 {
				count = 1
			}
			if count > 15 {
				count = 15
			}
			// Position: keys "4" and "5" if present, else reactor pos
			x := int16(getInt(e, "4", int(r.pos.x)))
			y := int16(getInt(e, "5", int(r.pos.y)))
			spawnPos := pos{x: x, y: y, foothold: r.pos.foothold}

			for i := 0; i < count; i++ {
				_ = pool.instance.lifePool.spawnMobFromID(int32(mobID), spawnPos, false, true, true, 0)
			}
		default:
			// Minimal: ignore other types here
		}
	}
}
