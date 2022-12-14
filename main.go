package main

import (
	"log"

	"github.com/mjehanno/gtop/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model.InitialModel())
	if  err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
