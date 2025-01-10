package main

import (
	"os"
	"strings"
	"testing"
)

func TestWordcount(t *testing.T) {
	var data []byte
	var err error
	var content []byte

	t.Run("parseFile", func(t *testing.T) {
		data, err = os.ReadFile("README.md")
		if err != nil {
			t.Errorf("unable to read file %v", err.Error())
		}

		content = parseFile(data)
		if strings.Contains(string(content), "#") {
			t.Error("parseFile returned incorrect content")
		}
	})

	var paragraphs []string
	t.Run("extractText", func(t *testing.T) {
		paragraphs, err = extractText(content)
		if err != nil {
			t.Errorf("failed to parse html %v", err.Error())
		}

		for _, p := range paragraphs {
			if strings.Contains(p, "<") {
				t.Errorf("markup should contain not html tag delimeters but it does not %v", p)
			}
		}
	})

	t.Run("countWords", func(t *testing.T) {
		for _, p := range paragraphs {
			count := countWords(p)

			if count < 1 {
				t.Errorf("there should not be empty paragraphs but there are: %v", p)
			}
		}
	})

	t.Run("main", func(t *testing.T) {
		args := os.Args[:1]
		args = append(args, "-f", "README.md")
		err = app.Run(args)
		if err != nil {
			t.Errorf("execution should not have failed %v", err.Error())
		}
	})
}
