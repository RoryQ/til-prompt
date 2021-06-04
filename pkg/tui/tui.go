package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kennygrant/sanitize"
)

type Model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode

	config Config
}

var (
	today = time.Now().Format("2006-01-02")

	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedSave   = focusedButton("Save")
	unfocusedSave = blurredButton("Save")
)

func focusedButton(name string) string {
	return focusedStyle.Copy().Render("[ " + name + " ]")
}

func blurredButton(name string) string {
	return fmt.Sprintf("[ %s ]", blurredStyle.Render(name))
}

func NewUI(config Config) Model {
	m := Model{
		inputs: make([]textinput.Model, 3),
		config: config,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "What I Learned Today"
		case 2:
			t.Placeholder = "Category"
		}

		m.inputs[i] = t
	}

	return m
}
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.saveFocused() {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) saveFocused() bool {
	return m.focusIndex == len(m.inputs)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	b := &strings.Builder{}
	m.renderInputs(b)
	m.renderButtons(b)
	m.renderFooter(b)
	return b.String()
}

func formatTitle(s string) string {
	if s == "" {
		return ""
	}

	return sanitize.Name("-" + s)
}

func formatDirectory(s string) string {
	if s == "" {
		return ""
	}

	return sanitize.Name(s) + string(os.PathSeparator)
}

func (m Model) renderFooter(b *strings.Builder) {
	b.WriteString(helpStyle.Render("your til will be saved to "))
	filename := fmt.Sprintf("%s%s%s.md", formatDirectory(m.inputs[2].Value()), today, formatTitle(m.inputs[0].Value()))
	b.WriteString(cursorModeHelpStyle.Render(filename))
}

func (m Model) renderButtons(b *strings.Builder) {
	button := &unfocusedSave
	if m.focusIndex == len(m.inputs) {
		button = &focusedSave
	}
	fmt.Fprintf(b, "\n\n%s\n\n", *button)
}

func (m Model) renderInputs(b *strings.Builder) {
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
}
