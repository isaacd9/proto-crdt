package g_set

import (
	"fmt"

	pb "github.com/isaacd9/proto-crdt/pb/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func New(msgs []proto.Message) (*pb.GSet, error) {
	set := &pb.GSet{
		Elements: []*anypb.Any{},
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

func Insert(set *pb.GSet, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	if !contains(set.Elements, encoded) {
		set.Elements = append(set.Elements, encoded)
	}
	return nil
}

func Contains(set *pb.GSet, el proto.Message) (bool, error) {
	encoded, err := anypb.New(el)
	if err != nil {
		return false, fmt.Errorf("could not encode message as any: %+v", err)
	}

	return contains(set.Elements, encoded), nil
}

func Len(set *pb.GSet) int {
	return len(set.Elements)
}

func Elements(set *pb.GSet) ([]proto.Message, error) {
	r := []proto.Message{}
	for _, item := range set.Elements {
		pb, err := item.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("invalid message: %+v", err)
		}
		r = append(r, pb)
	}
	return r, nil
}

func Merge(sets ...*pb.GSet) (*pb.GSet, error) {
	c := &pb.GSet{}
	for _, set := range sets {
		elements, err := Elements(set)
		if err != nil {
			return nil, fmt.Errorf("could not get elements of set %+v: %+v", set, err)
		}

		for _, item := range elements {
			if err := Insert(c, item); err != nil {
				return nil, fmt.Errorf("could not insert item: %+v", err)
			}
		}
	}
	return c, nil
}
