package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type entry struct {
	name     string
	isDir    bool
	fullPath string
	mode     string
	size     int64
	modTime  time.Time

	metadataRaw string
	displayName string
}

type dirLoadedMsg struct {
	path    string
	entries []entry
	err     error
}

type mainScreenModel struct {
	width, height int

	cursor      int
	renderStart int
	renderEnd   int

	currentDir string
	parentDir  string
	marked     map[int]struct{}
	entries    []entry

	deleteConfirming bool
}

type openConfirmationPopupMsg struct {
	target entry
}

func NewMainScreen() mainScreenModel {
	return mainScreenModel{
		deleteConfirming: false,
		marked:           make(map[int]struct{}),
	}
}

func (m mainScreenModel) visibleHeight() int {
	reservedLines := 9

	h := m.height - reservedLines
	if h < 1 {
		return 1
	}

	return h
}

func (m *mainScreenModel) refreshRenderFrame() {
	m.renderStart = 0
	m.renderEnd = min(m.visibleHeight()-1, len(m.entries)-1)
}

func (m *mainScreenModel) activateDeleteState() {
	m.deleteConfirming = true
}

func (m *mainScreenModel) deactivateDeleteState() {
	m.deleteConfirming = false
}

func (m mainScreenModel) View() string {
	var s strings.Builder
	s.WriteString("T-Dired\n\n")
	fmt.Fprintf(&s, "%s (%vx%v -- %v->%v) - c:%v\n\n", m.currentDir, m.width, m.height, m.renderStart, m.renderEnd, m.cursor)

	if m.renderEnd > m.renderStart {
		for i := m.renderStart; i <= m.renderEnd; i++ {
			entry := m.entries[i]

			if i == m.cursor {
				s.WriteString("> ")
				s.WriteString(metaCursorStyle.Render(entry.metadataRaw))
				s.WriteString(underlineStyle.Render("  " + entry.displayName))
				s.WriteByte('\n')
			} else {
				s.WriteString(metaStyle.Render("  " + entry.metadataRaw))
				s.WriteString("  ")
				s.WriteString(entry.displayName)
				s.WriteByte('\n')
			}
		}
	}

	return s.String()
}

func (m mainScreenModel) Update(msg tea.Msg) (mainScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.refreshRenderFrame()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down", "j", "n":
			if m.deleteConfirming {
				return m, nil
			}

			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
			if m.cursor > m.renderEnd && m.renderEnd < len(m.entries)-1 {
				m.renderStart++
				m.renderEnd++
			}
		case "up", "k", "p":
			if m.deleteConfirming {
				return m, nil
			}

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
		case "D":
			if len(m.entries) < 1 {
				return m, nil
			}

			selectedEntry := m.entries[m.cursor]
			// We don't delete parent dir or current dir here bro
			if selectedEntry.name == "." || selectedEntry.name == ".." {
				return m, nil
			}

			m.activateDeleteState()

			return m, func() tea.Msg {
				return openConfirmationPopupMsg{target: selectedEntry}
			}
		}
	case dirLoadedMsg:
		if msg.err != nil {
			return m, nil
		}

		dirBeforeLoad := m.currentDir
		m.currentDir = msg.path
		m.parentDir = filepath.Dir(msg.path)
		m.entries = msg.entries

		// if there is no entries in the current dir, the cursor stay at the .. entry but if there is at least 1 entry, the cursor should stay at the next item after .. for bet ux
		if m.currentDir != dirBeforeLoad {
			if len(m.entries) > 1 {
				m.cursor = 1
			} else {
				m.cursor = 0
			}
		} else {
			if m.cursor > len(m.entries)-1 {
				m.cursor = len(m.entries) - 1
			}
		}

		m.refreshRenderFrame()

		return m, nil
	case confirmDecisionMsg:
		m.deactivateDeleteState()

		// Agree delete
		if msg.confirmed {
			if msg.target.isDir {
				os.RemoveAll(msg.target.fullPath)
			} else {
				os.Remove(msg.target.fullPath)
			}
		}

		return m, loadDirCmd(m.currentDir)
	}

	return m, nil
}
