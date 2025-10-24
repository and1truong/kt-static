package read_command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// screenSelectTranslation lists available translations and prompts the user for a selection.
func screenSelectTranslation(baseDir string) (string, string, error) {
	translations, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Base directory %s does not exist. Run 'scan' first.\n", baseDir)
			return "", "", nil
		}
		return "", "", cli.Exit(fmt.Sprintf("Failed to read base directory %s: %v", baseDir, err), 1)
	}

	var translationDirs []os.DirEntry
	for _, entry := range translations {
		if entry.IsDir() {
			translationDirs = append(translationDirs, entry)
		}
	}

	if len(translationDirs) == 0 {
		fmt.Printf("No translations found in %s. Run 'scan' first.\n", baseDir)
		return "", "", nil
	}

	var selection int
	if len(translationDirs) == 1 {
		selection = 1
	} else {
		fmt.Println("Available Translations:")
		for i, dir := range translationDirs {
			fmt.Printf("%d. %s\n", i+1, dir.Name())
		}

		fmt.Print("Enter selection number: ")
		_, err = fmt.Scanln(&selection)
		if err != nil || selection < 1 || selection > len(translationDirs) {
			return "", "", cli.Exit("Invalid selection", 1)
		}
	}

	selectedTranslation := translationDirs[selection-1].Name()
	translationPath := filepath.Join(baseDir, selectedTranslation)
	return selectedTranslation, translationPath, nil
}
