package bubbletasks

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imyounas/tuitaskapp/internal/database"
)

type tableModel struct {
	table         table.Model
	tasks         []database.Task
	filteredTasks []database.Task
	focused       bool
	width         int
	height        int
	currentPage   int
	itemsPerPage  int
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)
	paginationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

func newTableModel(tasks []database.Task) tableModel {

	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Description", Width: 30},
		//{Title: "Assignee", Width: 10},
	}

	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = table.Row{
			fmt.Sprintf("%d", task.Id),
			task.Name,
			task.Description,
			//task.Assignee,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{
		tasks:         tasks,
		filteredTasks: tasks,
		table:         t,
		itemsPerPage:  5,
		currentPage:   1,
		focused:       true,
	}
	m.updateTableRows()
	return m
}

func (m *tableModel) updateTableRows() {
	if len(m.filteredTasks) == 0 {
		m.table.SetRows([]table.Row{})
		return
	}

	start := (m.currentPage - 1) * m.itemsPerPage
	end := start + m.itemsPerPage

	if start >= len(m.filteredTasks) {
		m.currentPage = (len(m.filteredTasks)-1)/m.itemsPerPage + 1
		start = (m.currentPage - 1) * m.itemsPerPage
		end = start + m.itemsPerPage
	}

	if end > len(m.filteredTasks) {
		end = len(m.filteredTasks)
	}

	rows := make([]table.Row, 0, end-start)
	for _, task := range m.filteredTasks[start:end] {
		rows = append(rows, table.Row{fmt.Sprintf("%d", task.Id), task.Name, task.Description})
	}
	m.table.SetRows(rows)

	// Ensure cursor is within bounds
	if m.table.Cursor() >= len(rows) {
		m.table.SetCursor(len(rows) - 1)
	}
}

func (m *tableModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.table.SetWidth(60)
	m.table.SetHeight(height - 18)
}

func (m *tableModel) UpdateTasks(tasks []database.Task) {
	m.tasks = tasks
	m.filteredTasks = tasks
	m.updateTableRows()
}

func (m *tableModel) filterTasks(filter string) {
	if filter == "" {
		m.filteredTasks = m.tasks
	} else {
		m.filteredTasks = []database.Task{}
		for _, task := range m.tasks {
			if strings.Contains(strings.ToLower(task.Name), strings.ToLower(filter)) ||
				strings.Contains(strings.ToLower(task.Description), strings.ToLower(filter)) ||
				strings.Contains(strings.ToLower(task.Assignee), strings.ToLower(filter)) {
				m.filteredTasks = append(m.filteredTasks, task)
			}
		}
	}
	m.currentPage = 1
	m.updateTableRows()
}

func (m *tableModel) totalPages() int {
	return (len(m.filteredTasks) + m.itemsPerPage - 1) / m.itemsPerPage
}

func (m tableModel) SelectedTask() *database.Task {
	index := m.table.Cursor() + (m.currentPage-1)*m.itemsPerPage
	if index >= 0 && index < len(m.filteredTasks) {
		return &m.filteredTasks[index]
	}
	return &database.Task{} // Return an empty task if index is out of bounds
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m *tableModel) SetFocus() {
	m.table.Focus()
	m.focused = true
}

func (m *tableModel) SetBlur() {
	m.table.Blur()
	m.focused = false
}

func (m tableModel) Update(msg tea.Msg) (tableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("left"))):
			if m.currentPage > 1 {
				m.currentPage--
				m.updateTableRows()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("right"))):
			if m.currentPage < m.totalPages() {
				m.currentPage++
				m.updateTableRows()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
			if len(m.filteredTasks) > 0 {
				index := m.table.Cursor() + (m.currentPage-1)*m.itemsPerPage
				m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
				m.filteredTasks = append(m.filteredTasks[:index], m.filteredTasks[index+1:]...)
				m.updateTableRows()
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	pagination := fmt.Sprintf("Page %d of %d", m.currentPage, m.totalPages())

	table := focusedStyle.Render(m.table.View())
	if !m.focused {
		table = blurredStyle.Render(m.table.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		table,
		paginationStyle.Render(pagination),
	)
}
