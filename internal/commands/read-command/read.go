package read_command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"htruong/kt-crawler/internal"

	"github.com/urfave/cli/v2"
)

// Category struct to parse _category_.json
type Category struct {
	Label string `json:"label"`
}

func RunRead(c *cli.Context, config *internal.Config) error {
	baseDir := config.Listeners.Store.Filesystem.Directory

	// 1. List available translations (directories in baseDir)
	translations, err := os.ReadDir(baseDir)
	if err != nil {
		// If the directory doesn't exist, it's not an error, just no content.
		if os.IsNotExist(err) {
			fmt.Printf("Base directory %s does not exist. Run 'scan' first.\n", baseDir)
			return nil
		}
		return cli.Exit(fmt.Sprintf("Failed to read base directory %s: %v", baseDir, err), 1)
	}

	var translationDirs []os.DirEntry
	for _, entry := range translations {
		if entry.IsDir() {
			translationDirs = append(translationDirs, entry)
		}
	}

	if len(translationDirs) == 0 {
		fmt.Printf("No translations found in %s. Run 'scan' first.\n", baseDir)
		return nil
	}

	var selection int
	if len(translationDirs) == 1 {
		selection = 1
	} else {
		// 2. First screen: list available translations
		fmt.Println("Available Translations:")
		for i, dir := range translationDirs {
			fmt.Printf("%d. %s\n", i+1, dir.Name())
		}

		// Get user selection
		fmt.Print("Enter selection number: ")
		_, err = fmt.Scanln(&selection)
		if err != nil || selection < 1 || selection > len(translationDirs) {
			return cli.Exit("Invalid selection", 1)
		}
	}

	selectedTranslation := translationDirs[selection-1].Name()
	translationPath := filepath.Join(baseDir, selectedTranslation)

	// 3. Second screen: list books
	books, err := os.ReadDir(translationPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read translation directory %s: %v", translationPath, err), 1)
	}

	fmt.Printf("\nBooks in %s:\n\n", selectedTranslation)
	var bookLabels []string
	for _, bookEntry := range books {
		if bookEntry.IsDir() {
			categoryPath := filepath.Join(translationPath, bookEntry.Name(), "_category_.json")

			// Read and parse _category_.json
			data, err := os.ReadFile(categoryPath)
			if err != nil {
				// Skip if _category_.json is not found or unreadable
				continue
			}

			var category Category
			if err := json.Unmarshal(data, &category); err != nil {
				// Skip if parsing fails
				continue
			}
			bookLabels = append(bookLabels, category.Label)
		}
	}

	if len(bookLabels) == 0 {
		fmt.Printf("No books found in %s\n", selectedTranslation)
		return nil
	}

	// Determine the number of columns and rows
	numBooks := len(bookLabels)
	numCols := 4
	numRows := (numBooks + numCols - 1) / numCols

	// Find the maximum width for each column
	colWidths := make([]int, numCols)
	for i := 0; i < numCols; i++ {
		for j := 0; j < numRows; j++ {
			index := j*numCols + i
			if index < numBooks {
				// Plus 4 for "xx. "
				width := len(bookLabels[index]) + 4
				if width > colWidths[i] {
					colWidths[i] = width
				}
			}
		}
	}

	// Print the books in columns
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			index := i + j*numRows
			if index < numBooks {
				label := fmt.Sprintf("%d. %s", index+1, bookLabels[index])
				fmt.Printf("% -*s", colWidths[j], label)
			}
		}
		fmt.Println()
	}
	fmt.Println()

	return nil
}
