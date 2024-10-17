package bubbletasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imyounas/tuitaskapp/internal/database"
)

type appState int

const (
	tableView appState = iota
	filterView
	detailView
	editlDetailView
	addDetailView
	deleteView
)

type taskMainModel struct {
	filter     filterModel
	table      tableModel
	detailView detailModel
	focusIndex int
	state      appState
	height     int
	width      int
}

func InitialMainModel(tasks []database.Task) taskMainModel {

	return taskMainModel{
		filter:     newFilterModel(),
		table:      newTableModel(tasks),
		detailView: newDetailViewModel(),
		focusIndex: 1, // Start with focus on the table
		state:      tableView,
	}
}

func (m taskMainModel) Init() tea.Cmd {
	return nil
}
func (m *taskMainModel) resetFilter() {

	m.table.filteredTasks = m.table.tasks
	m.filter.input.SetValue("")
	m.table.updateTableRows()
}

func (m *taskMainModel) resetViewFocus() {
	m.table.SetBlur()
	m.filter.SetBlur()
	m.detailView.SetBlur()
	m.detailView.doneChanges = false
	m.detailView.focusIndex = 0
	switch m.state {
	case filterView:
		m.filter.SetFocus()
	case tableView:
		m.table.SetFocus()
	case editlDetailView, addDetailView:
		m.detailView.SetFocus()
	}
}

func (m taskMainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.table.SetSize(msg.Width-20, msg.Height) // Reserve space for detail view

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "f":
			if m.state != filterView {
				m.state = filterView
				m.resetViewFocus()
			}

		case "r":
			if m.state != filterView {
				m.resetFilter()
				m.state = tableView
				m.resetViewFocus()
				return m, nil
			}
		case "n":
			if m.state == tableView {
				m.state = addDetailView
				m.detailView.editing = true
				m.detailView.setTask(&database.Task{})
				m.resetFilter()
				m.resetViewFocus()
				return m, nil
			}
		case "enter":
			if m.state == tableView {
				selectedTask := m.table.SelectedTask()
				m.detailView.setTask(selectedTask)
				m.state = editlDetailView
				m.detailView.editing = true
				m.resetViewFocus()
				return m, nil

			} else if m.state == filterView {
				filterText := m.filter.input.Value()
				m.state = tableView
				m.table.filterTasks(filterText)
				m.resetViewFocus()
				return m, nil
			}

		case "esc":
			{
				m.table.filterTasks("")
				m.state = tableView
				m.resetViewFocus()
				return m, nil
			}
		}
	}

	switch m.state {

	case tableView:
		{
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
			if len(m.table.tasks) > 0 {
				selectedTask := m.table.SelectedTask()
				m.detailView.setTask(selectedTask)
				m.detailView.SetBlur()
			}
		}
	case editlDetailView, addDetailView:
		{

			m.detailView, cmd = m.detailView.Update(msg)
			cmds = append(cmds, cmd)

			if m.detailView.doneChanges {

				if m.state == editlDetailView {
					selectedTask := m.table.SelectedTask()

					selectedTask.Name = m.detailView.inputs[0].Value()
					selectedTask.Description = m.detailView.inputs[1].Value()
					selectedTask.Notes = m.detailView.inputs[2].Value()
					selectedTask.Assignee = m.detailView.inputs[3].Value()
				} else if m.state == addDetailView {
					newTask := database.Task{
						Id:          len(m.table.tasks) + 1,
						Name:        m.detailView.inputs[0].Value(),
						Description: m.detailView.inputs[1].Value(),
						Notes:       m.detailView.inputs[2].Value(),
						Assignee:    m.detailView.inputs[3].Value(),
					}
					m.table.tasks = append(m.table.tasks, newTask)
				}

				m.table.UpdateTasks(m.table.tasks)
				m.state = tableView
				m.resetViewFocus()
			}
		}
	case filterView:
		{
			m.filter, cmd = m.filter.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m taskMainModel) View() string {
	filterView := m.filter.View()
	tableView := m.table.View()
	detailView := m.detailView.View()

	tableAndDetailContent := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tableView,
		detailView,
	)

	title := "Bubble TUI - Task ToDo List"
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title),
		filterView,
		tableAndDetailContent,
		"n: New task • enter: Edit task • d: Delete task • f: Filter tasks • ←→: Navigate pages • q: Quit",
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)

}
