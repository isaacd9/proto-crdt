package twophaseset

import (
	"fmt"

	pb "github.com/isaacd9/proto-crdt/pb/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func New(msgs []proto.Message) (*pb.TwoPhaseSet, error) {
	set := &pb.TwoPhaseSet{
		Added:   []*anypb.Any{},
		Removed: []*anypb.Any{},
	}

	for _, m := range msgs {
		if err := Insert(set, m); err != nil {
			return nil, fmt.Errorf("could not insert %+v: %+v", m, err)
		}
	}

	return set, nil
}

func contains(set []*anypb.Any, msg *anypb.Any) bool {
	for _, item := range set {
		if proto.Equal(item, msg) {
			return true
		}
	}
	return false
}

func Insert(set *pb.TwoPhaseSet, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	// Add to Added set
	if !contains(set.Added, encoded) {
		set.Added = append(set.Added, encoded)
	}
	return nil
}

func Remove(set *pb.TwoPhaseSet, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	// Add to Removed set
	if contains(set.Added, encoded) && !contains(set.Removed, encoded) {
		set.Removed = append(set.Removed, encoded)
	}
	return nil
}

func Contains(set *pb.TwoPhaseSet, el proto.Message) (bool, error) {
	encoded, err := anypb.New(el)
	if err != nil {
		return false, fmt.Errorf("could not encode message as any: %+v", err)
	}

	return contains(set.Added, encoded) && !contains(set.Removed, encoded), nil
}

func Len(set *pb.TwoPhaseSet) int {
	return len(set.Added) - len(set.Removed)
}

func Elements(set *pb.TwoPhaseSet) ([]proto.Message, error) {
	r := []proto.Message{}
	for _, item := range set.Added {
		// If we're removed this item, we don't want to include it
		if contains(set.Removed, item) {
			continue
		}

		pb, err := item.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("invalid message: %+v", err)
		}
		r = append(r, pb)
	}
	return r, nil
}

func Merge(sets ...*pb.TwoPhaseSet) (*pb.TwoPhaseSet, error) {
	c := &pb.TwoPhaseSet{}

	for _, set := range sets {
		for _, itemAny := range set.Added {
			item, err := itemAny.UnmarshalNew()
			if err != nil {
				return nil, fmt.Errorf("invalid message: %+v", err)
			}
			if err := Insert(c, item); err != nil {
				return nil, fmt.Errorf("could not insert item: %+v", err)
			}
		}

		for _, itemAny := range set.Removed {
			item, err := itemAny.UnmarshalNew()
			if err != nil {
				return nil, fmt.Errorf("invalid message: %+v", err)
			}
			if err := Remove(c, item); err != nil {
				return nil, fmt.Errorf("could not remove item: %+v", err)
			}
		}
	}

	return c, nil
}
