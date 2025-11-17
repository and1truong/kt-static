package read_command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

func chapterContentEnter(m modal) modal {
	m.chapterPath = filepath.Join(m.bookPath, m.chapters[m.cursor].Name())
	m.state = viewChapterContent
	m.chapterCursor = m.cursor
	m.cursor = 0

	if m.chapterContent == "" {
		contentBytes, err := os.ReadFile(m.chapterPath)
		if err != nil {
			m.err = fmt.Errorf("failed to read chapter file %s: %v", m.chapterPath, err)
		}

		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
		)

		m.chapterContent, _ = r.Render(string(contentBytes))

		m.viewport.SetContent(m.chapterContent)
	}

	return m
}

func chapterContentUpdate(m modal, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = viewChapters
			m.chapterContent = ""
			m.cursor = m.chapterCursor
			return m, nil
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(headerView(m))
		footerHeight := lipgloss.Height(footerView(m))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.chapterContent)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func chapterContentDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	header := headerView(m)
	footer := footerView(m)
	content := m.viewport.View()

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

func headerView(m modal) string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func footerView(m modal) string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
