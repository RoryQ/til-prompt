package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	gap "github.com/muesli/go-app-paths"

	"github.com/roryq/til-prompt/internal/core"
)

var (
	scope = gap.NewScope(gap.User, "til")
)

var CLI struct {
	New    struct{} `cmd default:"1" help:"Create a new TIL entry. (default command)"`
	Config struct {
		List struct{} `cmd default:"1" help:"List the current configuration. (default sub-command)"`
		Edit struct{} `cmd help:"Open the config in your configured $EDITOR"`
	} `cmd help:"Manage the current configuration."`
}

func main() {
	config, err := core.LoadConfig(scope)
	if err != nil {
		fmt.Printf("could not load config: %s\n", err)
		os.Exit(1)
	}

	ktx := kong.Parse(&CLI,
		kong.Name("til"),
		kong.Description("An interactive prompt for managing TIL entries."))

	switch ktx.Command() {
	case "config list":
		fmt.Println(config.Formatted())
	case "config edit":
		err := core.EditConfig(scope)
		if err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	case "edit":
		break
	default:
		if err := tea.NewProgram(core.NewUI(config)).Start(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
