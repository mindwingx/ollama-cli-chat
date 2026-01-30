package service

import (
	"errors"
	"fmt"
	prv "ollama-cli/internal/adapter/provider"
	"ollama-cli/internal/core/domain"
	svc "ollama-cli/internal/core/port"
	"ollama-cli/pkg"
	"strings"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
)

type Chat struct {
	std *pkg.Std
	api prv.IOllamaProvider
}

func NewOllamaChat() svc.IChatService {
	return &Chat{std: pkg.NewStd(), api: prv.NewOllamaProvider()}
}

func (c *Chat) Handshake() (res bool, err error) {
	resp, err := c.api.Handshake()
	if err != nil {
		return
	}

	if resp != "" {
		res = true
	}

	return
}

func (c *Chat) SendChatMessage(model, role, msg string, stream bool) (res domain.Message, err error) {
	payload := prv.ChatRequest{Model: model, Stream: stream}
	payload.Messages = append(payload.Messages, prv.Message{Role: role, Content: msg})

	resp, err := c.api.SendMessage(payload)
	if err != nil {
		return
	}

	res = *domain.NewMessage()
	res.SetContent(resp.Message.Content)
	res.SetCreatedAt(resp.CreatedAt)
	res.SetTotalDuration(resp.TotalDuration)

	return
}

func (c *Chat) GetModelsList() (res domain.Models, err error) {
	resp, err := c.api.GetModels()
	if err != nil {
		return
	}

	if len(resp.Models) == 0 {
		err = errors.New("no model found")
		return
	}

	items := make([]domain.Model, 0)

	for _, m := range resp.Models {
		items = append(items, domain.Model{Model: m.Model, Name: m.Name})
	}

	res = *domain.NewModels()
	res.SetList(items)

	return
}

func (c *Chat) PullModel(model string) (err error) {
	fmt.Println(hideCursor)
	defer fmt.Println(showCursor)

	payload := prv.ModelRequest{Model: model}

	ctx, streamChan, errChan := c.api.PullModel(payload)

	var (
		status   string  = ""
		total    int64   = 0
		progress float64 = 0
	)

	for {

		select {
		case err = <-errChan:
			if err != nil {
				fmt.Print("\r")
				c.std.Err(err.Error())
			}

			return
		case resp := <-streamChan:
			statusChunks := strings.Split(resp.Status, " ")
			status = statusChunks[0]

			if strings.HasPrefix(resp.Status, "pulling") == true {
				if total == 0 {
					total = toMb(resp.Total)
				}

				if total > 0 && total == toMb(resp.Total) && resp.Completed > 0 {
					progress = float64((resp.Completed * 100) / resp.Total)
				}
			}

			if total > 0 && (total == toMb(resp.Completed) || status == "success") {
				progress = 100.00
			}

			if status == "" {
				status = "process"
			}

			fmt.Print(fmt.Sprintf("\r⚓️ %s - %dMB(%.f%%)", status, total, progress))
		case <-ctx.Done():
			fmt.Print("\r")
			fmt.Print(fmt.Sprintf("\r⚓️ success - %dMB(100%%)%s", total, strings.Repeat(" ", 10)))
			fmt.Print("\n\n")
			return
		}
	}
}

func (c *Chat) DeleteModel(model string) (err error) {
	payload := prv.ModelRequest{Model: model}
	err = c.api.DeleteModel(payload)
	return
}

// HELPERS

func toMb(bytes int64) int64 {
	return bytes / (1024 * 1204)
}
