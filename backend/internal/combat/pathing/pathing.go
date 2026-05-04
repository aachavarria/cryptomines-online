// Package pathing computes tile-based movement paths on the battle map.
// Used by fleets each round to spend their movement budget reaching the
// optimal tile for engaging or kiting (GDD §8.6.3, §8.6.10).
package pathing

import (
	"container/heap"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// FindPath returns the shortest tile-by-tile path from start to goal using
// 8-directional A* with Chebyshev heuristic (matches GO2's "king moves"
// movement model).
//
// Behavior:
//   - The path INCLUDES start as the first element and goal as the last.
//   - If goal is blocked but reachable adjacent tiles exist, returns the
//     best partial path that ends adjacent to goal (used when the AI wants
//     to close to engagement range without entering the enemy tile).
//   - If maxSteps >= 0, the returned path is truncated to that many steps
//     of movement (path length-1). A non-empty truncated prefix is fine.
//   - Returns nil if no path exists at all.
func FindPath(world tilemap.Blockable, unitID string, start, goal tilemap.Position, maxSteps int) []tilemap.Position {
	if start.Equals(goal) {
		return []tilemap.Position{start}
	}

	// If the goal is blocked, route to whichever neighbor of goal yields
	// the shortest reachable approach.
	targetGoals := []tilemap.Position{goal}
	if world.IsBlocked(goal, unitID) {
		targetGoals = tilemap.Neighbors8(goal)
		if len(targetGoals) == 0 {
			return nil
		}
	}

	best := []tilemap.Position(nil)
	bestLen := -1
	for _, g := range targetGoals {
		if world.IsBlocked(g, unitID) && !g.Equals(start) {
			continue
		}
		p := aStar(world, unitID, start, g)
		if p == nil {
			continue
		}
		if bestLen < 0 || len(p) < bestLen {
			best = p
			bestLen = len(p)
		}
	}
	if best == nil {
		return nil
	}
	if maxSteps >= 0 && len(best)-1 > maxSteps {
		best = best[:maxSteps+1]
	}
	return best
}

type node struct {
	pos    tilemap.Position
	g      int // steps from start
	f      int // g + heuristic
	parent *node
	idx    int // heap index
}

type openSet []*node

func (h openSet) Len() int            { return len(h) }
func (h openSet) Less(i, j int) bool  { return h[i].f < h[j].f }
func (h openSet) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].idx = i; h[j].idx = j }
func (h *openSet) Push(x interface{}) { n := x.(*node); n.idx = len(*h); *h = append(*h, n) }
func (h *openSet) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func aStar(world tilemap.Blockable, unitID string, start, goal tilemap.Position) []tilemap.Position {
	open := &openSet{}
	heap.Init(open)

	startNode := &node{pos: start, g: 0, f: tilemap.Chebyshev(start, goal)}
	heap.Push(open, startNode)

	closed := map[tilemap.Position]int{}
	closed[start] = 0

	openByPos := map[tilemap.Position]*node{start: startNode}

	for open.Len() > 0 {
		cur := heap.Pop(open).(*node)
		delete(openByPos, cur.pos)
		if cur.pos.Equals(goal) {
			return reconstruct(cur)
		}
		for _, n := range tilemap.Neighbors8(cur.pos) {
			if !n.Equals(goal) && world.IsBlocked(n, unitID) {
				continue
			}
			ng := cur.g + 1
			if prev, seen := closed[n]; seen && prev <= ng {
				continue
			}
			closed[n] = ng
			next := &node{
				pos:    n,
				g:      ng,
				f:      ng + tilemap.Chebyshev(n, goal),
				parent: cur,
			}
			if existing, inOpen := openByPos[n]; inOpen {
				existing.g = ng
				existing.f = next.f
				existing.parent = cur
				heap.Fix(open, existing.idx)
			} else {
				heap.Push(open, next)
				openByPos[n] = next
			}
		}
	}
	return nil
}

func reconstruct(end *node) []tilemap.Position {
	out := []tilemap.Position{}
	for n := end; n != nil; n = n.parent {
		out = append(out, n.pos)
	}
	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
