package read_command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func chapterListEnter(m modal) modal {
	m.bookPath = filepath.Join(m.translationPath, m.books[m.cursor].DirEntry.Name())
	m.state = viewChapters
	m.bookCursor = m.cursor
	m.cursor = 0

	if m.chapters == nil {
		chapters, err := os.ReadDir(m.bookPath)
		if err != nil {
			m.err = fmt.Errorf("failed to read book directory %s: %v", m.bookPath, err)
		}

		for _, chapterEntry := range chapters {
			if chapterEntry.Name() != "_category_.json" {
				m.chapters = append(m.chapters, chapterEntry)
			}
		}

		sort.Slice(m.chapters, func(i, j int) bool {
			nameI := m.chapters[i].Name()
			nameJ := m.chapters[j].Name()

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

		m.chapterLabels = make([]string, len(m.chapters))
		for i, chapterEntry := range m.chapters {
			name := chapterEntry.Name()
			numStr := strings.TrimSuffix(name, filepath.Ext(name))
			m.chapterLabels[i] = fmt.Sprintf("Chapter %s", numStr)
		}
		m.chapterTotal = len(m.chapters)
	}

	return m
}

func chapterListDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if len(m.chapters) == 0 {
		return "No chapters found in " + m.books[m.cursor].Label + "\n"
	}

	// 1. Generate all checkbox strings
	var allCheckboxes []string
	c := m.cursor
	for i, label := range m.chapterLabels {
		allCheckboxes = append(allCheckboxes, checkbox(label, c == i))
	}

	var choices string
	if len(m.chapters) > 14 {
		// Multi-column logic: max 14 rows per column
		const maxRows = 14
		numChapters := len(allCheckboxes)
		numCols := (numChapters + maxRows - 1) / maxRows // Ceiling division

		// 1. Split into columns
		columns := make([][]string, numCols)
		for i := 0; i < numCols; i++ {
			start := i * maxRows
			end := (i + 1) * maxRows
			if end > numChapters {
				end = numChapters
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

	tpl := "Chapters in " + m.books[m.bookCursor].Label + ":\n\n%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	return fmt.Sprintf(tpl, choices)
}

func chapterListUpdate(m modal, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	const maxRows = 14
	numChapters := len(m.chapters)
	numCols := (numChapters + maxRows - 1) / maxRows

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = numChapters - 1
		}
	case "down", "j":
		if m.cursor < numChapters-1 {
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
				if m.cursor >= numChapters {
					m.cursor = numChapters - 1
				}
			}
		}
	case "right", "l":
		if numCols > 1 {
			currentCol := m.cursor / maxRows
			if currentCol < numCols-1 {
				m.cursor += maxRows
				if m.cursor >= numChapters {
					m.cursor = numChapters - 1
				}
			} else {
				// Wrap around to the first column
				row := m.cursor % maxRows
				m.cursor = row
			}
		}
	case "enter":
		return chapterContentEnter(m), nil
	}
	return m, nil
}
