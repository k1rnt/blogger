// cmd/delete/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/blogger/v3"
	"google.golang.org/api/option"
)

func main() {
	postID := flag.String("id", "", "Blogger postID to delete permanently")
	flag.Parse()

	if *postID == "" {
		log.Fatal("postID is required")
	}

	// OAuth setup
	conf := &oauth2.Config{
		ClientID:     os.Getenv("BLOGGER_CLIENT_ID"),
		ClientSecret: os.Getenv("BLOGGER_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{blogger.BloggerScope},
	}
	tok := &oauth2.Token{RefreshToken: os.Getenv("BLOGGER_REFRESH_TOKEN")}
	client := conf.Client(context.Background(), tok)
	svc, _ := blogger.NewService(context.Background(), option.WithHTTPClient(client))

	// Permanent delete
	blogID := os.Getenv("BLOG_ID")
	if err := svc.Posts.Delete(blogID, *postID).Do(); err != nil {
		log.Fatalf("delete %s: %v", *postID, err)
	}
	log.Println("deleted:", *postID)
}
