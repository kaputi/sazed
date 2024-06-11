package main

import (
	"fmt"
	"os"
	"sazed/config"
	"sazed/utils"
	"sazed/views"
	"strings"

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
	tree    views.TreeModel
	list    views.ListModel
	snippet views.TreeModel
	quiting bool
	loaded  bool
	config  config.Config
}

func newModel(sazedConfig config.Config) mainModel {
	help := help.New()
	help.ShowAll = true

	entries := utils.ReadDir(sazedConfig.Root())
	dirs := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	// snippetEntry := utils.ReadDir(fmt.Sprintf("%s/%s", sazedConfig.Root(), dirs[0]))[0]
	// TODO: read metadata etc, and this should run every time a directory is selected

	dir := strings.Join([]string{sazedConfig.Root(), dirs[0]}, string(os.PathSeparator))

	return mainModel{
		config:  sazedConfig,
		help:    help,
		focused: treeView,
		tree:    views.NewTree(dirs),
		list:    views.NewList(dir),
		snippet: views.NewTree(
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

		m.tree.Update(msg)

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
			m.tree = newModel.(views.TreeModel)
		case listView:
			newModel, _ := m.list.Update(msg)
			m.list = newModel.(views.ListModel)
		case snippetView:
			newModel, _ := m.snippet.Update(msg)
			m.snippet = newModel.(views.TreeModel)
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

	config := config.Load()

	utils.CreateDirIfNotExist(config.Root())
	utils.CreateDirIfNotExist(fmt.Sprintf("%s/%s", config.Root(), config.Filetype()))

	p := tea.NewProgram(newModel(config))
	_, err := p.Run()
	utils.CheckErr(err)

	// files := utils.ReadDir("/home/eduardo")
	// for _, file := range files {
	// 	fmt.Println("FILE: -----------")
	// 	fmt.Println("name", file.Name())
	// 	fmt.Println("ext", file.Ext())
	// 	fmt.Println("path", file.Path())
	// 	fmt.Println("isDir", file.IsDir())
	// }
}
