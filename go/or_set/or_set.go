package or_set

import (
	"fmt"

	"github.com/google/uuid"
	pb "github.com/isaacd9/proto-crdt/pb/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func New(msgs []proto.Message) (*pb.ORSet, error) {
	set := &pb.ORSet{}

	for _, m := range msgs {
		if err := Insert(set, m); err != nil {
			return nil, fmt.Errorf("could not insert %+v: %+v", m, err)
		}
	}

	return set, nil
}

func Insert(set *pb.ORSet, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	set.Added = append(set.Added, &pb.ORSet_Element{
		Value:      encoded,
		Identifier: uuid.New().String(),
	})
	return nil
}

func Remove(set *pb.ORSet, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	for _, added := range set.Added {
		// Remove all elements we've added which match the object
		if proto.Equal(encoded, added.Value) {
			set.Removed = append(set.Removed, added)
		}
	}
	return nil
}

func Contains(set *pb.ORSet, el proto.Message) (bool, error) {
	encoded, err := anypb.New(el)
	if err != nil {
		return false, fmt.Errorf("could not encode message as any: %+v", err)
	}

	identifiers := make(map[string]struct{})
	for _, added := range set.Added {
		if proto.Equal(encoded, added.Value) {
			identifiers[added.Identifier] = struct{}{}
		}
	}

	for _, removed := range set.Removed {
		if proto.Equal(encoded, removed.Value) {
			delete(identifiers, removed.Identifier)
		}
	}

	if len(identifiers) > 0 {
		return true, nil
	}

	return false, nil
}

func Merge(sets ...*pb.ORSet) (*pb.ORSet, error) {
	c := &pb.ORSet{}
	for _, set := range sets {
		c.Added = append(c.Added, set.Added...)
		c.Removed = append(c.Removed, set.Removed...)
	}
	return c, nil
}
