package flow

import (
	"fmt"
	"strconv"
	"time"
)

type SetModel struct {
	next IFlow
}

func newSetModel() *SetModel {
	return &SetModel{}
}

func (s *SetModel) Process(p *Process) {
	models, err := p.api.GetModelsList()
	if err != nil {
		p.SetStop(true)
		fmt.Print(fmt.Sprintf("🚫 %s\n", err.Error()))
		return
	}

	//
	fmt.Print("📌 Available LLMs:\n")

	for i, model := range models.List() {
		fmt.Printf(fmt.Sprintf("%d) %s\n", i+1, model.Name))
	}

	fmt.Println()

	//

	reply := p.std.Ask("Enter the model number to set: ", false)

	{
		model := models.List()[0]

		if reply != "" {
			if index, convErr := strconv.Atoi(reply); convErr == nil && index >= 0 && index <= len(models.List()) {
				model = models.List()[index-1]
			} else {
				p.std.ClearScreen()
				fmt.Print(fmt.Sprintf("⚠️ Set model failed. default model(%s) set\n", model.Name))
				time.Sleep(1500 * time.Millisecond)
			}
		}

		p.SetModel(model.Model)
		p.SetName(model.Name)
	}

	p.std.ClearScreen()
	fmt.Println()

	if s.next != nil {
		s.next.Process(p)
	} else {
		return
	}
}

func (s *SetModel) Next(next IFlow) {
	s.next = next
}
