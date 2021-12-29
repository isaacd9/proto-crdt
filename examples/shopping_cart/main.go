package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/isaacd9/proto-crdt/examples/shopping-cart/pb"

	"github.com/isaacd9/proto-crdt/or_set"
	crdt_pb "github.com/isaacd9/proto-crdt/pb/v1"
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
	set   *crdt_pb.ORSet
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

	set, err := or_set.New([]proto.Message{})
	if err != nil {
		log.Fatalf("could not initialize ORSet: %+v", err)
	}

	s := &ShoppingCart{
		peers: peers,
		set:   set,
	}
	pb.RegisterShoppingCartServer(grpcServer, s)

	log.Printf("starting server on %q", *addr)
	grpcServer.Serve(lis)
}
