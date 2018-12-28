package main

import (
	"context"
	"fmt"
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
	fmt.Println("[blog][client][main] Creating the blog ...")
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
	fmt.Printf("\nCreated blogId: %v", blogID)

	////////////////////////////////////////////////////////////////////////////////
	// READ BLOG
	////////////////////////////////////////////////////////////////////////////////
	fmt.Printf("\nNow reading the blog we just created ...")

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

	fmt.Println("[blog][client][main] ... END")
}
