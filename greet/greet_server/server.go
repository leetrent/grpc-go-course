package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/LeeTrent/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("\n[Greet][server.go][(*server)Greet] => *greetpb.GreetRequest: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := &greetpb.GreetResponse{
		Result: result,
	}

	return response, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("\n[Greet][server.go][(*server)GreetManyTimes] => *greetpb.GreetManyTimesRequest: %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	for ii := 0; ii < 10; ii++ {
		result := "Hello " + firstName + " (#" + strconv.Itoa(ii) + ")"
		response := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(response)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil // no error
}

func main() {
	fmt.Println("[server.go][main] Starting server ...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("[server.go][main][net.Listen]: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[server.go][main][Server.Serve]: %v", err)
	}
}
