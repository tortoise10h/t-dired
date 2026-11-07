package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type entry struct {
	name     string
	isDir    bool
	fullPath string
}

type dirLoadedMsg struct {
	path    string
	entries []entry
	err     error
}

type model struct {
	width int
	height int

	cursor int
	renderStart int
	renderEnd int
	offset int

	currentDir string
	parentDir  string
	marked     map[int]struct{}
	entries    []entry
}

func InitModel() model {
	return model{
		cursor: 0,
		marked: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return loadDirCmd(".")
}

func (m model) visibleHeight() int {
	reservedLines := 9

	h := m.height - reservedLines
	if h < 1 {
		return 1
	}

	return h
}

func (m *model) refreshRenderFrame() {
	m.renderStart = 0
	m.renderEnd = min(m.visibleHeight() - 1, len(m.entries) - 1)
}

func (m model) View() tea.View {
	var s strings.Builder
	s.WriteString("T-Dired\n\n")
	fmt.Fprintf(&s, "%s (%vx%v -- %v->%v) - c:%v\n\n", m.currentDir, m.width, m.height, m.renderStart, m.renderEnd, m.cursor)

	if m.renderEnd > m.renderStart {
		for i := m.renderStart; i<=m.renderEnd; i++ {
			entry := m.entries[i]

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

	}


	s.WriteString("\nPress 'q' to quit.\n")

	v := tea.NewView(s.String())
	v.AltScreen = true

	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.refreshRenderFrame()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down", "j", "n":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
			if m.cursor > m.renderEnd && m.renderEnd < len(m.entries) - 1 {
				m.renderStart++
				m.renderEnd++
			}
		case "up", "k", "p":
			if m.cursor > 0 {
				m.cursor--
			}
			if m.cursor < m.renderStart && m.renderStart > 0 {
				m.renderStart--
				m.renderEnd--
			}
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.entries) {
				selectedEntry := m.entries[m.cursor]

				// Directory entry? zoom in it
				if selectedEntry.isDir {
					return m, loadDirCmd(selectedEntry.fullPath)
				}

				// File entry? open it
				return m, openFileCmd(selectedEntry.fullPath)
			}

			return m, nil
		case "^":
			return m, loadDirCmd(m.parentDir)
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

		m.refreshRenderFrame()

		return m, nil
	}

	return m, nil
}
