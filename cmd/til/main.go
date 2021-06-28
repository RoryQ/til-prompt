package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"
	"strings"

	"github.com/roryq/til-prompt/pkg/editor"
	"github.com/roryq/til-prompt/pkg/viewer"

	"github.com/roryq/til-prompt/internal/manager"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	gap "github.com/muesli/go-app-paths"

	"github.com/roryq/til-prompt/internal/core"
)

var (
	scope = gap.NewScope(gap.User, "til")
)

var cli struct {
	New    struct{} `cmd default:"1" help:"Create a new TIL entry. (default command)"`
	Config struct {
		List struct{} `cmd default:"1" help:"List the current configuration. (default sub-command)"`
		Edit struct{} `cmd help:"Open the config in your configured $EDITOR"`
	} `cmd help:"Manage the current configuration."`
	Edit struct {
		Keyword string `arg optional help:"optional keyword to search for."`
	} `cmd help:"Open today's TIL in your configured editor. You can optionally pass a keyword to search for."`
	View struct {
		Keyword string `arg optional help:"optional keyword to search for."`
	} `cmd help:"Display a styled view of today's TIL. You can optionally pass a keyword to search for."`
}

func main() {
	config, err := core.LoadConfig(scope)
	if err != nil {
		fmt.Printf("could not load config: %s\n", err)
		os.Exit(1)
	}

	ktx := kong.Parse(&cli,
		kong.Name("til"),
		kong.Description("An interactive prompt for managing TIL entries."))

	switch c := ktx.Command(); c {
	case "config list":
		fmt.Println(config.Formatted())
	case "config edit":
		err := core.EditConfig(scope)
		if err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	case "edit", "edit <keyword>", "view", "view <keyword>":
		selection := make(chan manager.Selection, 1)
		program := new(tea.Program)
		if strings.HasPrefix(c, "edit") {
			program = tea.NewProgram(manager.NewUI(config, manager.Edit, selection, cli.Edit.Keyword))
		} else {
			program = tea.NewProgram(manager.NewUI(config, manager.View, selection, cli.View.Keyword))
		}

		if err := program.Start(); err != nil {
			fmt.Printf("could not start editor: %s\n", err)
			os.Exit(1)
		}

		select {
		case c := <-selection:
			if c.Mode == manager.Edit {
				editor.Launch(config, c.Path)
			} else {
				viewer.Launch(c.Path)
			}
		default:
			os.Exit(0)
		}
	default:
		if err := tea.NewProgram(core.NewUI(config)).Start(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
