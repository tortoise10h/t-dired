package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type entry struct {
	name     string
	isDir    bool
	fullPath string
}

type model struct {
	cursor     int
	currentDir string
	marked     map[int]struct{}
	entries    []entry
}

type dirLoadedMsg struct {
	path    string
	entries []entry
	err     error
}

func initModel() model {
	return model{
		cursor: 0,
		marked: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return loadDirFiles()
}

func loadDirFiles() tea.Cmd {
	return func() tea.Msg {
		cwd, err := os.Getwd()
		if err != nil {
			return dirLoadedMsg{
				err: err,
			}
		}

		dirEntries, err := os.ReadDir(cwd)
		if err != nil {
			return dirLoadedMsg{
				err: err,
			}
		}

		var entries []entry

		parentDir := filepath.Dir(cwd)
		entries = append(entries, entry{
			name:     "..",
			isDir:    true,
			fullPath: parentDir,
		})

		for _, de := range dirEntries {
			entries = append(entries, entry{
				name:     de.Name(),
				isDir:    de.IsDir(),
				fullPath: filepath.Join(cwd, de.Name()),
			})
		}

		return dirLoadedMsg{
			path:    cwd,
			entries: entries,
		}
	}
}

func (m model) View() tea.View {
	var s strings.Builder
	s.WriteString("T-Dired\n\n")

	for _, entry := range m.entries {
		if entry.isDir {
			fmt.Fprintf(&s, "%v/\n", entry.name)
		} else {
			fmt.Fprintf(&s, "%v\n", entry.name)
		}
	}

	s.WriteString("\nPress 'q' to quit.\n")

	v := tea.NewView(s.String())
	v.AltScreen = true

	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case dirLoadedMsg:
		if msg.err != nil {
			return m, nil
		}

		m.currentDir = msg.path
		m.entries = msg.entries
		return m, nil
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
