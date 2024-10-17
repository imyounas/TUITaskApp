package bubbletasks

import "github.com/charmbracelet/lipgloss"

var (
	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	blurredStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	focusedButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render("[ Submit ]")

	blurredButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("[ Submit ]")
)
