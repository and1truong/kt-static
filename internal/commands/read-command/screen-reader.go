package read_command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func ReadChapter(
	chapterPath string,
	translationLabel string,
	bookLabel string,
	chapterLabel string,
	totalChapters int,
) error {
	// Read and display chapter content
	contentBytes, err := os.ReadFile(chapterPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to read chapter file %s: %v", chapterPath, err), 1)
	}

	// Construct the full text with content and bottom status bar
	var buffer bytes.Buffer
	buffer.Write(contentBytes)

	// Status bar at the bottom: Translation; Book; Number of chapters; Current Chapter
	statusBar := fmt.Sprintf(
		"\n\n--- Status: Translation: %s | Book: %s | Chapters: %d | Current: %s ---",
		translationLabel,
		bookLabel,
		totalChapters,
		chapterLabel,
	)
	buffer.WriteString(statusBar)
	buffer.WriteString("\n") // Ensure a final newline

	// Use a pager to display the content
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less" // Default pager
	}

	cmd := exec.Command(pager)
	cmd.Stdin = &buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the pager command
	err = cmd.Run()
	if err != nil {
		// Fallback to simple print if pager fails
		fmt.Println(buffer.String())
	}

	return nil
}

func getChapterPath(bookPath string, chapterSelection int, chapterFiles []os.DirEntry) (string, error) {
	if chapterSelection < 1 || chapterSelection > len(chapterFiles) {
		return "", fmt.Errorf("invalid chapter selection")
	}

	selectedChapterFile := chapterFiles[chapterSelection-1]
	return filepath.Join(bookPath, selectedChapterFile.Name()), nil
}
