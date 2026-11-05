package main

import (
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
	return tea.NewView("Hello, this is the first view")
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
