package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/isaacd9/proto-crdt/examples/shopping-cart/pb"
	"github.com/isaacd9/proto-crdt/or_set"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	addr = flag.String("addr", ":8080", "address to listen on")
	peer = flag.String("peers", "", "peer list as a comma seperated list of strings. If not provided, assume bootstrap")
)

type ShoppingCart struct {
	*pb.UnimplementedShoppingCartServer

	peers []string
	items map[string]uint64
}

func (s *ShoppingCart) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	s.items[r.Item.Name] = r.Item.Quantity
	if err := s.replicate(); err != nil {
		return nil, fmt.Errorf("replication err: %+v", err)
	}
	return &pb.AddResponse{}, nil
}

func (s *ShoppingCart) Remove(ctx context.Context, r *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	delete(s.items, r.Item)
	if err := s.replicate(); err != nil {
		return nil, fmt.Errorf("replication err: %+v", err)
	}
	return &pb.RemoveResponse{}, nil
}

func (s *ShoppingCart) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	var items []*pb.CartItem
	for n, q := range s.items {
		items = append(items, &pb.CartItem{
			Name:     n,
			Quantity: q,
		})
	}
	return &pb.GetResponse{Items: items}, nil
}

func (s *ShoppingCart) replicate() error {
	return nil
}

func (s *ShoppingCart) UpdateCart(ctx context.Context, r *pb.CartRequest) (*pb.CartResponse, error) {
	var items = []proto.Message{}
	for name, quantity := range s.items {
		items = append(items, &pb.CartItem{
			Name:     name,
			Quantity: quantity,
		})
	}
	thisOrSet, err := or_set.New(items)
	if err != nil {
		return nil, fmt.Errorf("could not create ORSet: %+v", s)
	}
	merged, err := or_set.Merge(thisOrSet, r.Set)
	if err != nil {
		return nil, fmt.Errorf("could not merge ORSets: %+v", s)
	}
	_ = merged

	for name, q := range s.items {
		contains, err := or_set.Contains(merged, &pb.CartItem{
			Name:     name,
			Quantity: q,
		})
		if err != nil {
			return nil, fmt.Errorf("could not determine if merged contains item %q: %+v", name, s)
		}

		if !contains {
			delete(s.items, name)
		}
	}
	return nil, nil
}

func getPeers(peerSt string) []string {
	if peerSt == "" {
		return []string{}
	}
	return strings.Split(peerSt, ",")
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%+v", *addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	peers := getPeers(*peer)

	if err != nil {
		log.Fatalf("could not initialize ORSet: %+v", err)
	}

	s := &ShoppingCart{
		peers: peers,
	}

	pb.RegisterShoppingCartServer(grpcServer, s)

	log.Printf("starting server on %q", *addr)
	grpcServer.Serve(lis)
}
