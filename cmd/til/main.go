package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	gap "github.com/muesli/go-app-paths"

	"github.com/roryq/til-prompt/pkg/tui"
)

var (
	scope = gap.NewScope(gap.User, "til")
)

var CLI struct {
	Config struct{} `cmd:"config" help:"Displays the current configuration."`
}

func main() {
	config, err := tui.LoadConfig(scope)
	if err != nil {
		fmt.Printf("could not load config: %s\n", err)
		os.Exit(1)
	}

	ktx := kong.Parse(&CLI, kong.Name("til"))
	switch ktx.Command() {
	case "config":
		fmt.Println(config.Sprint())
	default:
		if err := tea.NewProgram(tui.NewUI(config)).Start(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
