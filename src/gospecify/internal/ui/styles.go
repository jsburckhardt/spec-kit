// Package ui provides terminal user interface components
package ui

import "github.com/charmbracelet/lipgloss"

// Common styles for consistent UI
var (
	// Colors
	ColorCyan   = lipgloss.Color("cyan")
	ColorGreen  = lipgloss.Color("green")
	ColorRed    = lipgloss.Color("red")
	ColorYellow = lipgloss.Color("yellow")
	ColorGray   = lipgloss.Color("240")
	ColorWhite  = lipgloss.Color("white")
	ColorBlack  = lipgloss.Color("black")

	// Text styles
	BoldStyle = lipgloss.NewStyle().Bold(true)

	CyanStyle = lipgloss.NewStyle().Foreground(ColorCyan)

	GreenStyle = lipgloss.NewStyle().Foreground(ColorGreen)

	RedStyle = lipgloss.NewStyle().Foreground(ColorRed)

	YellowStyle = lipgloss.NewStyle().Foreground(ColorYellow)

	GrayStyle = lipgloss.NewStyle().Foreground(ColorGray)

	WhiteStyle = lipgloss.NewStyle().Foreground(ColorWhite)

	// Panel styles
	InfoPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCyan).
			Padding(1, 2)

	WarningPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorYellow).
			Padding(1, 2)

	ErrorPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorRed).
			Padding(1, 2)

	SuccessPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGreen).
			Padding(1, 2)

	// Banner style
	BannerStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true).
			Align(lipgloss.Center)

	// Tagline style
	TaglineStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Italic(true).
			Align(lipgloss.Center)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorGray)

	// Selection styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Bold(true)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorWhite)

	// Progress symbols
	ProgressDone    = GreenStyle.Render("●")
	ProgressRunning = CyanStyle.Render("○")
	ProgressPending = GrayStyle.Render("○")
	ProgressError   = RedStyle.Render("●")
	ProgressSkipped = YellowStyle.Render("○")
)
