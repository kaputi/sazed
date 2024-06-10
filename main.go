package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focused int

const margin = 4

const (
	treeView focused = iota
	listView
	snippetView
)

var (
	columnStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder())

	treeStyle           lipgloss.Style
	focusedTreeStyle    lipgloss.Style
	listStyle           lipgloss.Style
	focusedListStyle    lipgloss.Style
	snippetStyle        lipgloss.Style
	focusedSnippetStyle lipgloss.Style
)

type mainModel struct {
	help    help.Model
	focused focused
	tree    treeModel
	list    treeModel
	snippet treeModel
	quiting bool
	loaded  bool
}

func newModel() mainModel {
	help := help.New()
	help.ShowAll = true
	return mainModel{
		help:    help,
		focused: treeView,
		tree: initTreeModel(
			[]string{
				"/abc/def/ghi",
				"/home/eduardo/pupu",
			},
		),
		list: initTreeModel(
			[]string{
				"/Hello/World",
				"/Hello/World/Again",
			},
		),
		snippet: initTreeModel(
			[]string{
				"/lala",
				"/lalalala",
			},
		),
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.tree.Init(), m.list.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width - margin
		treeWidth := int(float32(msg.Width) * 0.2)
		treeStyle = columnStyle.Width(treeWidth).Height(msg.Height - 2)
		focusedTreeStyle = treeStyle.BorderForeground(lipgloss.Color("#008080"))

		listWidth := int(float32(msg.Width) * 0.2)
		listStyle = columnStyle.Width(listWidth).Height(msg.Height - 2)
		focusedListStyle = listStyle.BorderForeground(lipgloss.Color("#008080"))

		snippetWidth := int(float32(msg.Width)*0.6) - 4
		snippetStyle = columnStyle.Width(snippetWidth).Height(msg.Height - 2)
		focusedSnippetStyle = snippetStyle.BorderForeground(lipgloss.Color("#008080"))

		m.loaded = true

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ", "l", "right":
			m.focused++
			if m.focused > snippetView {
				m.focused = snippetView
			}
		case "h", "left", "backspace", "esc":
			m.focused--
			if m.focused < treeView {
				m.focused = treeView
			}
		}

		switch m.focused {
		case treeView:
			newModel, _ := m.tree.Update(msg)
			m.tree = newModel.(treeModel)
		case listView:
			newModel, _ := m.list.Update(msg)
			m.list = newModel.(treeModel)
		case snippetView:
			newModel, _ := m.snippet.Update(msg)
			m.snippet = newModel.(treeModel)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if m.quiting {
		return ""
	}

	if !m.loaded {
		return "Loading..."
	}

	s := ""

	switch m.focused {
	case treeView:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedTreeStyle.Render(m.tree.View()),
			listStyle.Render(m.list.View()),
			snippetStyle.Render(m.snippet.View()),
		)
	case listView:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			treeStyle.Render(m.tree.View()),
			focusedListStyle.Render(m.list.View()),
			snippetStyle.Render(m.snippet.View()),
		)
	case snippetView:
		s += lipgloss.JoinHorizontal(
			lipgloss.Left,
			treeStyle.Render(m.tree.View()),
			listStyle.Render(m.list.View()),
			focusedSnippetStyle.Render(m.snippet.View()),
		)
	}

	return s
}

func main() {
	// clear the terminal
	fmt.Print("\033[H\033[2J")
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
