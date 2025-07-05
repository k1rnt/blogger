package convert

import (
	"regexp"
	"strings"
	"unicode"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/text/unicode/norm"
)

var imgRE = regexp.MustCompile(`src="images/([^"]+)"`)

func MarkdownToHTML(md, rawBase string) string {
	var buf strings.Builder
	gm := goldmark.New(goldmark.WithExtensions(extension.GFM))
	_ = gm.Convert([]byte(md), &buf)
	return imgRE.ReplaceAllString(buf.String(), `src="`+rawBase+`/images/$1"`)
}

func HTMLToMarkdown(html string) (string, error) {
	conv := h2m.NewConverter("", true, nil)
	return conv.ConvertString(html)
}

func Slugify(in string) string {
	s := strings.ToLower(norm.NFKD.String(in))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(strings.Join(strings.FieldsFunc(b.String(), func(r rune) bool { return r == '-' }), "-"), "-")
}
