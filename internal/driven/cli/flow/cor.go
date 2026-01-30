package flow

import (
	"ollama-cli/internal/core/port"
	"ollama-cli/pkg"
	"time"
)

// CoR

const modelDefaultRole = "user"

type IFlow interface {
	Process(p *Process)
	Next(next IFlow)
}

type (
	Process struct {
		std      pkg.Std
		api      port.IChatService
		stop     bool
		model    string
		name     string
		role     string
		message  string
		question string
		lastQA   string
		stream   bool
		err      error
		info     Info
		flows    Flows
	}

	Info struct {
		questionsCount        int
		totalResponseDuration time.Duration
	}

	Flows struct {
		general *General
		init    *InitCheck
		model   *SetModel
		message *SendMessage
		pull    *PullModel
		delete  *DeleteModel
	}
)

func New(api port.IChatService) *Process {
	return &Process{
		std:    *pkg.NewStd(),
		api:    api,
		role:   modelDefaultRole,
		stream: false,
		info:   Info{questionsCount: 0, totalResponseDuration: time.Duration(0)},
		flows: Flows{
			general: newGeneral(),
			init:    newInitCheck(),
			model:   newSetModel(),
			message: newSendMessage(),
			pull:    newPullModel(),
			delete:  newDeleteModel(),
		},
	}
}

func (p *Process) Stop() bool {
	return p.stop
}

func (p *Process) SetStop(stop bool) {
	p.stop = stop
}

func (p *Process) Model() string {
	return p.model
}

func (p *Process) SetModel(model string) {
	p.model = model
}

func (p *Process) Name() string {
	return p.name
}

func (p *Process) SetName(name string) {
	p.name = name
}

func (p *Process) Role() string {
	return p.role
}

func (p *Process) SetRole(role string) {
	p.role = role
}

func (p *Process) Message() string {
	return p.message
}

func (p *Process) SetMessage(message string) {
	p.message = message
}

func (p *Process) Question() string {
	return p.question
}

func (p *Process) SetQuestion(q string) {
	p.question = q
}

func (p *Process) LastQA() string {
	return p.lastQA
}

func (p *Process) SetLastQA(qa string) {
	p.lastQA = qa
}

func (p *Process) Stream() bool {
	return p.stream
}

func (p *Process) SetStream(stream bool) {
	p.stream = stream
}

func (p *Process) Err() error {
	return p.err
}

func (p *Process) SetErr(err error) {
	p.err = err
}

func (p *Process) Info() *Info {
	return &p.info
}

func (p *Process) Flows() *Flows {
	return &p.flows
}

// info

func (i *Info) TotalQuestionsCount() int {
	return i.questionsCount
}

func (i *Info) IncreaseQuestionsCount() {
	i.questionsCount += 1
}

func (i *Info) TotalResponseDuration() time.Duration {
	return i.totalResponseDuration
}

func (i *Info) AddDuration(d time.Duration) {
	i.totalResponseDuration += d
}

// flows

func (f *Flows) General() General {
	return *f.general
}

func (f *Flows) InitCheck() InitCheck {
	return *f.init
}

func (f *Flows) GetModels() SetModel {
	return *f.model
}

func (f *Flows) SendMessage() SendMessage {
	return *f.message
}

func (f *Flows) PullModel() PullModel {
	return *f.pull
}

func (f *Flows) DeleteModel() DeleteModel {
	return *f.delete
}
