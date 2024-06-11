package views

import (
	"fmt"
	"sazed/snippets"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	directory string
	Metadata  snippets.Metadata
	cursor    int
	focus     bool
}

func NewList(directory string) ListModel {
	metadata := snippets.NewMetadata("go", directory)
	if len(metadata.Snippets) == 0 {
		metadata.AddSnippet("test", "this is just a test")
		metadata.SetCode("test", "this is just a test")
	}
	if len(metadata.Snippets) == 1 {
		metadata.AddSnippet("YAYAYA", "hello there")
		metadata.SetCode("YAYAYA", "wawawiwa")
	}
	return ListModel{
		directory: directory,
		Metadata:  metadata,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			m.cursor++
			if m.cursor >= len(m.Metadata.Snippets) {
				m.cursor = 0
			}
		case "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.Metadata.Snippets) - 1
			}
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	s := ""

	if m.focus {
		s += focusedTitleStyle.Render(m.directory) + "\n\n"
	} else {
		s += titleStyle.Render(m.directory) + "\n\n"
	}

	for i, snippet := range m.Metadata.Snippets {
		if m.cursor == i {
			s += titleStyle.Render("│ "+snippet.Name+"    "+snippet.Date) + "\n"
			s += titleStyle.Render("│ ") + otherStyle.Render(snippet.Description) + "\n"
		} else {
			s += fmt.Sprintf("  %s    %s\n", snippet.Name, snippet.Date)
			s += fmt.Sprintf("  %s\n", snippet.Description)
		}

		s += "\n"
	}

	return s
}
