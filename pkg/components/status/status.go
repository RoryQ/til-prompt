package status

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	// Status Bar.
	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusTextStyle = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

func Render(leftBlock, description string) string {
	width, _, _ := term.GetSize(0)
	w := lipgloss.Width

	statusKey := statusStyle.Render(strings.ToUpper(leftBlock))
	appName := fishCakeStyle.Render("TIL-Prompt")

	statusVal := statusTextStyle.Copy().
		Width(width - w(statusKey) - w(appName)).
		Render(description)

	bar := lipgloss.JoinHorizontal(lipgloss.Bottom,
		statusKey,
		statusVal,
		appName,
	)

	return statusBarStyle.Width(width).Render(bar)
}
