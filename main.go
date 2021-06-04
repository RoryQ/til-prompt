package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"

	gap "github.com/muesli/go-app-paths"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/roryq/til-prompt/pkg/tui"
)

var (
	scope = gap.NewScope(gap.User, "til")
)

func main() {
	config, err := tui.LoadConfig(scope)
	if err != nil {
		fmt.Printf("could not load config: %s\n", err)
		os.Exit(1)
	}

	if err := tea.NewProgram(tui.NewUI(config)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
