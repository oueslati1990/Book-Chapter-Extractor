package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/oueslati1990/Book-Chapter-Extractor/internal/extractor"
)

func main() {
	var input string
	var outputDir string
	var pattern string
	var verbose bool

	flag.StringVar(&input, "input", "", "Input PDF file (required)")
	flag.StringVar(&input, "i", "", "Shorthand of --input")
	flag.StringVar(&outputDir, "output", "Chapters", "Output Directory")
	flag.StringVar(&outputDir, "o", "Chapters", "Shorthand of --output")
	flag.StringVar(&pattern, "pattern", "", "Pattern of chapter")
	flag.StringVar(&pattern, "p", "", "Shorthand of --pattern")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Shorthand of --verbose")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `pdf-chapter-extractor - Extract chapters from a PDF book into separate files
Usage:  pdf-chapter-extractor [options]
Options:`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if input == "" {
		fmt.Println("Error : input file required")
		os.Exit(1)
	}

	if err := os.MkdirAll(pattern, 0755); err != nil {
		fmt.Printf("Error creating output folder %v\n", err)
		os.Exit(1)
	}

	ext := extractor.New(input, outputDir, pattern, verbose)
	if err := ext.ExtractChapters(); err != nil {
		fmt.Printf("Can't extract chapters : %v\n", err)
	}
}
