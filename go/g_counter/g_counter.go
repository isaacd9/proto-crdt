package g_counter

import (
	pb "github.com/isaacd9/proto-crdt/pb/v1"
)

type Counter = pb.GCounter

func New(id string) *Counter {
	return &pb.GCounter{
		Identifier: id,
		Counts:     make(map[string]uint64),
	}
}

func Increment(p *Counter, n uint64) {
	p.Counts[p.Identifier] += n
}

func Value(p *Counter) uint64 {
	var t uint64
	for _, c := range p.Counts {
		t += c
	}
	return t
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func Merge(id string, counters ...*Counter) *Counter {
	c := New(id)

	for _, counter := range counters {
		for id, count := range counter.Counts {
			c.Counts[id] = max(count, c.Counts[id])
		}
	}

	return c
}
