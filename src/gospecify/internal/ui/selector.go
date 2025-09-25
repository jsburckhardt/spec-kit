// Package ui provides terminal user interface components
package ui

import (
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Selector provides an interactive selection interface
type Selector struct {
	list     list.Model
	selected string
	quitting bool
}

// NewSelector creates a new interactive selector
func NewSelector(prompt string, options map[string]string, defaultKey string) *Selector {
	var items []list.Item
	var keys []string

	// Sort keys for consistent ordering
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create list items
	for _, key := range keys {
		items = append(items, selectorItem{
			key:   key,
			value: options[key],
		})
	}

	// Create the list
	l := list.New(items, selectorDelegate{}, 0, 0)
	l.Title = prompt
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	// Set default selection
	if defaultKey != "" {
		for i, key := range keys {
			if key == defaultKey {
				l.Select(i)
				break
			}
		}
	}

	return &Selector{
		list: l,
	}
}

// Run starts the interactive selection
func (s *Selector) Run() (string, error) {
	p := tea.NewProgram(s)
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	finalModel := result.(*Selector)
	return finalModel.selected, nil
}

// Init initializes the Bubbletea model
func (s *Selector) Init() tea.Cmd {
	return nil
}

// Update handles user input
func (s *Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item := s.list.SelectedItem(); item != nil {
				s.selected = item.(selectorItem).key
				s.quitting = true
				return s, tea.Quit
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "ctrl+c"))):
			s.quitting = true
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		s.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// View renders the selector
func (s *Selector) View() string {
	if s.quitting {
		return ""
	}

	return docStyle.Render(s.list.View())
}

// selectorItem represents an item in the selector
type selectorItem struct {
	key   string
	value string
}

// FilterValue returns the value for filtering (not used)
func (i selectorItem) FilterValue() string {
	return i.key
}

// selectorDelegate handles item rendering
type selectorDelegate struct{}

func (d selectorDelegate) Height() int {
	return 1
}

func (d selectorDelegate) Spacing() int {
	return 0
}

func (d selectorDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d selectorDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(selectorItem)

	var style lipgloss.Style
	if index == m.Index() {
		style = selectedItemStyle
	} else {
		style = itemStyle
	}

	fmt.Fprintf(w, style.Render(fmt.Sprintf("%s (%s)", i.key, i.value)))
}

// Styles
var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("white"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("cyan")).
				Bold(true)
)
