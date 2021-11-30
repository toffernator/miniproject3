package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/toffernator/miniproject3/api"
	"google.golang.org/grpc"
)

var clients = make([]*CombinedClient, 0)

type AuctionServer struct {
	api.UnimplementedAuctionServer
}

func (s AuctionServer) Bid(ctx context.Context, msg *api.BidMsg) (*api.Ack, error) {
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
		return &api.Ack{Status: api.Ack_SUCCESS}, nil
	} else {
		return &api.Ack{Status: api.Ack_FAILED}, nil
	}

}

func (s AuctionServer) Result(ctx context.Context, empty *api.Empty) (*api.Outcome, error) {
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
	replicas := os.Args[1:]
	for _, replica := range replicas {
		clients = append(clients, newClient(replica))
	}

	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterAuctionServer(grpc.NewServer(), AuctionServer{})
	log.Printf("Auction Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
