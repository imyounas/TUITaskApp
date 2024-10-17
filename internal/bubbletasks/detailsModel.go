package bubbletasks

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imyounas/tuitaskapp/internal/database"
)

type detailModel struct {
	inputs      []textinput.Model
	taskId      textinput.Model
	currentTask *database.Task
	editing     bool
	focusIndex  int
	focused     bool
	doneChanges bool
}

func newDetailViewModel() detailModel {
	m := detailModel{
		inputs:  make([]textinput.Model, 4),
		focused: false,
	}

	m.taskId.Placeholder = "Id"
	m.taskId.Width = 25

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 25
		t.Width = 25

		switch i {

		case 0:
			t.Placeholder = "Name"
		case 1:
			t.Placeholder = "Description"
		case 2:
			t.Placeholder = "Notes"
		case 3:
			t.Placeholder = "Assignee"
		}

		m.inputs[i] = t
	}

	return m
}

func (m detailModel) Init() tea.Cmd {
	return nil
}

func (m *detailModel) SetFocus() {
	m.taskId.Blur()
	m.inputs[0].Focus()
	m.focused = true
}

func (m *detailModel) SetBlur() {
	for i := 0; i < len(m.inputs); i++ {
		m.inputs[i].Blur()

	}
	m.taskId.Blur()
	m.focusIndex = 0
	m.focused = false
}

func (m *detailModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *detailModel) setTask(task *database.Task) {
	m.taskId.SetValue(strconv.Itoa(task.Id))
	m.inputs[0].SetValue(task.Name)
	m.inputs[1].SetValue(task.Description)
	m.inputs[2].SetValue(task.Notes)
	m.inputs[3].SetValue(task.Assignee)
	m.currentTask = task
}

func (m detailModel) Update(msg tea.Msg) (detailModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down", "esc":

			if m.focused {
				s := msg.String()

				if s == "enter" && m.focusIndex == len(m.inputs) {
					m.doneChanges = true
					return m, nil
				}

				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

				if s == "esc" {
					m.doneChanges = true
				}

				cmds = make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
						continue
					}
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m detailModel) View() string {
	var b string

	b += fmt.Sprintf(
		"%s:%s\n\n",
		m.taskId.Placeholder,
		m.taskId.View(),
	)

	for _, input := range m.inputs {
		b += fmt.Sprintf(
			"%s:%s\n\n",
			input.Placeholder,
			input.View(),
		)
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	b += fmt.Sprintf("%s\n\n", *button)

	details := focusedStyle.Render(b)
	if !m.focused {
		details = blurredStyle.Render(b)
	}

	fixedWidthDetailView := lipgloss.NewStyle().Width(45).PaddingBottom(1).Render(details)
	//return details
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		fixedWidthDetailView,
	)
}
