package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type treeModel struct {
	directories []string
	paths       map[string]string
	cursor      int
	focus       bool
}

func initTreeModel(pathList []string) treeModel {
	paths := make(map[string]string)
	directories := []string{}
	for _, path := range pathList {
		parts := strings.Split(path, string(os.PathSeparator))
		dir := parts[len(parts)-1]
		paths[dir] = path
		directories = append(directories, dir)
	}
	return treeModel{
		paths:       paths,
		directories: directories,
	}
}

func (m treeModel) Init() tea.Cmd {
	return nil
}

func (m treeModel) Focus() {
	m.focus = true
}

func (m treeModel) Blur() {
	m.focus = false
}

func (m treeModel) Focused() bool {
	return m.focus
}

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m treeModel) View() string {
	s := ""
	for i, dir := range m.directories {
		cursor := "  "
		if m.cursor == i {
			cursor = "->"
		}
		s += fmt.Sprintf("%s %s\n", cursor, dir)
	}
	return s
}
