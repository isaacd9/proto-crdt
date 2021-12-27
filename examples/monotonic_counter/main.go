package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/isaacd9/proto-crdt/examples/monotonic-counter/pb"
	"github.com/isaacd9/proto-crdt/g_counter"
	crdt_pb "github.com/isaacd9/proto-crdt/pb/v1"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
)

var (
	addr       = flag.String("addr", ":8080", "address to listen on")
	peer       = flag.String("peers", "", "peer list as a comma seperated list of strings. If not provided, assume bootstrap")
	identifier = flag.String("identifier", "", "identifier for this host. default is addr")
)

type counter struct {
	id string

	pb.CounterServer
	stopCh chan interface{}

	sync.Mutex
	counter *crdt_pb.GCounter
}

func (c *counter) Peer(s pb.Counter_PeerServer) error {
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
				c.counter = g_counter.Merge(c.id, msg.Counter, c.counter)
			case <-ctx.Done():
				return nil
			}
		}
		return nil
	})

	return g.Wait()
}

func (c *counter) Value(context.Context, *pb.ValueRequest) (*pb.ValueResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *counter) Tick(ticker *time.Ticker) error {
	for {
		log.Printf("value: %d", g_counter.Value(c.counter))
		select {
		case <-ticker.C:
			c.Lock()
			g_counter.Increment(c.counter, 1)
			c.Unlock()
		case <-c.stopCh:
			return nil
		}
	}
}

func (c *counter) Stop() {
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

	c := counter{
		counter: g_counter.New(id),
		id:      id,
	}

	go func() {
		c.Tick(time.NewTicker(1 * time.Second))
	}()

	pb.RegisterCounterServer(grpcServer, &c)

	log.Printf("starting server on %q", *addr)
	grpcServer.Serve(lis)
}
