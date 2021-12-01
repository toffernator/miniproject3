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
var ongoing bool

func main() {
	ongoing = true
	startRMServer()
}

type RMServer struct {
	isRegistered map[string]bool
	api.UnimplementedRMServer
}

func (as *RMServer) Bid(ctx context.Context, bm *api.BidMsg) (*api.Ack, error) {
	if ongoing {
		lock.Lock()
		defer lock.Unlock()

		log.Printf("Received Bid")

		if !as.isRegistered[bm.User] {
			as.isRegistered[bm.User] = true
			log.Printf("Registered %s to the auction", bm.User)
		}

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
	} else {
		return &api.Ack{Status: api.Ack_ENDED}, errors.New("cannot bid as the auction has ended")
	}
}

func (as *RMServer) Result(context.Context, *api.Empty) (*api.Outcome, error) {
	if ongoing {
		log.Printf("The highest current bid: %d", highestBid)
	} else {
		log.Printf("The auction has ended and the winner is: %s, with a bid of: %d", user, highestBid)
	}
	return &api.Outcome{ResultOrHighest: highestBid, Winner: user}, nil
}

func (r *RMServer) ForceBid(ctx context.Context, bm *api.BidMsg) (*api.Ack, error) {
	if ongoing {
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
	log.Printf("cannot bid as the auction has ended")
	return &api.Ack{Status: api.Ack_ENDED}, nil
}

func (r *RMServer) EndAuction(context.Context, *api.Empty) (*api.Ack, error) {
	ongoing = false
	return &api.Ack{Status: api.Ack_SUCCESS}, nil
}

func startRMServer() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed ot listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := RMServer{
		isRegistered: make(map[string]bool, 0),
	}

	api.RegisterRMServer(grpcServer, &server)
	log.Printf("RM Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
