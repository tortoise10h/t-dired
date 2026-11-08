package tui

import "charm.land/lipgloss/v2"

var (
	underlineStyle = lipgloss.NewStyle().Underline(true)

	metaStyle       = lipgloss.NewStyle().Faint(true)
	metaCursorStyle = metaStyle.Underline(true)
)
