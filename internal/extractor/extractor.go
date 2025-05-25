package extractor

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/oueslati1990/Book-Chapter-Extractor/internal/bookmark"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type Extractor interface {
	ExtractChapters() error
}

type PdfExtractor struct {
	Input   string
	Output  string
	Pattern string
	Verbose bool
}

func New(i, o, p string, v bool) *PdfExtractor {
	return &PdfExtractor{
		Input:   i,
		Output:  o,
		Pattern: p,
		Verbose: v,
	}
}

func (e *PdfExtractor) ExtractChapters() error {
	if e.Verbose {
		fmt.Printf("Processing PDF : %s\n", e.Input)
	}

	b := bookmark.New(e.Input)
	// Get bookmarks from PDF
	bookmarks, err := b.ExtractBookmarks()
	if err != nil {
		return fmt.Errorf("failed to extract bookmarks : %v", err)
	}

	if len(bookmarks) == 0 {
		return fmt.Errorf("no bookmarks found in the PDF file")
	}

	//Find the specific bookmark that matches the pattern
	var targetBookmark *pdfcpu.Bookmark
	var matches []*pdfcpu.Bookmark
	if e.Pattern != "" {
		pattern, err := regexp.Compile(e.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern : %v", err)
		}
		for _, bookmark := range bookmarks {
			if pattern.MatchString(bookmark.Title) {
				matches = append(matches, bookmark)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("no bookmark matching the pattern")
		}

		if len(matches) > 1 {
			var titles []string
			for _, bm := range matches {
				titles = append(titles, bm.Title)
			}
			return fmt.Errorf("multiple bookmarks matched the pattern , please be more specific, Maches : %s", strings.Join(titles, "\n-"))
		}

		targetBookmark = matches[0]
	} else {
		return fmt.Errorf("no pattern provided for bookmark matching")
	}

	// Process just the bookmark chapter
	startPage := targetBookmark.PageFrom
	endPage := targetBookmark.PageThru

	if e.Verbose {
		fmt.Printf("Extracting pages %d-%d \n", startPage, endPage)
	}

	// Clean up title for filename
	chapterTitle := b.SanitizeFilename(targetBookmark.Title)
	if chapterTitle == "" {
		chapterTitle = "Extracted_Chapter"
	}

	outputFile := filepath.Join(e.Output, fmt.Sprintf("%s.pdf", chapterTitle))

	if e.Verbose {
		fmt.Printf("Extracting chapter : %s (pages %d-%d) to %s\n",
			targetBookmark.Title, startPage, endPage, outputFile)
	}
	// Extract pages to new PDF
	if err := api.TrimFile(e.Input, outputFile, []string{fmt.Sprintf("%d-%d", startPage, endPage)}, nil); err != nil {
		return fmt.Errorf("failed to extract chapter %s : %v", targetBookmark.Title, err)
	}

	fmt.Printf("Sucessfully extracted chapter '%s' to %s\n", targetBookmark.Title, outputFile)
	return nil
}
