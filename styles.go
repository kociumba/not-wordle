package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// wrong          = lipgloss.NewStyle().Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#FF567B"))
	wrong          = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#000000"))
	rightPlace     = lipgloss.NewStyle().Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#BDFF52"))
	containsLetter = lipgloss.NewStyle().Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#FFE96B"))
	normalText     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
)
