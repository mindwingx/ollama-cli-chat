package flow

import (
	"fmt"
)

type InitCheck struct {
	next IFlow
}

func newInitCheck() *InitCheck {
	return &InitCheck{}
}

func (i *InitCheck) Process(p *Process) {
	fmt.Printf("🟢 Started\n(For multiline: Start with .. [enter], end with .. [enter])\n\n")

	res, err := p.api.Handshake()
	if err != nil || res == false {
		p.SetStop(true)
		fmt.Print("🚫 Service unavailable\n")
		return
	}

	if i.next != nil {
		i.next.Process(p)
	} else {
		return
	}
}

func (i *InitCheck) Next(next IFlow) {
	i.next = next
}
