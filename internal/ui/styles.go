package ui

import "github.com/charmbracelet/lipgloss"

var (
	activePane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2"))

	inactivePane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	statusBar = lipgloss.NewStyle().
			Reverse(true).
			Bold(true)

	filterBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true)

	queryBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true)

	helpKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3"))

	helpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6"))

	selectedItem = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10"))
)
