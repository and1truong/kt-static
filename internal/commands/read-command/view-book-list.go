package read_command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func bookListEnter(m modal) modal {
	m.translationPath = filepath.Join(m.baseDir, m.translations[m.cursor].Name())
	m.state = viewBooks
	m.cursor = 0

	books, err := os.ReadDir(m.translationPath)
	if err != nil {
		m.err = fmt.Errorf("failed to read translation directory %s: %v", m.translationPath, err)
	}

	for _, bookEntry := range books {
		if bookEntry.IsDir() {
			categoryPath := filepath.Join(m.translationPath, bookEntry.Name(), "_category_.json")

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
			m.books = append(m.books, bookInfo{
				Label:    category.Label,
				DirEntry: bookEntry,
				Position: category.Position,
			})
		}
	}

	sort.Slice(m.books, func(i, j int) bool {
		return m.books[i].Position < m.books[j].Position
	})

	return m
}

func bookListDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if len(m.books) == 0 {
		return "No books found in " + filepath.Base(m.translationPath) + "\n"
	}

	// 1. Generate all checkbox strings
	var allCheckboxes []string
	c := m.cursor
	for i, b := range m.books {
		allCheckboxes = append(allCheckboxes, checkbox(b.Label, c == i))
	}

	var choices string
	if len(m.books) > 14 {
		// Multi-column logic: max 10 rows per column
		const maxRows = 14
		numBooks := len(allCheckboxes)
		numCols := (numBooks + maxRows - 1) / maxRows // Ceiling division

		// 1. Split into columns
		columns := make([][]string, numCols)
		for i := 0; i < numCols; i++ {
			start := i * maxRows
			end := (i + 1) * maxRows
			if end > numBooks {
				end = numBooks
			}
			columns[i] = allCheckboxes[start:end]
		}

		// 2. Find max width for each column
		maxWidths := make([]int, numCols)
		for i := 0; i < numCols; i++ {
			for _, s := range columns[i] {
				w := lipgloss.Width(s)
				if w > maxWidths[i] {
					maxWidths[i] = w
				}
			}
		}

		// 3. Combine columns row by row
		var lines []string
		separator := "  " // 2 spaces separator
		for r := 0; r < maxRows; r++ {
			var rowStrings []string
			for c := 0; c < numCols; c++ {
				col := columns[c]
				if r < len(col) {
					s := col[r]
					// Pad the string to its column's max width
					s = lipgloss.NewStyle().Width(maxWidths[c]).Render(s)
					rowStrings = append(rowStrings, s)
				} else {
					// Add empty padding for shorter columns
					rowStrings = append(rowStrings, lipgloss.NewStyle().Width(maxWidths[c]).Render(""))
				}
			}
			// Join the row strings with the separator
			lines = append(lines, strings.Join(rowStrings, separator))
		}
		choices = strings.Join(lines, "\n")

	} else {
		// Single-column logic (existing)
		choices = strings.Join(allCheckboxes, "\n")
	}

	tpl := "Books in " + filepath.Base(m.translationPath) + ":\n\n%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	return fmt.Sprintf(tpl, choices)
}

func bookListUpdate(m modal, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	const maxRows = 14
	numBooks := len(m.books)
	numCols := (numBooks + maxRows - 1) / maxRows

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = numBooks - 1
		}
	case "down", "j":
		if m.cursor < numBooks-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}
	case "left", "h":
		if numCols > 1 {
			currentCol := m.cursor / maxRows
			if currentCol > 0 {
				m.cursor -= maxRows
			} else {
				// Wrap around to the last column, trying to keep the same row
				row := m.cursor % maxRows
				m.cursor = (numCols-1)*maxRows + row
				if m.cursor >= numBooks {
					m.cursor = numBooks - 1
				}
			}
		}
	case "right", "l":
		if numCols > 1 {
			currentCol := m.cursor / maxRows
			if currentCol < numCols-1 {
				m.cursor += maxRows
				if m.cursor >= numBooks {
					m.cursor = numBooks - 1
				}
			} else {
				// Wrap around to the first column
				row := m.cursor % maxRows
				m.cursor = row
			}
		}
	case "enter":
		return chapterListEnter(m), nil
	case "esc":
		m.state = viewTranslations
		m.books = []bookInfo{}
		m.bookCursor = 0
		m.cursor = 0
		return m, nil
	}

	return m, nil
}
