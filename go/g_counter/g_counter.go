package g_counter

import (
	pb "github.com/isaacd9/proto-crdt/pb/v1"
)

func New(id string) *pb.GCounter {
	return &pb.GCounter{
		Identifier: id,
		Counts:     make(map[string]uint64),
	}
}

func Increment(p *pb.GCounter, n uint64) {
	p.Counts[p.Identifier] += n
}

func Count(p *pb.GCounter) uint64 {
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

func Merge(id string, counters ...*pb.GCounter) *pb.GCounter {
	c := New(id)

	for _, counter := range counters {
		for id, count := range counter.Counts {
			c.Counts[id] = max(count, c.Counts[id])
		}
	}

	return c
}
