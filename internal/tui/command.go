package tui

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
)

func loadDirCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, targetDir, err := loadDirEntries(path)
		return dirLoadedMsg{
			entries: entries,
			path:    targetDir,
			err:     err,
		}
	}
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
