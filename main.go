package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// load markdown file
	mdFile, err := ioutil.ReadFile("./example.md")
	if err != nil {
		log.Fatal(err)
	}
	// convert markdown to html
	html := blackfriday.MarkdownCommon(mdFile)
	// replace code-parts with syntax-highlighted parts
	replaced, err := replaceCodeParts(html)
	if err != nil {
		log.Fatal(err)
	}
	// write to stdout
	t, err := template.ParseFiles("./template.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, struct{ Content string }{Content: replaced})
	if err != nil {
		log.Fatal(err)
	}
}

func replaceCodeParts(mdFile []byte) (string, error) {
	byteReader := bytes.NewReader(mdFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", err
	}
	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		oldCode := s.Text()
		formatted, err := syntaxhighlight.AsHTML([]byte(oldCode))
		if err != nil {
			log.Fatal(err)
		}
		s.SetHtml(string(formatted))
	})
	new, err := doc.Html()
	if err != nil {
		return "", err
	}
	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}
