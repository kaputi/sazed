package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focused int

const (
	tree focused = iota
	list
	snippets
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(10).
			BorderStyle(lipgloss.NormalBorder())

	activeModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(10).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("205"))
)

// type treeModel struct {
// }

// type listModel struct {
// }

// type snippetModel struct {
// }

type testModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initTestModel() testModel {
	return testModel{
		choices:  []string{"one", "two", "three", "four", "five"},
		selected: make(map[int]struct{}),
	}
}

func (m testModel) Init() tea.Cmd {
	return nil
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			if m.cursor > 0 {
				m.cursor--
			}
		case "k", "up":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ", "l":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m testModel) View() string {
	s := ""

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "->"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

type mainModel struct {
	active  focused
	tree    treeModel
	list    testModel
	snippet testModel
}

func newModel() mainModel {
	return mainModel{
		active: tree,
		tree: initTreeModel(
			[]string{
				"/abc/def/ghi",
				"/home/eduardo/pupu",
			},
		),
		list:    initTestModel(),
		snippet: initTestModel(),
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.tree.Init(), m.list.Init(), m.snippet.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// fmt.Println(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "down", "j":
			switch m.active {
			case tree:
				m.tree, cmd = m.tree.Update(msg)
				// m.tree, cmd = m.tree.Update(msg)
				// cmds = append(cmds, cmd)
			}
		}

		switch m.active {
		case tree:
			// append commands
		case list:
			// append commands
		case snippets:
			// append commands
		}
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	s := ""
	// model := m.currentFocusedModel()

	switch m.active {
	case tree:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			activeModelStyle.Render(m.tree.View()),
			modelStyle.Render(m.list.View()),
			modelStyle.Render(m.snippet.View()),
		)
	case list:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			modelStyle.Render(m.tree.View()),
			activeModelStyle.Render(m.list.View()),
			modelStyle.Render(m.snippet.View()),
		)
	case snippets:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			modelStyle.Render(m.tree.View()),
			modelStyle.Render(m.list.View()),
			modelStyle.Render(m.snippet.View()),
		)
	}

	return s
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
