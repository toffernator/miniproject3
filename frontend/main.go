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
		if err != nil || ack.Status != api.Ack_SUCCESS {
			failed_clients = append(failed_clients, client)
		}
		if err != nil || ack.Status == api.Ack_ENDED {
			return &api.Ack{Status: api.Ack_ENDED}, nil
		}
	}

	// This frontend won the consensus.
	if len(failed_clients) <= len(clients)/2 {
		for _, client := range failed_clients {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			ack, err := client.ForceBid(ctx, msg)
			if err != nil {
				log.Printf("Replica %s not responding. It is probably dead.", client.server_addr)
			} else if ack.Status == api.Ack_FAILED {
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
	for i := 0; i < 5; i++ {
		replica_idx := rand.Intn(len(clients))
		log.Printf("Trying to reach replica %v", replica_idx)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		outcome, err := clients[replica_idx].Result(ctx, &api.Empty{})
		if err != nil {
			log.Printf("Could not reach replica %s. Trying another.", clients[replica_idx].server_addr)
		} else {
			log.Printf("Returning result to client %v", outcome)
			return outcome, nil
		}
	}
	for i, client := range clients {
		log.Printf("Trying to reach replica %v", i)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		outcome, err := client.Result(ctx, &api.Empty{})
		if err != nil {
			log.Printf("Could not reach replica %s. Trying another.", client.server_addr)
		} else {
			log.Printf("Returning result to client %v", outcome)
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

func EndAuction() {

	duration := time.Duration(30) * time.Second

	f := func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		for _, c := range clients {
			c.EndAuction(ctx, &api.Empty{})
		}
	}

	timer := time.AfterFunc(duration, f)
	defer timer.Stop()

	time.Sleep(time.Second * 30)

	log.Printf("The auction has ended!")
}

func main() {
	replicas := os.Args[1:]
	rand.Seed(time.Now().UnixMicro())
	for _, replica := range replicas {
		clients = append(clients, newClient(replica))
	}

	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	go EndAuction()

	api.RegisterAuctionServer(grpcServer, AuctionServer{})
	log.Printf("Auction Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
