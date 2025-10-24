package read_command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// listChapters lists available chapters for a given book.
func listChapters(bookPath string, bookLabel string) error {
	chapters, err := os.ReadDir(bookPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read book directory %s: %v", bookPath, err), 1)
	}

	fmt.Printf("\nChapters in %s:\n\n", bookLabel)
	var chapterFiles []os.DirEntry
	for _, chapterEntry := range chapters {
		if chapterEntry.Name() != "_category_.json" {
			chapterFiles = append(chapterFiles, chapterEntry)
		}
	}

	if len(chapterFiles) == 0 {
		fmt.Printf("No chapters found in %s\n", bookLabel)
		return nil
	}

	sort.Slice(chapterFiles, func(i, j int) bool {
		nameI := chapterFiles[i].Name()
		nameJ := chapterFiles[j].Name()

		numStrI := strings.TrimSuffix(nameI, filepath.Ext(nameI))
		numStrJ := strings.TrimSuffix(nameJ, filepath.Ext(nameJ))

		var numI, numJ int
		var errI, errJ error

		numI, errI = strconv.Atoi(numStrI)
		numJ, errJ = strconv.Atoi(numStrJ)

		if errI != nil && errJ != nil {
			return nameI < nameJ // both non-numeric, sort alphabetically
		}
		if errI != nil {
			return true // i is non-numeric, j is numeric, non-numeric comes first
		}
		if errJ != nil {
			return false // j is non-numeric, i is numeric, non-numeric comes first
		}

		return numI < numJ // both are numeric
	})

	chapterLabels := make([]string, len(chapterFiles))
	for i, chapterEntry := range chapterFiles {
		name := chapterEntry.Name()
		numStr := strings.TrimSuffix(name, filepath.Ext(name))
		chapterLabels[i] = fmt.Sprintf("Chapter %s", numStr)
	}

	if err := formatInColumns(chapterLabels, 4); err != nil {
		return err
	}

	// Get user selection for chapter
	fmt.Print("Enter chapter number to read: ")
	var chapterSelection int
	_, err = fmt.Scanln(&chapterSelection)
	if err != nil || chapterSelection < 1 || chapterSelection > len(chapterFiles) {
		return cli.Exit("Invalid chapter selection", 1)
	}

	selectedChapterFile := chapterFiles[chapterSelection-1]
	chapterPath := filepath.Join(bookPath, selectedChapterFile.Name())

	// Read and display chapter content
	content, err := os.ReadFile(chapterPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read chapter file %s: %v", chapterPath, err), 1)
	}

	fmt.Printf("\n--- %s - %s ---\n\n", bookLabel, chapterLabels[chapterSelection-1])
	fmt.Println(string(content))
	fmt.Printf("\n--- End of Chapter ---\n")

	return nil
}
