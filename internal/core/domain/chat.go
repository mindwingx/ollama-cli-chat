package domain

import "time"

type Message struct {
	content       string
	createdAt     time.Time
	totalDuration int
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) Content() string {
	return m.content
}

func (m *Message) SetContent(content string) {
	m.content = content
}

func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) SetCreatedAt(at time.Time) {
	m.createdAt = at
}

func (m *Message) TotalDuration() int {
	return m.totalDuration
}

func (m *Message) SetTotalDuration(duration int) {
	m.totalDuration = duration
}

//

type (
	Models struct {
		list []Model
	}

	Model struct {
		Model string
		Name  string
	}
)

func NewModels() *Models {
	return &Models{}
}

func (m *Models) List() []Model {
	return m.list
}

func (m *Models) SetList(list []Model) {
	m.list = list
}
