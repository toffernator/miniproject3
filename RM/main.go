package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/toffernator/miniproject3/api"
	"google.golang.org/grpc"
)

const port = ":50000"

var highestBid int32
var user string
var lock sync.Mutex

func main() {
	go startAuctionServer()
	startReplicationControlServer()
}

type AuctionServer struct {
	api.UnimplementedAuctionServer
}

func (as *AuctionServer) Bid(ctx context.Context, bm *api.BidMsg) (*api.Ack, error) {
	lock.Lock()
	defer lock.Unlock()

	log.Printf("Received Bid")

	if bm.Amount > highestBid {
		log.Printf("Bid was accepted with a value: %d | Previous highest bid: %d", bm.Amount, highestBid)
		highestBid = bm.Amount
		user = bm.User
		return &api.Ack{Status: api.Ack_SUCCESS}, nil
	} else if bm.Amount <= highestBid {
		log.Printf("Bid was not accepted. New bid: %d, was equal to or lower than previous bid: %d", bm.Amount, highestBid)
		return &api.Ack{Status: api.Ack_FAILED}, nil
	}
	return &api.Ack{Status: api.Ack_EXCEPTION}, errors.New("magically broke simple algebra")
}

func (as *AuctionServer) Result(context.Context, *api.Empty) (*api.Outcome, error) {
	log.Printf("The highest current bid: %d", highestBid)

	return &api.Outcome{ResultOrHighest: highestBid, Winner: user}, nil
}

type ReplicationControlServer struct {
	api.UnimplementedReplicationControlServer
}

func (r *ReplicationControlServer) ForceBid(ctx context.Context, bm *api.BidMsg) (*api.Ack, error) {
	lock.Lock()
	defer lock.Unlock()

	log.Printf("Forcefully changed the highest bid")

	if bm.Amount == highestBid {
		highestBid = bm.Amount
		user = bm.User
		return &api.Ack{Status: api.Ack_SUCCESS}, nil
	} else {
		log.Printf("tried forcing bid with bid amount not equal to highest bid")
		return &api.Ack{Status: api.Ack_FAILED}, nil
	}
}

func startAuctionServer() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed ot listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := AuctionServer{}

	api.RegisterAuctionServer(grpcServer, &server)
	log.Printf("Auction Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startReplicationControlServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed ot listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := ReplicationControlServer{}

	api.RegisterReplicationControlServer(grpcServer, &server)
	log.Printf("Replication Contrl Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
