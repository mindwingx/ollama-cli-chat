package flow

type General struct {
	next IFlow
}

func newGeneral() *General {
	return &General{}
}

func (g *General) Process(p *Process) {
	if g.next != nil {
		g.next.Process(p)
	} else {
		return
	}
}

func (g *General) Next(next IFlow) {
	g.next = next
}
