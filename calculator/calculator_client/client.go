package main

import (
	"context"
	"fmt"
	"log"

	"github.com/simplesteph/grpc-go-course/calculator/calculatorpb"

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

	doUnary(client)
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
