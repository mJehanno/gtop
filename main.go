package main

import (
	"log"

	"github.com/mjehanno/gtop/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
