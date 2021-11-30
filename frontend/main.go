package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/toffernator/miniproject3/api"
	"google.golang.org/grpc"
)

var clients = make([]*CombinedClient, 0)

type RMServer struct {
	api.UnimplementedRMServer
}

func (s *RMServer) Bid(msg *api.BidMsg) *api.Ack {
	failed_clients := make([]*CombinedClient, 0)
	for _, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ack, err := client.Bid(ctx, msg)
		if err != nil {
			log.Printf("Could not reach %s. Error: %v", client.server_addr, err)
		}
		if ack.Status != api.Ack_SUCCESS || err != nil {
			failed_clients = append(failed_clients, client)
		}
	}
	// This client won the consensus.
	if len(failed_clients) <= len(clients)/2 {
		for _, client := range failed_clients {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			ack, err := client.ForceBid(ctx, msg)
			if err != nil {
				log.Printf("Replica %s not responding. It is probably dead.", client.server_addr)
			}
			if ack.Status == api.Ack_FAILED {
				log.Println("Force sync failed. This means a competing write bid higher, rendering the force sync unnecessary.")
			} else if ack.Status == api.Ack_EXCEPTION {
				log.Println("Unknown exception occured in the replica.")
			}
		}
		return &api.Ack{Status: api.Ack_SUCCESS}
	} else {
		return &api.Ack{Status: api.Ack_FAILED}
	}

}

func (s *RMServer) Result(empty *api.Empty) (*api.Outcome, error) {
	replica_idx := rand.Intn(len(clients))
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		outcome, err := clients[replica_idx].Result(ctx, &api.Empty{})
		if err != nil {
			log.Printf("Could not reach replica %s. Trying another.", clients[replica_idx].server_addr)
		} else {
			return outcome, nil
		}
	}
	return nil, errors.New("No replicas available.")
}

type CombinedClient struct {
	api.RMClient
	server_addr string
}

func newClient(server string) *CombinedClient {
	conn, err := grpc.Dial(server, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("Failed to connect to next node")
	}

	rmClient := api.NewRMClient(conn)

	client := CombinedClient{
		RMClient:    rmClient,
		server_addr: server,
	}

	return &client
}

func main() {
	for _, replica := range os.Args {
		clients = append(clients, newClient(replica))
	}
}
