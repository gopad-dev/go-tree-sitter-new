package highlight

import (
	"context"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/tree-sitter/tree-sitter-go/bindings/go"

	"github.com/tree-sitter/go-tree-sitter"
)

const (
	ResetStyle   = "\x1b[m"
	GreenStyle   = "\x1b[32m"
	BlueStyle    = "\x1b[34m"
	MagentaStyle = "\x1b[35m"
)

var (
	theme = map[string]string{
		"keyword":  MagentaStyle,
		"variable": BlueStyle,
		"string":   GreenStyle,
	}

	captureNames = []string{
		"keyword",
		"variable",
		"string",
	}

	source = []byte(`package main

func main() {
	println("Hello, World!")
}
`)
)

func TestHighlighter_Highlight(t *testing.T) {
	highlightsQuery, err := os.ReadFile("test_data/highlights.scm")
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())

	cfg, err := NewConfiguration(language, "go", highlightsQuery, nil, nil)
	if err != nil {
		log.Fatalf("failed to create highlight config: %v", err)
	}
	if cfg == nil {
		log.Fatalf("tree-sitter grammar not found")
	}

	cfg.Configure(captureNames)

	highlighter := New()
	highlights := highlighter.Highlight(context.Background(), *cfg, source, func(name string) *Configuration {
		log.Println("loading highlight config for", name)
		return nil
	})

	var (
		highlightNames []string
		highlightName  *string
	)
	for event, err := range highlights {
		if err != nil {
			log.Panicf("failed to highlight source: %v", err)
		}
		switch e := event.(type) {
		case EventSource:
			var style string
			if highlightName != nil {
				style, _ = theme[*highlightName]
			}
			renderStyle(style, string(source[e.Start:e.End]))

		case EventStart:
			highlightName = &captureNames[e.Highlight]
			if !slices.Contains(highlightNames, *highlightName) {
				highlightNames = append(highlightNames, *highlightName)
			}
		case EventEnd:
			highlightName = nil
		}
	}

	t.Logf("used highlight names: %v", highlightNames)
}

func renderStyle(style string, source string) {
	print(style + source + ResetStyle)
}
