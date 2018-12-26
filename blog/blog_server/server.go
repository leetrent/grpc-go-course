package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LeeTrent/grpc-go-course/blog/blogpb"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("[blog][server][CreateBlog]: BEGIN")

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetAuthorId(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("[blog][server][ReadBlog] => BEGIN")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("[blog][server][ReadBlog] => primitive.ObjectIDFromHex(blogID): %v", err),
		)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("[blog][server][ReadBlog] => collection.FindOne(): %v", err),
		)
	}

	fmt.Println("[blog][server][ReadBlog] => END (returning Blog")

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPB(data),
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("[blog][server][UpdateBlog] => BEGIN")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("[blog][server][UpdateBlog] => primitive.ObjectIDFromHex(blogID): %v", err),
		)
	}
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("[blog][server][UpdateBlog] => collection.FindOne(): %v", err),
		)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("[blog][server][UpdateBlog] => Cannot update object in MongoDB: %v", updateErr),
		)
	}

	func dataToBlogPB(data *blogItem) *blogpb.Blog {

		return &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		}
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPB(data),
	}, nil
}

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
