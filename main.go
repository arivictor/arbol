package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var verbose bool

func formatMetadata(file os.FileInfo) string {
	modTime := file.ModTime().Format(time.RFC3339)
	size := file.Size()
	return fmt.Sprintf("(size: %d bytes, modified: %s)", size, modTime)
}

func printDir(path string, prefix string, parentPrefix string) {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open directory '%s': %s\n", path, err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		fmt.Printf("Failed to read directory '%s': %s\n", path, err)
		return
	}

	for i, file := range files {
		newPrefix := parentPrefix
		visual := "├── "

		if i == len(files)-1 {
			visual = "└── "
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		displayName := file.Name()
		if file.IsDir() {
			displayName += "/"
		}
		if verbose {
			displayName += " " + formatMetadata(file)
		}

		colorCode := ""
		resetCode := ""
		if file.IsDir() {
			colorCode = "\x1b[34m"
			resetCode = "\x1b[0m"
		}

		fmt.Printf("%s%s%s%s%s\n", prefix, visual, colorCode, displayName, resetCode)

		if file.IsDir() {
			printDir(filepath.Join(path, file.Name()), newPrefix, newPrefix)
		}
	}
}

func main() {
	// Parse command-line flags
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output (shorthand)")
	flag.Parse()

	// Get the directory to start from
	startDir := "."
	if len(flag.Args()) > 0 {
		startDir = flag.Args()[0]
	}

	fmt.Println(startDir)
	printDir(startDir, "", "")
}
