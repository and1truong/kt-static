package read_command

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func translationListEnter(m modal) modal {
	if m.translations == nil {
		entries, err := os.ReadDir(m.baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				m.err = fmt.Errorf("base directory %s does not exist. Run 'scan' first", m.baseDir)
			} else {
				m.err = fmt.Errorf("failed to read base directory %s: %v", m.baseDir, err)
			}
		}

		for _, entry := range entries {
			if entry.IsDir() {
				m.translations = append(m.translations, entry)
			}
		}
	}

	return m
}

func displayTranslations(m modal) string {
	translationListEnter(m)

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if len(m.translations) == 0 {
		return "No translations found. Run 'scan' first.\n"
	}

	tpl := "Translations:\n\n%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	checkboxes := []string{}
	c := m.cursor
	for i, t := range m.translations {
		checkboxes = append(checkboxes, checkbox(t.Name(), c == i))
	}

	choices := strings.Join(checkboxes, "\n")

	return fmt.Sprintf(tpl, choices)
}

func translationListUpdate(m modal, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = len(m.translations) - 1
		}
	case "down", "j":
		if m.cursor < len(m.translations)-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}
	case "enter":
		return bookListEnter(m), nil
	}
	return m, nil
}
