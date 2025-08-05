package data

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Library struct {
	db *SQLite
}

func (l *Library) GetBooks() []string {
	files, err := os.ReadDir("./data/library")
	if err != nil {
		return nil
	}

	var books []string
	for _, file := range files {
		if !file.IsDir() {
			books = append(books, file.Name())
		}
	}

	return books
}

func (l *Library) GetChapter(name, part, chapter string) (string, error) {
	filePath := fmt.Sprintf("./data/library/%s/content/%s/chapter_%s.md", name, part, chapter)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Process markdown content to fix line break issues
	processedContent := l.processMarkdownLineBreaks(string(content))

	return processedContent, nil
}

// processMarkdownLineBreaks fixes the rendering issue where single \n characters
// create unnecessary line breaks in the final HTML output, making text appear choppy.
// This function preserves intentional breaks (after headings, references, paragraphs)
// while joining broken sentences that were split across multiple lines.
func (l *Library) processMarkdownLineBreaks(content string) string {
	// Step 1: Protect intentional paragraph breaks (\n\n)
	content = strings.ReplaceAll(content, "\n\n", "|||PARAGRAPH_BREAK|||")

	// Step 2: Protect line breaks after markdown structural elements
	// These patterns preserve the formatting structure of the document

	// 2a. Headings followed by reference citations: "## Title\n(reference)\n"
	headingWithRefPattern := regexp.MustCompile(`(#{1,6}[^\n]+)\n(\([^\n]+\)\.?)\n`)
	content = headingWithRefPattern.ReplaceAllString(content, "${1}|||HEADING_BREAK|||${2}|||REFERENCE_BREAK|||")

	// 2b. Standalone headings: "## Title\n"
	headingPattern := regexp.MustCompile(`(#{1,6}[^\n]+)\n`)
	content = headingPattern.ReplaceAllString(content, "${1}|||HEADING_BREAK|||")

	// 2c. Lines ending with colon (typically section introducers): "Text:\n"
	colonPattern := regexp.MustCompile(`([^\n]+:)\n`)
	content = colonPattern.ReplaceAllString(content, "${1}|||COLON_BREAK|||")

	// 2d. Standalone reference lines: "(reference)\n"
	parenthesesPattern := regexp.MustCompile(`(\([^\n]+\)\.?)\n`)
	content = parenthesesPattern.ReplaceAllString(content, "${1}|||REFERENCE_BREAK|||")

	// Step 3: Replace remaining single line breaks with spaces
	// This joins broken sentences that were unnecessarily split across lines
	content = strings.ReplaceAll(content, "\n", " ")

	// Step 4: Restore all protected line breaks
	// Order matters here to ensure proper reconstruction
	replacements := map[string]string{
		"|||HEADING_BREAK|||":   "\n",
		"|||COLON_BREAK|||":     "\n",
		"|||REFERENCE_BREAK|||": "\n",
		"|||PARAGRAPH_BREAK|||": "\n\n",
	}

	for placeholder, replacement := range replacements {
		content = strings.ReplaceAll(content, placeholder, replacement)
	}

	return content
}
