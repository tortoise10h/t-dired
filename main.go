package main

import (
	"strings"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct{
	cursor int
	marked map[int]struct{}
	files []string
}

func initModel() model {
	return model{
		cursor: 0,
		marked: make(map[int]struct{}),
		files: []string{"..", "main.go", "README.md"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() tea.View {
	var s strings.Builder; s.WriteString("T-Dired\n\n")

	for _, file := range m.files {
		fmt.Fprintf(&s, "%v\n", file)
	}

	s.WriteString("\nPress 'q' to quit.\n")

	return tea.NewView(s.String())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func main() {
	p := tea.NewProgram(initModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There are some error on startup %v\n", err)
		os.Exit(1)
	}
}
