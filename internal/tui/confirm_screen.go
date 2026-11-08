package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type confirmScreenModel struct {
	target entry
}

type confirmDecisionMsg struct {
	target    entry
	confirmed bool
}

func NewConfirmModel(e entry) confirmScreenModel {
	return confirmScreenModel{
		target: e,
	}
}

func (c confirmScreenModel) View() string {
	targetName := c.target.name

	if c.target.isDir {
		targetName += "/"
	}

	var prompt string
	if c.target.isDir {
		prompt = fmt.Sprintf("\nDelete this directory \"%s\" and all the entries inside it? (y or n) ", targetName)
	} else {
		prompt = fmt.Sprintf("\nDelete this file \"%s\"? (y or n) ", targetName)
	}

	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("203")).
		Render(prompt)
}

func (c confirmScreenModel) Update(msg tea.Msg) (confirmScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "y", "Y":
			return c, func() tea.Msg {
				return confirmDecisionMsg{target: c.target, confirmed: true}
			}
		case "n", "N", "esc", "ctrl+g": // Emacs C-g aborts
			return c, func() tea.Msg {
				return confirmDecisionMsg{target: c.target, confirmed: false}
			}
		}
	}

	return c, nil
}
