package views

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	purple1 = "#e5d4ed"
	purple2 = "#6d72c3"
	purple3 = "#5941a9"
	purple4 = "#514f59"
	purple5 = "#1d1128"

	oxfordBlue    = "#191d32"
	spaceCadet    = "#282f44"
	englishViolet = "#453a49"
	wine          = "#6d3b47"
	magentaDye    = "#ba2c73"
)

func InterpolateColor(a, b [3]int8, t float64) [3]int8 {
	var out [3]int8
	for i := 0; i < 3; i++ {
		a[i] = int8(float64(a[i]) + float64(b[i]-a[i])*t)
	}
	return out
}

var (
	focusedTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(magentaDye)).
				Foreground(lipgloss.Color(oxfordBlue)).
				Inline(true)

	titleStyle = lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color(purple2))

	otherStyle = lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color(purple3))
)

type focusCmd bool
