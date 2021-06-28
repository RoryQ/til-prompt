package manager

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/roryq/til-prompt/internal/core"
	"github.com/roryq/til-prompt/pkg/components/status"
)

type Model struct {
	config core.Config

	keyword        string
	entries        []string
	selectionIndex int
	selection      chan Selection
	mode           Mode

	exited bool
}

//go:generate stringer -type=Mode
type Mode int

const (
	Edit Mode = iota
	View
)

type Selection struct {
	Path string
	Mode Mode
}

var (
	today = time.Now().Format("2006-01-02")
)

func NewUI(config core.Config, mode Mode, selection chan Selection, keyword string) Model {
	searchTerm := today
	if keyword != "" {
		searchTerm = keyword
	}
	entries, err := findEntries(config.SaveDirectory, searchTerm)
	if err != nil {
		return Model{}
	}

	return Model{
		config:    config,
		keyword:   keyword,
		entries:   entries,
		mode:      mode,
		selection: selection,
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
		fmt.Println("No entries found")
		return tea.Quit
	}
	if len(m.entries) == 1 {
		return m.launch()
	}
	return textinput.Blink
}

// launch editor or viewer for selected file
func (m Model) launch() tea.Cmd {
	return func() tea.Msg {
		selectedPath := path.Join(m.config.SaveDirectory, m.entries[m.selectionIndex])
		m.selection <- Selection{
			Path: selectedPath,
			Mode: m.mode,
		}
		return tea.Quit()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "e":
			m.mode = Edit
		case "v":
			m.mode = View

		case "enter":
			return m, m.launch()

		case "down", "j":
			m.selectionIndex++
			m.selectionIndex %= len(m.entries)
			if m.selectionIndex >= len(m.entries) {
				m.selectionIndex = 0
			}

		case "up", "k":
			m.selectionIndex--
			if m.selectionIndex < 0 {
				m.selectionIndex = len(m.entries) - 1
			}
		}
	}

	return m, tea.Batch()
}

func (m Model) View() string {
	b := new(strings.Builder)

	if len(m.entries) > 1 {
		m.renderEntries(b)
		m.renderStatusBar(b)
	}

	return b.String()
}

func (m Model) renderEntries(b *strings.Builder) {
	for i := range m.entries {
		if m.selectionIndex == i {
			b.WriteString("[✔] ")
		} else {
			b.WriteString("[ ] ")
		}
		b.WriteString(m.entries[i])
		b.WriteRune('\n')
	}
}

func (m Model) renderStatusBar(b *strings.Builder) {
	b.WriteString(status.Render(m.mode.String(), m.mode.Help()))
	b.WriteRune('\n')
}

func (i Mode) Help() string {
	switch i {
	case Edit:
		return "Edit selected file. Press v to switch to View mode."
	case View:
		return "View selected file. Press e to switch to Edit mode."
	}
	return ""
}
