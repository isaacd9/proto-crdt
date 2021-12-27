package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/isaacd9/proto-crdt/examples/monotonic-counter/pb"
	"github.com/isaacd9/proto-crdt/examples/monotonic_counter/pb/pb"
	"github.com/isaacd9/proto-crdt/g_counter"
	crdt_pb "github.com/isaacd9/proto-crdt/pb/v1"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
)

var (
	addr       = flag.String("addr", ":8080", "address to listen on")
	peer       = flag.String("peers", "", "peer list as a comma seperated list of strings. If not provided, assume bootstrap")
	identifier = flag.String("identifier", "", "identifier for this host. default is addr")
	initial    = flag.Uint64("init", 0, "initial value")
	tick       = flag.Duration("tick", 1*time.Second, "duration between ticks")
)

func getPeers(peerSt string) []string {
	if peerSt == "" {
		return []string{}
	}
	return strings.Split(peerSt, ",")
}

type Counter struct {
	Id string

	pb.CounterServer
	stopCh      chan interface{}
	counter     *crdt_pb.GCounter
	peerClients []pb.CounterClient

	sync.Mutex
}

func NewCounter(id string, initial uint64, peers []string) *Counter {
	peerClients := make([]pb.CounterClient, len(peers))
	for i, peer := range peers {
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to create peer client %q: %v", peer, err)
		}
		peerClients[i] = pb.NewCounterClient(conn)
	}

	gcounter := g_counter.New(id)
	g_counter.Increment(gcounter, initial)

	return &Counter{
		Id:          id,
		peerClients: peerClients,
		counter:     gcounter,
	}
}

func (c *Counter) Peer(s pb.Counter_PeerServer) error {
	g, ctx := errgroup.WithContext(s.Context())

	msgCh := make(chan *pb.MergeRequest)

	g.Go(func() error {
		for {
			msg, err := s.Recv()
			if err != nil {
				return err
			}
			msgCh <- msg
		}
	})

	g.Go(func() error {
		for {
			select {
			case msg := <-msgCh:
				c.Lock()
				// log.Printf("recv from call: %+v. us: %+v", msg, c.counter)
				c.counter = g_counter.Merge(c.Id, msg.Counter, c.counter)
				s.Send(&pb.MergeResponse{
					Counter: c.counter,
				})
				c.Unlock()
			case <-ctx.Done():
				return nil
			}
		}
		return nil
	})

	return g.Wait()
}

func (c *Counter) Value(context.Context, *pb.ValueRequest) (*pb.ValueResponse, error) {
	return &pb.ValueResponse{
		value: g_counter.Value(c.counter),
	}, nil
}

func (c *Counter) Tick(ctx context.Context, ticker *time.Ticker) error {
	g, ctx := errgroup.WithContext(ctx)

	var (
		recvCh = make(chan *crdt_pb.GCounter)
	)

	peerStreams := make([]pb.Counter_PeerClient, len(c.peerClients))
	for i, client := range c.peerClients {
		stream, err := client.Peer(ctx)
		if err != nil {
			return err
		}
		peerStreams[i] = stream

		g.Go(func() error {
			for {
				msg, err := stream.Recv()
				if err != nil {
					return err
				}
				recvCh <- msg.Counter
			}
		})
	}

	g.Go(func() error {
		for {
			select {
			case <-ticker.C:
				c.Lock()
				// log.Printf("value (%+v): %d", c.counter, g_counter.Value(c.counter))
				log.Printf("value: %d", g_counter.Value(c.counter))
				g_counter.Increment(c.counter, 1)
				for _, peer := range peerStreams {
					err := peer.Send(&pb.MergeRequest{
						Counter: c.counter,
					})
					if err != nil {
						return err
					}
				}
				c.Unlock()
			case counter := <-recvCh:
				c.Lock()
				// log.Printf("recv from peer: %+v. us: %+v", counter, c.counter)
				c.counter = g_counter.Merge(c.Id, counter, c.counter)
				c.Unlock()
			}
		}
	})

	return g.Wait()
}

func (c *Counter) Stop() {
	c.stopCh <- struct{}{}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%+v", *addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	id := *identifier
	if id == "" {
		id = *addr
	}

	peers := getPeers(*peer)
	c := NewCounter(id, *initial, peers)

	go func() {
		if err := c.Tick(context.Background(), time.NewTicker(*tick)); err != nil {
			log.Fatalf("tick failed: %+v", err)
		}
	}()

	pb.RegisterCounterServer(grpcServer, c)

	log.Printf("starting server on %q", *addr)
	grpcServer.Serve(lis)
}
