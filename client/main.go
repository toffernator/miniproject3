package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/toffernator/miniproject3/api"
	"google.golang.org/grpc"
)

var randomizedName string

var (
	nameFlag    = flag.String("name", "name", "The name of your user")
	addressFlag = flag.String("address", "localhost:50000", "The address of the front-end to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addressFlag, grpc.WithInsecure(), grpc.WithBlock())
	must(err)
	client := api.NewAuctionClient(conn)

	randomizedName = randomizeName()
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
		Label: "Amount",
	}

	isDone := false
	var amount int
	for !isDone {
		inputtedAmount, err := bidPrompt.Run()
		must(err)

		amount, err = strconv.Atoi(inputtedAmount)
		isDone = err == nil
		if !isDone {
			log.Println("You must enter a valid integer to bid")
		}
	}

	payload := api.BidMsg{
		Amount: int32(amount),
		User:   randomizedName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ack, err := c.Bid(ctx, &payload)
	if err != nil {
		log.Fatalf("There was an internal error while executing the bid request: %v", err)
	}

	log.Println("Bid ", ack.Status.String())
}

func Result(c api.AuctionClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	payload := api.Empty{}
	outcome, err := c.Result(ctx, &payload)
	if err != nil {
		log.Fatalf("There was an internal error while processing the result request: %v", err)
	}

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

func randomizeName() string {
	rand.Seed(time.Now().UnixNano())
	return *nameFlag + "-" + strconv.Itoa(rand.Int())
}
