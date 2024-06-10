package views

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TreeModel struct {
	directories []string
	paths       map[string]string
	cursor      int
	focus       bool
}

func NewTree(pathList []string) TreeModel {
	paths := make(map[string]string)
	directories := []string{}
	for _, path := range pathList {
		parts := strings.Split(path, string(os.PathSeparator))
		dir := parts[len(parts)-1]
		paths[dir] = path
		directories = append(directories, dir)
	}
	return TreeModel{
		paths:       paths,
		directories: directories,
	}
}

func (m TreeModel) Focus() tea.Msg {
	m.focus = true
	return focusCmd(m.focus)
}

func (m TreeModel) Blur() tea.Msg {
	m.focus = false
	return focusCmd(m.focus)
}

func (m TreeModel) Path() string {
	return m.paths[m.directories[m.cursor]]
}

func (m TreeModel) Init() tea.Cmd {
	return nil
}

func (m TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			m.cursor++
			if m.cursor >= len(m.directories) {
				m.cursor = 0
			}
		case "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.directories) - 1
			}
		}
	}
	return m, nil
}

func (m TreeModel) View() string {
	s := ""
	if m.focus {
		s += focusedTitleStyle.Render("Directories:")
	} else {
		s += titleStyle.Render("Directories:")
	}

	s += "\n\n"

	for i, dir := range m.directories {
		if m.cursor == i {
			s += fmt.Sprintf("%s\n", titleStyle.Render(fmt.Sprintf("│ %s", dir)))
			s += titleStyle.Render("│")
			s += "\n"
		} else {
			s += fmt.Sprintf("  %s\n\n", dir)
		}
	}
	return s
}
