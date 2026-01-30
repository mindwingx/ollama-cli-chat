package flow

type PullModel struct {
	next IFlow
}

func newPullModel() *PullModel {
	return &PullModel{}
}

func (pm *PullModel) Process(p *Process) {
	reply := p.std.Ask("Enter the model name to pull: ", false)

	if reply != "" {
		err := p.api.PullModel(reply)
		if err != nil {
			p.std.Err(err.Error())
			return
		}
	}

	if pm.next != nil {
		pm.next.Process(p)
	} else {
		return
	}
}

func (pm *PullModel) Next(next IFlow) {
	pm.next = next
}
