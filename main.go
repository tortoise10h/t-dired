package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var cursorStyle = lipgloss.NewStyle().Underline(true)

type entry struct {
	name     string
	isDir    bool
	fullPath string
}

type model struct {
	cursor     int
	currentDir string
	parentDir  string
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
	return loadDirEntries(".")
}

func loadDirEntries(path string) tea.Cmd {
	return func() tea.Msg {
		targetDir, err := filepath.Abs(path)
		if err != nil {
			return dirLoadedMsg{
				path: path,
				err:  err,
			}
		}

		dirEntries, err := os.ReadDir(targetDir)
		if err != nil {
			return dirLoadedMsg{
				path: path,
				err:  err,
			}
		}
		var entries []entry

		parentDir := filepath.Dir(targetDir)
		entries = append(entries, entry{
			name:     "..",
			isDir:    true,
			fullPath: parentDir,
		})

		for _, de := range dirEntries {
			entries = append(entries, entry{
				name:     de.Name(),
				isDir:    de.IsDir(),
				fullPath: filepath.Join(targetDir, de.Name()),
			})
		}

		return dirLoadedMsg{
			path:    targetDir,
			entries: entries,
		}
	}
}

func (m model) View() tea.View {
	var s strings.Builder
	s.WriteString("T-Dired\n\n")
	fmt.Fprintf(&s, "%s\n\n", m.currentDir)

	for i, entry := range m.entries {
		var renderedEntry string
		if entry.isDir {
			renderedEntry = fmt.Sprintf("%v/", entry.name)
		} else {
			renderedEntry = fmt.Sprintf("%v", entry.name)
		}

		if i == m.cursor {
			fmt.Fprintf(&s, "> %s\n", cursorStyle.Render(renderedEntry))
		} else {
			fmt.Fprintf(&s, "  %s\n", renderedEntry)
		}
	}

	s.WriteString("\nPress 'q' to quit.\n")

	v := tea.NewView(s.String())
	v.AltScreen = true

	return v
}

func openFileCmd(path string) tea.Cmd {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return func() tea.Msg {
			return fmt.Errorf("Unsupported platform\n")
		}
	}

	err := cmd.Run()
	if err != nil {
		return func() tea.Msg {
			return fmt.Errorf("Something went wrong when opening file %v\n", err)
		}
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down", "j", "n":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "up", "k", "p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.entries) {
				selectedEntry := m.entries[m.cursor]

				// Directory entry? zoom in it
				if selectedEntry.isDir {
					return m, loadDirEntries(selectedEntry.fullPath)
				}

				// File entry? open it
				return m, openFileCmd(selectedEntry.fullPath)
			}

			return m, nil
		case "^":
			return m, loadDirEntries(m.parentDir)
		}
	case dirLoadedMsg:
		if msg.err != nil {
			return m, nil
		}

		m.currentDir = msg.path
		m.parentDir = filepath.Dir(msg.path)
		m.entries = msg.entries

		// if there is no entries in the current dir, the cursor stay at the .. entry but if there is at least 1 entry, the cursor should stay at the next item after .. for better ux
		if len(m.entries) > 1 {
			m.cursor = 1
		} else {
			m.cursor = 0
		}

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
