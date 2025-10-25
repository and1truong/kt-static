package read_command

import (
	"fmt"
	"htruong/kt-crawler/internal"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type viewState int

const (
	viewTranslations viewState = iota
	viewBooks
	viewChapters
	viewChapterContent
	dotChar = " • "
)

var (
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

func RunRead(c *cli.Context, config *internal.Config) error {
	m := modal{
		config:  config,
		baseDir: config.Listeners.Store.Filesystem.Directory,
		state:   viewTranslations,
	}

	// entering to first screen
	m = translationListEnter(m)

	// start program until error
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()

	return err
}

type modal struct {
	config *internal.Config

	state    viewState
	cursor   int
	baseDir  string
	err      error
	quitting bool

	// translations
	translations    []os.DirEntry
	translationPath string

	// books
	books      []bookInfo
	bookCursor int
	bookPath   string

	// chapters
	chapters      []os.DirEntry
	chapterCursor int
	chapterPath   string
	chapterLabels []string
	chapterTotal  int

	// content
	chapterContent string
}

func (m modal) Init() tea.Cmd {
	return nil
}

func (m modal) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var s string
	switch m.state {
	case viewTranslations:
		s = displayTranslations(m)
	case viewBooks:
		s = bookListDisplay(m)
	case viewChapters:
		s = chapterListDisplay(m)
	case viewChapterContent:
		s = chapterContentDisplay(m)
	default:
		panic("unhandled default case")
	}

	return mainStyle.Render("\n" + s + "\n\n")
}

func (m modal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case viewTranslations:
			return translationListUpdate(m, msg)
		case viewBooks:
			return bookListUpdate(m, msg)
		case viewChapters:
			return chapterListUpdate(m, msg)
		case viewChapterContent:
			return chapterContentUpdate(m, msg)
		default:
			// panic("unhandled default case")
		}
	}

	return m, nil
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
