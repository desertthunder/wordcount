// package wordcount is a simple command line tool that
// prints out the number of words in the prose sections
// of a markdown file.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/urfave/cli/v2"

	xHTML "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var logger = slog.Default()

var app *cli.App = &cli.App{
	Name:     "wordcount",
	Version:  "v1.0.0",
	Compiled: time.Now(),
	Authors: []*cli.Author{{
		Name:  "Owais J.",
		Email: "desertthunder.dev@gmail.com.com",
	}},
	Copyright: "(c) 2025 lol",
	HelpName:  "wordcount",
	Usage:     "command-line wordcounter for markdown files",
	UsageText: "wordcount - how to input markdown",
	ArgsUsage: "[args]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f", "path", "p"},
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		fpath := ctx.String("file")
		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}

		markup := parseFile(data)
		paragraphs, err := extractText(markup)

		fmt.Printf("%v Contents \n%v", fpath, string(markup))

		if err != nil {
			return err
		}

		total := int32(0)
		for _, p := range paragraphs {
			total = total + countWords(p)
		}

		logger.Info(fmt.Sprintf(
			"there are %v words among %v paragraphs in your document",
			total, len(paragraphs),
		))

		return nil
	},
}

func parseFile(c []byte) []byte {
	p := parser.NewWithExtensions(parser.CommonExtensions)
	ast := markdown.Parse(c, p)
	html := markdown.Render(ast, &html.Renderer{})
	return html
}

func extractText(m []byte) ([]string, error) {
	s := string(m)
	paragraphs := []string{}

	doc, err := xHTML.Parse(strings.NewReader(s))
	if err != nil {
		return paragraphs, err
	}

	for n := range doc.Descendants() {
		if n.Type ==
			xHTML.ElementNode {
			if n.DataAtom == atom.P ||
				n.Data == "p" {
				txt := n.FirstChild.Data
				paragraphs = append(paragraphs, txt)
			}
		}
	}

	logger.Debug(fmt.Sprintf("there are %v paragraphs in the source file",
		len(paragraphs)))

	return paragraphs, nil
}

func countWords(p string) int32 {
	var by_space []string
	by_punctuation := strings.Split(p, ".")
	words := []string{}
	for _, sentence := range by_punctuation {
		by_space = strings.Split(sentence, " ")
		words = append(words, by_space...)
	}

	return int32(len(words))
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(fmt.Sprintf("execution failed %v", err.Error()))
	}
}
