package read_command

import (
	"fmt"
	"htruong/kt-crawler/internal"

	"github.com/urfave/cli/v2"
)

func RunRead(c *cli.Context, config *internal.Config) error {
	baseDir := config.Listeners.Store.Filesystem.Directory

	// 1. Select Translation
	selectedTranslation, translationPath, err := screenSelectTranslation(baseDir)
	if err != nil {
		return err
	}
	if selectedTranslation == "" {
		return nil // No translations found, already printed message
	}

	// 2. Select Book
	selectedBook, bookPath, err := screenSelectBook(translationPath, selectedTranslation)
	if err != nil {
		return err
	}
	if selectedBook == nil {
		return nil // No books found, already printed message
	}

	// 3. List Chapters
	return screenSelectChapter(bookPath, selectedBook.Label, selectedTranslation)
}

// formatInColumns prints a list of labels in the specified number of columns.
func formatInColumns(labels []string, numCols int) error {
	numItems := len(labels)
	if numItems == 0 {
		return nil
	}

	numRows := (numItems + numCols - 1) / numCols

	// Find the maximum width for each column
	colWidths := make([]int, numCols)
	for i := 0; i < numCols; i++ {
		for j := 0; j < numRows; j++ {
			index := j*numCols + i
			if index < numItems {
				// Calculate width including index and padding
				label := fmt.Sprintf("%d. %s", index+1, labels[index])
				width := len(label) + 3 // Extra padding
				if width > colWidths[i] {
					colWidths[i] = width
				}
			}
		}
	}

	// Print the items in columns
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			index := i + j*numRows
			if index < numItems {
				label := fmt.Sprintf("%d. %s", index+1, labels[index])
				fmt.Printf("% -*s", colWidths[j], label)
			}
		}
		fmt.Println()
	}
	fmt.Println()
	return nil
}
