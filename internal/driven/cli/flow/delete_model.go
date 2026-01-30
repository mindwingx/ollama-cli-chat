package flow

import (
	"fmt"
	"strconv"
)

type DeleteModel struct {
	next IFlow
}

func newDeleteModel() *DeleteModel {
	return &DeleteModel{}
}

func (d *DeleteModel) Process(p *Process) {
	models, err := p.api.GetModelsList()
	if err != nil {
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

	reply := p.std.Ask("Enter the model number to delete: ", false)

	var model string

	if reply != "" {
		if index, convErr := strconv.Atoi(reply); convErr == nil && index >= 0 && index <= len(models.List()) {
			model = models.List()[index-1].Model

			if p.Model() == model {
				p.std.Err("This model is in use!")
				return
			}
		} else {
			p.std.Err("Invalid model!")
			return
		}
	}

	//

	err = p.api.DeleteModel(model)
	if err != nil {
		p.std.SetEmoji("🚫").Err(err.Error())
		return
	}

	//

	p.std.ClearScreen()
	fmt.Print("")

	//

	if d.next != nil {
		d.next.Process(p)
	} else {
		return
	}
}

func (d *DeleteModel) Next(next IFlow) {
	d.next = next
}
