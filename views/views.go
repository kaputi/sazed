package views

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#FF00FF")).
				Foreground(lipgloss.Color("#000000")).
				Inline(true)

	titleStyle = lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("#FF00FF"))
)

type focusCmd bool
