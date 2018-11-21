package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/LeeTrent/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("[greet][server.go][(*server)LongGreet()] => BEGIN ...")

	result := ""

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// We have finished reading the client stream
				return stream.SendAndClose(&greetpb.LongGreetResponse{
					Result: result,
				})
			} else {
				log.Fatalf("\n[greet][server.go][(*server)LongGreet()] => stream.Recv() error: %v", err)
			}
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("[greet][server.go][(*server)GreetEveryone()] => BEGIN ...")

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// We're done - return nil (no error)
				return nil
			} else {
				log.Fatalf("\n[greet][server.go][(*server)GreetEveryone()] => stream.Recv() error: %v", err)
				return err
			}
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		sendError := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendError != nil {
			log.Fatalf("\n[greet][server.go][(*server)GreetEveryone()] => stream.Send() error: %v", err)
			return sendError
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Println("[Greet][server.go][(*server)GreetWithDeadline] => BEGIN ...")

	for ii := 0; ii < 3; ii++ {
		if ctx.Err() == context.Canceled {
			// Client cancelled request
			fmt.Println("[Greet][server.go][(*server)GreetWithDeadline] => Client canceled request")
			return nil, status.Error(codes.Canceled, "Client canceled request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return response, nil
}

func main() {
	fmt.Println("[greet][server.go][main()] Starting server ...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("[greet][server.go][main()] => Error encountered when invoking net.Listen(): %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[greet][server.go][main()] => Error encountered when invoking Server.Serve(): %v", err)
	}
}
