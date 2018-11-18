package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/LeeTrent/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("[Calculator][client.go][main] => BEGIN")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("FATAL => [Calculator][client.go][main] => grpc.Dial: Error: %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(client)
	//doServerStreaming(client)
	//doClientStreaming(client)
	doBiDiStreaming(client)

	fmt.Println("[Calculator][client.go][main] => END")
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("[Calculator][client.go][doUnary] => BEGIN")

	request := &calculatorpb.SumRequest{
		FirstNumber:  54,
		SecondNumber: 64,
	}

	response, err := client.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("FATAL => [Calculator][client.go][doUnary] => %v: ", err)
	}

	log.Printf("[Calculator][client.go][doUnary] => calculatorpb.SumResponse: %v", response.SumResult)

	fmt.Println("[Calculator][client.go][doUnary] => END")
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("[Calculator][client.go][doServerStreaming] => BEGIN")

	request := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12390392840,
	}

	stream, err := client.PrimeNumberDecomposition(context.Background(), request)
	if err != nil {
		log.Fatalf("FATAL => [Calculator][client.go][doServerStreaming] => %v: ", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("\nFATAl => [Calculator][client.go][doServerStreaming] => calculatorpb.CalculatorServiceClient.PrimeNumberDecomposition.Recv(): Error: %v", err)
			}
		}
		log.Printf("[Calculator][client.go][doServerStreaming] => calculatorpb.PrimeNumberDecompositionResponse.GetPrimeFactor(): %v", resp.GetPrimeFactor())
	}

	fmt.Println("[Calculator][client.go][doServerStreaming] => END")
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("[Calculator][client.go][doClientStreaming] => BEGIN")

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("\n[Calculator][client.go][doClientStreaming] => Error encountered when invoking client.CalculateAverage(): %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("\n[Calculator][client.go][doClientStreaming] => Error encountered when invoking stream.CloseAndRecv(): %v", err)
	}

	fmt.Printf("\n[Calculator][client.go][doClientStreaming] => CalculateAverageResponse.getAverage(): %v\n", response.GetAverage())
}

func doBiDiStreaming(client calculatorpb.CalculatorServiceClient) {

	fmt.Println("[Calculator][client.go][doBiDiStreaming] => BEGIN")

	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("\n[Calculator][client.go][doBiDiStreaming] => Error encountered when invoking client.FindMaximum(): %v", err)
	}

	waitChannel := make(chan struct{})

	// Send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive go routine
	go func() {
		for {
			respone, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatalf("\n[Calculator][client.go][doBiDiStreaming] => Error encountered when invoking stream.Recv(): %v", err)
					break
				}
			}
			maximum := respone.GetMaximum()
			fmt.Printf("New maximum is: %v:\n", maximum)
		}
		close(waitChannel)
	}()
	<-waitChannel

}
