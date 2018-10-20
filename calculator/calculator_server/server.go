package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/LeeTrent/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("\n[Calculator][server.go][(*server)Sum] => *calculatorpb.SumRequest: %v", req)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber
	response := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	fmt.Printf("\n[Calculator][server.go][(*server)Sum] => *calculatorpb.SumResponse: %v", response)
	return response, nil
}

func main() {
	fmt.Println("[Calculator][server.go][main()] => LISTENING ...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("FATAL => [Calculator][server.go][main] => net.Listen(): %v\n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("FATAL => [Calculator][server.go][main] => grpc.Server.Serve(): %v\n", err)
	}
	fmt.Println("[Calculator][server.go][main()] => DONE")
}
