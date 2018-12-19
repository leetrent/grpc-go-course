package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/LeeTrent/grpc-go-course/blog/blogpb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type blogItem struct{}

func main() {
	fmt.Println("[blog][server][main][main()]: BEGIN ...")

	///////////////////////////////////////////////////
	// This will provide the file name and line number
	// if our Go code crashes
	///////////////////////////////////////////////////
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	///////////////////////////////////////////////////

	///////////////////////////////////////////////////
	// Connect to MongoDB
	///////////////////////////////////////////////////
	fmt.Println("Connecting to MongoDB ...")
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("grpc-go-course").Collection("blog")

	///////////////////////////////////////////////////
	// Open the Blog Listener
	///////////////////////////////////////////////////
	fmt.Println("Opening the blog listener ...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("[blog][server][main][main()]: %v", err)
	}

	///////////////////////////////////////////////////
	// Start the Blog Server
	///////////////////////////////////////////////////
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting blog server ...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[blog][server][main][main()]: %v", err)
		}
	}()

	// Wait for 'Control C' to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	fmt.Println("Stopping the blog server ...")
	s.Stop()

	fmt.Println("Closing the blog listener ...")
	lis.Close()

	fmt.Println("Closing MongoDB Connection ...")
	client.Disconnect(context.TODO())

	fmt.Println("[blog][server][main][main()]: ... END")
}
