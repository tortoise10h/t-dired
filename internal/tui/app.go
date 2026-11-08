package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type screenMode int

const (
	modeMain screenMode = iota // main dired UI
	modeConfirm
)

type appModel struct {
	width, height int

	appMode       screenMode
	mainScreen    mainScreenModel
	confirmScreen confirmScreenModel
}

func InitModel() appModel {
	return appModel{
		appMode:    modeMain,
		mainScreen: NewMainScreen(),
	}
}

func (a appModel) Init() tea.Cmd {
	// First init start with main mode - Dired UI -> load current directory entries for view data
	return loadDirCmd(".")
}

func (a appModel) statusMessage() string {
	switch a.appMode {
	case modeMain:
		return "\nPress 'q' to quit.\n"
	case modeConfirm:
		return a.confirmScreen.View()
	default:
		return "\nPress 'q' to quit.\n"
	}
}

func (a appModel) View() tea.View {
	var content strings.Builder

	content.WriteString(a.mainScreen.View())
	content.WriteString(a.statusMessage())

	// Full screen mode
	v := tea.NewView(content.String())
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
	case tea.KeyPressMsg:
		// Highest priority keybindings for entire app, can activate anywhere
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	case openConfirmationPopupMsg:
		a.confirmScreen = NewConfirmModel(msg.target)
		a.appMode = modeConfirm
	case confirmDecisionMsg:
		a.mainScreen, cmd = a.mainScreen.Update(msg)
		a.appMode = modeMain
		return a, cmd
	}

	switch a.appMode {
	case modeMain:
		a.mainScreen, cmd = a.mainScreen.Update(msg)
		return a, cmd
	case modeConfirm:
		a.confirmScreen, cmd = a.confirmScreen.Update(msg)
		return a, cmd
	default:
		a.mainScreen, cmd = a.mainScreen.Update(msg)
		return a, cmd
	}
}
