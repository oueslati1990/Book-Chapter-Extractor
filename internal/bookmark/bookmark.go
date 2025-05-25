package bookmark

import (
	"os"
	"regexp"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type Bmk interface {
	ExtractBookmarks() ([]*pdfcpu.Bookmark, error)
	SanitizeFilename(title string) error
}

type Bookmark struct {
	Input string
}

func New(i string) *Bookmark {
	return &Bookmark{
		Input: i,
	}
}

func (b *Bookmark) ExtractBookmarks() ([]*pdfcpu.Bookmark, error) {
	file, err := os.Open(b.Input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Extract top-level bookmarks
	bookmarks, err := api.Bookmarks(file, nil)
	if err != nil {
		return nil, err
	}

	// Recursive flattening inside the same method to go throw Bookmark Children
	var flat []*pdfcpu.Bookmark
	var flatten func([]pdfcpu.Bookmark)
	flatten = func(nodes []pdfcpu.Bookmark) {
		for i := range nodes {
			bm := &nodes[i]
			flat = append(flat, bm)

			if len(bm.Kids) > 0 {
				flatten(bm.Kids)
			}
		}
	}

	flatten(bookmarks)
	return flat, nil
}

func (b *Bookmark) SanitizeFilename(title string) string {
	// Remove invalid characters
	name := strings.TrimSpace(title)
	reg := regexp.MustCompile(`[\\/:*?"<>|]`)
	name = reg.ReplaceAllString(name, "_")

	return name
}
