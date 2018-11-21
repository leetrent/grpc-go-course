package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/LeeTrent/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doClientStreaming(client)
	//doBiDiStreaming(client)
	doUnaryWithDeadline(client, 5*time.Second) // should complete
	doUnaryWithDeadline(client, 1*time.Second) // should timeout

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
	fmt.Println("[greet][client.go][doClientStreaming] => BEGIN")

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

func doBiDiStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("[greet][client.go][doBiDiStreaming] => BEGIN")

	///////////////////////////////////////////////////////////////
	// Create client data
	///////////////////////////////////////////////////////////////
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lee",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Linda",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Casey",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pooh Bear",
			},
		},
	}

	///////////////////////////////////////////////////////////////
	// PSEUDOCODE:
	///////////////////////////////////////////////////////////////
	// 1. Create a stream by invoking the client
	// 2. Send mulitple messages to client using a go routine
	// 3. Receive mulitple messages from client using a go routine
	// 4. Block until everything has completed
	///////////////////////////////////////////////////////////////

	waitChannel := make(chan struct{})

	// 1. Create a stream by invoking the client
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("\n[greet][client.go][doBiDiStreaming] => Error encountered when invoking client.GreetEveryone(): %v", err)
	}

	// 2. Send mulitple messages to client using a go routine
	go func() {
		for _, req := range requests {
			fmt.Printf("\nSending request: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// 3. Receive mulitple messages from client using a go routine
	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatalf("\n[greet][client.go][doBiDiStreaming] => Error encountered when invoking stream.Recv(): %v", err)
					break
				}
			}
			fmt.Printf("Received: %v:", response.GetResult())
		}
		close(waitChannel)
	}()

	// 4. Block until everything has completed
	<-waitChannel

	fmt.Println("[greet][client.go][doBiDiStreaming] => END")
}

func doUnaryWithDeadline(client greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Println("\n[greet][client.go][doUnaryWithDeadline] => BEGIN")

	request := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lee",
			LastName:  "Ceccato",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := client.GreetWithDeadline(ctx, request)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("\n[greet][client.go][doUnaryWithDeadline] => Deadline Exceeded")
				fmt.Printf("\n[greet][client.go][doUnaryWithDeadline] => status.FromError.Code(): %v", statusErr.Code())
			} else {
				fmt.Printf("\n[greet][client.go][doUnaryWithDeadline] => Unexpected error when calling client.GreetWithDeadline: %v", statusErr)
			}
		} else {
			log.Fatalf("[greet][client.go][doUnaryWithDeadline] => Unexpected error when calling client.GreetWithDeadline: %v", err)
		}
		return
	}

	log.Printf("[client.go][doUnaryWithDeadline] Response from client.GreetWithDeadline: %v", response.Result)
	fmt.Println("\n[greet][client.go][doUnaryWithDeadline] => END")
}
