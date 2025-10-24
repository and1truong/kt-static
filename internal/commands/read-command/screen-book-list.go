package read_command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/urfave/cli/v2"
)

type (
	// bookCategory struct to parse _category_.json
	bookCategory struct {
		Label    string `json:"label"`
		Position int    `json:"position"`
		Weight   int    `json:"weight"`
	}

	bookInfo struct {
		Label    string
		DirEntry os.DirEntry
		Position int
	}
)

// screenSelectBook lists available books for a given translation and prompts the user for a selection.
func screenSelectBook(translationPath string, selectedTranslation string) (*bookInfo, string, error) {
	books, err := os.ReadDir(translationPath)
	if err != nil {
		return nil, "", cli.Exit(fmt.Sprintf("Failed to read translation directory %s: %v", translationPath, err), 1)
	}

	fmt.Printf("\nBooks in %s:\n\n", selectedTranslation)

	var availableBooks []bookInfo
	for _, bookEntry := range books {
		if bookEntry.IsDir() {
			categoryPath := filepath.Join(translationPath, bookEntry.Name(), "_category_.json")

			data, err := os.ReadFile(categoryPath)
			if err != nil {
				// Skip if _category_.json is not found or unreadable
				continue
			}

			var category bookCategory
			if err := json.Unmarshal(data, &category); err != nil {
				// Skip if parsing fails
				continue
			}
			availableBooks = append(availableBooks, bookInfo{
				Label:    category.Label,
				DirEntry: bookEntry,
				Position: category.Position,
			})
		}
	}

	if len(availableBooks) == 0 {
		fmt.Printf("No books found in %s\n", selectedTranslation)
		return nil, "", nil
	}

	sort.Slice(availableBooks, func(i, j int) bool {
		return availableBooks[i].Position < availableBooks[j].Position
	})

	bookLabels := make([]string, len(availableBooks))
	for i, book := range availableBooks {
		bookLabels[i] = book.Label
	}

	if err := formatInColumns(bookLabels, 4); err != nil {
		return nil, "", err
	}

	// Get user selection for book
	fmt.Print("Enter book number to read: ")
	var bookSelection int
	_, err = fmt.Scanln(&bookSelection)
	if err != nil || bookSelection < 1 || bookSelection > len(availableBooks) {
		return nil, "", cli.Exit("Invalid book selection", 1)
	}

	selectedBook := availableBooks[bookSelection-1]
	bookPath := filepath.Join(translationPath, selectedBook.DirEntry.Name())
	return &selectedBook, bookPath, nil
}
