package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/LeeTrent/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("\n[Calculator][server.go][(*server)PrimeNumberDecomposition] => BEGIN")

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("\n[Calculator][server.go][(*server)PrimeNumberDecomposition] => Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("[Calculator][server.go][(*server)CalculateAverage()] => BEGIN ...")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// We have finished reading the stream
				average := float64(sum) / float64(count)
				return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
					Average: average,
				})
			} else {
				log.Fatalf("\n[Calculator][server.go][(*server)CalculateAverage()] => stream.Recv() error: %v", err)
			}
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("[Calculator][server.go][(*server)FindMaximum()] => BEGIN ...")
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				log.Fatalf("\n[Calculator][server.go][(*server)FindMaximum()] => stream.Recv() error: %v", err)
			}
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendErr != nil {
				log.Fatalf("\n[greet][server.go][(*server)FindMaximum()] => stream.Send() error: %v", err)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	fmt.Printf("[Calculator][server.go][(*server)SquareRoot()] => BEGIN ...")

	number := req.GetNumber()

	if number < 0 {
		return nil,
			status.Error(codes.InvalidArgument, fmt.Sprintf("Received negative number: %v", number))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
