package g_set

import (
	"fmt"

	pb "github.com/isaacd9/proto-crdt/pb/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Set = pb.GSet

func New() *Set {
	return &Set{
		Elements: []*anypb.Any{},
	}
}

func contains(set []*anypb.Any, msg *anypb.Any) bool {
	for _, item := range set {
		if proto.Equal(item, msg) {
			return true
		}
	}
	return false
}

func Insert(set *Set, el proto.Message) error {
	encoded, err := anypb.New(el)
	if err != nil {
		return fmt.Errorf("could not encode message as any: %+v", err)
	}

	if !contains(set.Elements, encoded) {
		set.Elements = append(set.Elements, encoded)
	}
	return nil
}

func Contains(set *Set, el proto.Message) (bool, error) {
	encoded, err := anypb.New(el)
	if err != nil {
		return false, fmt.Errorf("could not encode message as any: %+v", err)
	}

	return contains(set.Elements, encoded), nil
}

func Len(set *Set) int {
	return len(set.Elements)
}

func Elements(set *Set) ([]proto.Message, error) {
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

func Merge(sets ...*Set) (*Set, error) {
	c := &Set{}
	for _, set := range sets {
		for _, item := range set.Elements {
			if err := Insert(c, item); err != nil {
				return nil, fmt.Errorf("could not insert item: %+v", err)
			}
		}
	}
	return c, nil
}
