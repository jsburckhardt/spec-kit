// Package ui provides terminal user interface components
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
)

// ProgressRenderer handles rendering of progress tracking
type ProgressRenderer struct {
	tracker *config.StepTracker
}

// NewProgressRenderer creates a new progress renderer
func NewProgressRenderer(tracker *config.StepTracker) *ProgressRenderer {
	return &ProgressRenderer{
		tracker: tracker,
	}
}

// Render generates the progress tree display
func (pr *ProgressRenderer) Render() string {
	var output strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("cyan")).
		Bold(true).
		Render(pr.tracker.GetTitle())
	output.WriteString(title + "\n\n")

	// Steps
	steps := pr.tracker.GetSteps()
	for _, step := range steps {
		line := pr.renderStep(step)
		output.WriteString(line + "\n")
	}

	return output.String()
}

// renderStep renders a single step
func (pr *ProgressRenderer) renderStep(step config.Step) string {
	var symbol, style string

	switch step.Status {
	case config.StatusDone:
		symbol = "●"
		style = "green"
	case config.StatusRunning:
		symbol = "○"
		style = "cyan"
	case config.StatusError:
		symbol = "●"
		style = "red"
	case config.StatusSkipped:
		symbol = "○"
		style = "yellow"
	default: // StatusPending
		symbol = "○"
		style = "240" // dim gray
	}

	// Style the symbol
	styledSymbol := lipgloss.NewStyle().
		Foreground(lipgloss.Color(style)).
		Render(symbol)

	// Build the label
	label := step.Label
	if step.Detail != "" {
		if step.Status == config.StatusPending {
			label = fmt.Sprintf("%s (%s)", label, step.Detail)
		} else {
			label = fmt.Sprintf("%s (%s)",
				lipgloss.NewStyle().Render(label),
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render(step.Detail))
		}
	}

	// Apply pending style to entire line if needed
	if step.Status == config.StatusPending {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("%s %s", styledSymbol, label))
	}

	return fmt.Sprintf("%s %s", styledSymbol, label)
}

// LiveProgress provides live progress updates
type LiveProgress struct {
	renderer *ProgressRenderer
}

// NewLiveProgress creates a new live progress display
func NewLiveProgress(tracker *config.StepTracker) *LiveProgress {
	return &LiveProgress{
		renderer: NewProgressRenderer(tracker),
	}
}

// Render returns the current progress display
func (lp *LiveProgress) Render() string {
	return lp.renderer.Render()
}
