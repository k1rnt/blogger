package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/blogger/v3"
	"google.golang.org/api/option"

	"github.com/k1rnt/blogger/internal/convert"
	"github.com/k1rnt/blogger/internal/meta"
)

func main() {
	mdPath := flag.String("path", "", "Markdown file")
	publish := flag.Bool("publish", true, "Publish (true) or draft (false)")
	flag.Parse()

	if *mdPath == "" {
		fmt.Fprintln(os.Stderr, "path required")
		os.Exit(1)
	}
	mdBytes, _ := ioutil.ReadFile(*mdPath)

	var fm meta.FrontMatter
	body, _ := frontmatter.Parse(strings.NewReader(string(mdBytes)), &fm)

	html := convert.MarkdownToHTML(string(body), os.Getenv("RAW_BASE"))

	// OAuth client
	conf := &oauth2.Config{
		ClientID:     os.Getenv("BLOGGER_CLIENT_ID"),
		ClientSecret: os.Getenv("BLOGGER_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{blogger.BloggerScope},
	}
	tok := &oauth2.Token{RefreshToken: os.Getenv("BLOGGER_REFRESH_TOKEN")}
	client := conf.Client(context.Background(), tok)
	svc, _ := blogger.NewService(context.Background(), option.WithHTTPClient(client))

	post := &blogger.Post{
		Title:   fm.Title,
		Content: html,
		Labels:  fm.Labels,
	}

	if fm.BloggerID != "" {
		// 更新（posts.update）
		_, err := svc.Posts.Update(os.Getenv("BLOG_ID"), fm.BloggerID, post).
			IsDraft(!*publish).Do()
		must(err)
		fmt.Println("updated:", fm.BloggerID)
	} else {
		// 新規
		res, err := svc.Posts.Insert(os.Getenv("BLOG_ID"), post).
			IsDraft(!*publish).Do()
		must(err)
		fmt.Println("inserted:", res.Id)
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
