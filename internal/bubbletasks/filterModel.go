package bubbletasks

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type filterModel struct {
	input   textinput.Model
	focused bool
}

func newFilterModel() filterModel {
	input := textinput.New()
	input.Placeholder = "Filter tasks..."
	input.Width = 50

	//input.Focus()
	return filterModel{input: input, focused: false}
}

func (m filterModel) Init() tea.Cmd {

	return textinput.Blink
}

func (m *filterModel) SetFocus() {
	m.input.Focus()
	m.focused = true
}

func (m *filterModel) SetBlur() {
	m.input.Blur()
	m.focused = false
}

func (m filterModel) Update(msg tea.Msg) (filterModel, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m filterModel) View() string {

	filter := focusedStyle.Render(m.input.View())
	if !m.focused {
		filter = blurredStyle.Render(m.input.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		filter,
	)

}
