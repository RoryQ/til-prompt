package edit

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/roryq/til-prompt/pkg/editor"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/roryq/til-prompt/internal/core"
)

type Model struct {
	config core.Config

	keyword   string
	entries   []string
	selection int
}

var (
	today = time.Now().Format("2006-01-02")
)

func NewUI(config core.Config, keyword string) Model {
	searchTerm := today
	if keyword != "" {
		searchTerm = keyword
	}
	entries, err := findEntries(config.SaveDirectory, searchTerm)
	if err != nil {
		return Model{}
	}

	return Model{
		config:  config,
		keyword: keyword,
		entries: entries,
	}
}

func findEntries(saveDirectory, term string) ([]string, error) {
	var entries []string
	if err := fs.WalkDir(os.DirFS(saveDirectory), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.Contains(path, term) {
			entries = append(entries, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return entries, nil
}

func (m Model) Init() tea.Cmd {
	if len(m.entries) == 0 {
		return tea.Quit
	}
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			selectedPath := path.Join(m.config.SaveDirectory, m.entries[m.selection])
			editor.Launch(m.config, selectedPath)
			return m, tea.Quit

		case "down", "j":
			m.selection++
			m.selection %= len(m.entries)
			if m.selection >= len(m.entries) {
				m.selection = 0
			}

		case "up", "k":
			m.selection--
			if m.selection < 0 {
				m.selection = len(m.entries) - 1
			}
		}
	}

	return m, tea.Batch()
}

func (m Model) View() string {
	b := new(strings.Builder)
	m.renderEntries(b)
	return b.String()
}

func (m Model) renderEntries(b *strings.Builder) {
	for i := range m.entries {
		if m.selection == i {
			b.WriteString("[✔] ")
		} else {
			b.WriteString("[ ] ")
		}
		b.WriteString(m.entries[i])
		b.WriteRune('\n')
	}
}
