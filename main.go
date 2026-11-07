package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tortoise10h/t-dired/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.InitModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There are some error on startup %v\n", err)
		os.Exit(1)
	}
}
