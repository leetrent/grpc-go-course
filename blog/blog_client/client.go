package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/LeeTrent/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("[blog][client][main] BEGIN ...")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("[blog][client][main][grpc.Dial]: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	////////////////////////////////////////////////////////////////////////////////
	// CREATE BLOG
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	fmt.Println("CREATE BLOG")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")

	blog := &blogpb.Blog{
		AuthorId: "Lee",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	cbr, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("[blog][client][main][client.CreateBlog]: %v", err)
	}
	fmt.Printf("Blog has been created:\n%v\n", cbr)
	blogID := cbr.GetBlog().GetId()
	fmt.Printf("\nCreated blogId: %v\n", blogID)

	////////////////////////////////////////////////////////////////////////////////
	// READ BLOG
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n////////////////////////////////////////////////////////////////////////////////")
	fmt.Println("READ BLOG")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	fmt.Printf("Now reading the blog we just created ...")

	_, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "7c1fcc269e28f6b46666838f"})
	if readErr != nil {
		fmt.Printf("\nError happened while reading: %v", readErr)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)

	if readBlogErr != nil {
		fmt.Printf("\nError happened while reading: %v \n", readBlogErr)
	}

	fmt.Printf("\nBlog was read: %v \n", readBlogRes)

	////////////////////////////////////////////////////////////////////////////////
	// UPDATE BLOG
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n////////////////////////////////////////////////////////////////////////////////")
	fmt.Println("UPDATE BLOG")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")

	updatedBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: updatedBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	////////////////////////////////////////////////////////////////////////////////
	// DELETE BLOG
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n////////////////////////////////////////////////////////////////////////////////")
	fmt.Println("DELETE BLOG")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")

	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", updateErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	////////////////////////////////////////////////////////////////////////////////
	// LIST ALL BLOGS
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n////////////////////////////////////////////////////////////////////////////////")
	fmt.Println("LIST ALL BLOGS")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected error encountered when iterating over stream result: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
	fmt.Println("\n[blog][client][main] ... END")
}
