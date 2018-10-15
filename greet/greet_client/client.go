package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LeeTrent/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("[client.go][main] BEGIN")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[client.go][main]=>[grpc.Dial]: Error: %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("[client.go][main]=>[greetpb.NewGreetServiceClient]: Created client (%f)", client)

	doUnary(client)
	fmt.Println("[client.go][main] END")
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("[client.go][doUnary] BEGIN")

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

	fmt.Println("[client.go][doUnary] END")
}
