package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/blogger/v3"
	"google.golang.org/api/option"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type FrontMatter struct {
	Title     string   `yaml:"title"`
	Labels    []string `yaml:"labels"`
	BloggerID string   `yaml:"blogger_id"`
	Slug      string   `yaml:"slug,omitempty"`
	Status    string   `yaml:"status,omitempty"`
	Published string   `yaml:"published,omitempty"`
}

func markdownToHTML(md, rawBase string) string {
	var buf strings.Builder
	gm := goldmark.New(goldmark.WithExtensions(extension.GFM))
	_ = gm.Convert([]byte(md), &buf)
	// 画像パスを GitHub Raw URL に変換
	return strings.ReplaceAll(buf.String(), `src="images/`, `src="`+rawBase+`/images/`)
}

func main() {
	mdPath := flag.String("path", "", "Markdown file to publish")
	publish := flag.Bool("publish", true, "true: publish / false: draft")
	flag.Parse()

	if *mdPath == "" {
		log.Fatal("path is required")
	}

	mdBytes, err := ioutil.ReadFile(*mdPath)
	if err != nil {
		log.Fatalf("read %s: %v", *mdPath, err)
	}

	var fm FrontMatter
	body, err := frontmatter.Parse(strings.NewReader(string(mdBytes)), &fm)
	if err != nil {
		log.Fatalf("front-matter parse error in %s: %v", *mdPath, err)
	}

	html := markdownToHTML(string(body), os.Getenv("RAW_BASE"))

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

	blogID := os.Getenv("BLOG_ID")
	if fm.BloggerID != "" {
		// ---- update existing ----
		call := svc.Posts.Patch(blogID, fm.BloggerID, post)
		if *publish {
			call = call.Publish(true)
		}
		if _, err := call.Do(); err != nil {
			log.Fatalf("update %s: %v", *mdPath, err)
		}
		fmt.Println("updated:", fm.BloggerID)

	} else {
		// ---- insert new ----
		res, err := svc.Posts.Insert(blogID, post).
			IsDraft(!*publish).Do()
		if err != nil {
			log.Fatalf("insert %s: %v", *mdPath, err)
		}
		fmt.Println("inserted:", res.Id)
	}
}
