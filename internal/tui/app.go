package tui

import (
	tea "charm.land/bubbletea/v2"
)

type screenMode int

const (
	modeMain screenMode = iota // main dired UI
	modeConfirmDelete
)

type appModel struct {
	width, height  int

	appMode screenMode
	mainScreen mainScreenModel
}

func InitModel() appModel {
	return appModel{
		appMode: modeMain,
		mainScreen: NewMainScreen(),
	}
}

func (a appModel) Init() tea.Cmd {
	// First init start with main mode - Dired UI -> load current directory entries for view data
	return loadDirCmd(".")
}

func (a appModel) View() tea.View {
	var content string
	switch a.appMode {
	case modeMain:
		content = a.mainScreen.View()
	default:
		content = a.mainScreen.View()
	}

	// Full screen mode
	v := tea.NewView(content)
	v.AltScreen = true

	return v
}

func (a appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		a.mainScreen, cmd = a.mainScreen.Update(msg)
		return a, cmd
	}

	switch a.appMode {
	case modeMain:
		a.mainScreen, cmd = a.mainScreen.Update(msg)
		return a, cmd
	default:
		a.mainScreen, cmd = a.mainScreen.Update(msg)
		return a, cmd
	}
}
