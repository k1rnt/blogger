// cmd/export/main.go
package main

import (
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/blogger/v3"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"

	"github.com/k1rnt/blogger/internal/convert"
	"github.com/k1rnt/blogger/internal/meta"
)

var (
	imgTagRE   = regexp.MustCompile(`(?i)<img[^>]+src="([^"]+)"`)
	httpClient = &http.Client{Timeout: 20 * time.Second}
)

func main() {
	outDir := flag.String("out", "posts", "output directory for Markdown")
	dlImg := flag.Bool("dlimg", true, "download images & rewrite src")
	flag.Parse()

	conf := &oauth2.Config{
		ClientID:     mustEnv("BLOGGER_CLIENT_ID"),
		ClientSecret: mustEnv("BLOGGER_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{blogger.BloggerReadonlyScope},
	}
	tok := &oauth2.Token{RefreshToken: mustEnv("BLOGGER_REFRESH_TOKEN")}
	svc, _ := blogger.NewService(context.Background(),
		option.WithHTTPClient(conf.Client(context.Background(), tok)))

	call := svc.Posts.List(mustEnv("BLOG_ID")).MaxResults(500)
	for {
		resp, err := call.Do()
		if err != nil {
			log.Fatalf("Blogger API error: %v", err)
		}
		for _, p := range resp.Items {
			html := p.Content
			var imgMap map[string]string

			if *dlImg {
				imgMap = downloadImages(html)
				for orig, local := range imgMap {
					html = strings.ReplaceAll(html, orig, "images/"+local)
				}
			}

			md, _ := convert.HTMLToMarkdown(html)

			date, _ := time.Parse(time.RFC3339, p.Published)
			slug := convert.Slugify(p.Title)
			fn := fmt.Sprintf("%s/%s-%s.md", *outDir, date.Format("2006-01-02"), slug)

			fm := meta.FrontMatter{
				Title:     p.Title,
				Labels:    p.Labels,
				BloggerID: p.Id,
				Slug:      slug,
				Status:    "publish",
				Published: p.Published,
			}
			writeMarkdown(fn, fm, md)
			fmt.Println("exported:", fn)
		}
		if resp.NextPageToken == "" {
			break
		}
		call.PageToken(resp.NextPageToken)
	}
}

// downloadImages() returns map[originalURL]localFilename
func downloadImages(html string) map[string]string {
	if err := os.MkdirAll("images", 0o755); err != nil {
		panic(err)
	}
	imgMap := make(map[string]string)
	matches := imgTagRE.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		orig := m[1]
		if _, done := imgMap[orig]; done {
			continue
		}
		local := saveImage(orig)
		if local != "" {
			imgMap[orig] = local
		}
	}
	return imgMap
}

func saveImage(url string) string {
	resp, err := httpClient.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "skip img:", url, err)
		return ""
	}
	defer resp.Body.Close()

	ext := filepath.Ext(url)
	if len(ext) > 10 || len(ext) == 0 {
		ext = ".jpg"
	}

	h := sha1.Sum([]byte(url))
	local := fmt.Sprintf("%x%s", h[:6], ext) // 12 hex = 6byte
	fpath := filepath.Join("images", local)

	out, _ := os.Create(fpath)
	_, _ = io.Copy(out, resp.Body)
	out.Close()
	return local
}

func writeMarkdown(fn string, fm meta.FrontMatter, md string) {
	if err := os.MkdirAll(filepath.Dir(fn), 0o755); err != nil {
		panic(err)
	}
	b, _ := yaml.Marshal(fm)
	content := fmt.Sprintf("---\n%s---\n%s", b, md)
	if err := os.WriteFile(fn, []byte(content), 0o644); err != nil {
		panic(err)
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(fmt.Sprintf("env %s is required", k))
	}
	return v
}
