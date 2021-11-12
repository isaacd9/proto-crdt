package pn_counter

import (
	pb "github.com/isaacd9/proto-crdt/pb/v1"
)

type Counter = pb.PNCounter

func New(id string) *pb.PNCounter {
	return &pb.PNCounter{
		Identifier: id,
		Increments: make(map[string]uint64),
		Decrements: make(map[string]uint64),
	}
}

func Increment(p *pb.PNCounter, n uint64) {
	p.Increments[p.Identifier] += n
}

func Decrement(p *pb.PNCounter, n uint64) {
	p.Decrements[p.Identifier] += n
}

func Value(p *pb.PNCounter) int64 {
	var t int64
	for _, c := range p.Increments {
		t += int64(c)
	}
	for _, c := range p.Decrements {
		t -= int64(c)
	}
	return t
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func Merge(id string, counters ...*pb.PNCounter) *pb.PNCounter {
	c := New(id)

	for _, counter := range counters {
		for id, count := range counter.Increments {
			c.Increments[id] = max(count, c.Increments[id])
		}
	}

	for _, counter := range counters {
		for id, count := range counter.Decrements {
			c.Decrements[id] = max(count, c.Decrements[id])
		}
	}

	return c
}
