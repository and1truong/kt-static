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
func listChapters(bookPath string, bookLabel string, translationLabel string) error {
	chapters, err := os.ReadDir(bookPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read book directory %s: %v", bookPath, err), 1)
	}

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

		numI, errI := strconv.Atoi(numStrI)
		numJ, errJ := strconv.Atoi(numStrJ)

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

	for {
		fmt.Printf("\nChapters in %s:\n\n", bookLabel)
		if err := formatInColumns(chapterLabels, 4); err != nil {
			return err
		}

		// Get user selection for chapter
		fmt.Print("Enter chapter number to read (or 'q' to quit): ")
		var input string
		_, err = fmt.Scanln(&input)
		if err != nil {
			return cli.Exit("Invalid input", 1)
		}

		if strings.ToLower(input) == "q" {
			return nil
		}

		chapterSelection, err := strconv.Atoi(input)
		if err != nil || chapterSelection < 1 || chapterSelection > len(chapterFiles) {
			fmt.Println("Invalid chapter selection, please try again.")
			continue
		}

		chapterPath, err := getChapterPath(bookPath, chapterSelection, chapterFiles)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Failed to get chapter path: %v", err), 1)
		}

		err = ReadChapter(
			chapterPath,
			translationLabel,
			bookLabel,
			chapterLabels[chapterSelection-1],
			len(chapterFiles),
		)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Failed to read chapter: %v", err), 1)
		}
	}
}
