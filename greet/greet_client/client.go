package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/LeeTrent/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("[greet][client.go][main()] => BEGIN")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[greet][client.go][main]=>[grpc.Dial]: Error: %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("[greet][client.go][main()] => Created greetpb.NewGreetServiceClient ...")

	//doUnary(client)
	//doServerStreaming(client)
	doClientStreaming(client)
	fmt.Println("\n[greet][client.go][main()] => END")
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("\n[greet][client.go][doUnary] => BEGIN")

	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lee",
			LastName:  "Ceccato",
		},
	}

	response, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("[client.go][doUnary] Error encountered when invoking Greet RPC %v: ", err)
	}

	log.Printf("[client.go][doUnary] Response from Greet RPC: %v", response.Result)
	fmt.Println("\n[greet][client.go][doUnary] => END")
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("\n[greet][client.go][doServerStreaming] => BEGIN")

	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lee",
			LastName:  "Ceccato",
		},
	}

	resultStream, err := client.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("\n[greet][client.go][doServerStreaming] => greetpb.GreetServiceClient.GreetManyTimes(): Error: %v", err)
	}

	for {
		msg, err := resultStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("\n[greet][client.go][doServerStreaming] => greetpb.GreetServiceClient.GreetManyTimes.Recv(): Error: %v", err)
			}
		}
		log.Printf("%v", msg.GetResult())
	}
}

func doClientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("\n[greet][client.go][doClientStreaming] => BEGIN\n")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lee",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Linda",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Casey",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pooh Bear",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("\n[greet][client.go][doClientStreaming] => Error encountered when invoking client.LongGreet(): %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("\n[client.go][doClientStreaming] => Error encountered when invoking stream.CloseAndRecv(): %v", err)
	}
	fmt.Printf("\n[greet][client.go][doClientStreaming] => LongGreetResponse %v", response)
}
