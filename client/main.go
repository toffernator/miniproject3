package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/toffernator/miniproject3/api"
	"google.golang.org/grpc"
)

var (
	nameFlag    = flag.String("name", "name", "The name of your user")
	addressFlag = flag.String("address", "localhost:50000", "The address of the front-end to connect to")
)

func main() {
	conn, err := grpc.Dial(*addressFlag, grpc.WithInsecure(), grpc.WithBlock())
	must(err)
	client := api.NewAuctionClient(conn)

	clientLoop(client)
}

func clientLoop(c api.AuctionClient) {
	isDone := false
	for !isDone {
		prompt := promptui.Select{
			Label: "Choose an option",
			Items: []string{"Bid", "Result", "Quit"},
		}

		_, result, err := prompt.Run()
		must(err)

		switch result {
		case "Bid":
			Bid(c)
		case "Result":
			Result(c)
		case "Quit":
			Leave()
		}
	}
}

func Bid(c api.AuctionClient) {
	bidPrompt := promptui.Prompt{
		Label: "Amount: ",
	}

	amount, err := bidPrompt.Run()
	must(err)

	intAmount, err := strconv.Atoi(amount)
	must(err)

	payload := api.BidMsg{
		Amount: int32(intAmount),
		User:   *nameFlag,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	c.Bid(ctx, &payload)
}

func Result(c api.AuctionClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	payload := api.Empty{}
	outcome, _ := c.Result(ctx, &payload)

	log.Printf("Result: %d by %s", outcome.ResultOrHighest, outcome.Winner)
}

func Leave() {
	os.Exit(0)
}

func must(err error) {
	if err != nil {
		log.Fatalf("Failed with err code: %v", err)
	}
}
