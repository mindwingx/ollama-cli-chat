package flow

import (
	"context"
	"fmt"
	"ollama-cli/internal/core/domain"
	"strings"
	"time"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
)

type SendMessage struct {
	next IFlow
}

func newSendMessage() *SendMessage {
	return &SendMessage{}
}

func (s *SendMessage) Process(p *Process) {
	for {
		sendMessage(s, p)

		if p.Stop() == true {
			break
		}
	}

	p.std.ClearScreen()
	fmt.Println(printInsights(p))
	return
}

func (s *SendMessage) Next(next IFlow) {
	s.next = next
}

// HELPERS

func sendMessage(s *SendMessage, p *Process) {
	reply := p.std.SetEmoji("▶️ ").Ask("Ask: \n", true)

	if strings.HasPrefix(reply, "\\") && len(reply) <= 2 {
		helper(s, p, reply)
		return
	}

	//

	ctx, cancel := context.WithCancel(context.Background())
	go requestLoader(ctx)

	// this step helps to chain questions and answers to keep the concept seamless

	proof := reply

	if p.LastQA() != "" {
		proof = strings.Join([]string{
			"Context is limited to the Primary question and the Last question–answer pair below.",
			"Answer the new question using only this context. Do not add new information.\n",
			fmt.Sprintf("Primary question:%s\n%s\nNew question:\n%s", p.Question(), p.LastQA(), reply),
		}, "\n")
	}

	//

	answer, err := p.api.SendChatMessage(p.Model(), p.Role(), proof, p.Stream())
	fmt.Print(showCursor)
	cancel()

	if err != nil {
		p.SetErr(err)
		return
	}

	//

	{
		var qa string

		if p.Question() == "" {
			p.SetQuestion(reply)
			qa = fmt.Sprintf("last answer:\n%s", answer.Content())
		} else {
			qa = fmt.Sprintf("last question:\n%s\nlast answer:\n%s", reply, answer.Content())
		}

		p.SetLastQA(qa)
	}

	//

	printChatResponse(p, answer)

	return
}

// helper act as factory method(merged with CoR)
func helper(s *SendMessage, p *Process, reply string) {
	switch reply {
	case "\\h":
		items := []string{
			"For multiline: Start with .. [enter], end with .. [enter]\n",
			"\\h: show help",
			"\\q: quit the session",
			"\\i: show chat information",
			"\\c: show the currently selected model",
			"\\m: select a different model",
			"\\p: pull/download an available llm model",
			"\\d: delete an existing model",
			"\\n: start a new session (clears last q/a)",
			strings.Repeat("-", 10),
			"\n",
		}

		help := strings.Join(items, "\n")
		fmt.Print(help)

		return
	case "\\q":
		p.SetStop(true)
		return
	case "\\i":
		info := fmt.Sprintf("%s\n%s\n",
			printInsights(p),
			strings.Repeat("-", 10),
		)

		fmt.Print(info)
		return
	case "\\c":
		currentModel := fmt.Sprintf("🐬 current model: %s\n%s\n\n",
			p.Model(), strings.Repeat("-", 10))

		fmt.Print(currentModel)
		return
	case "\\m":
		getModels := p.Flows().GetModels()
		s.Next(&getModels)
		s.next.Process(p)
		return
	case "\\p":
		pullModel := p.Flows().PullModel()
		s.Next(&pullModel)
		s.next.Process(p)
		return
	case "\\d":
		deleteModel := p.Flows().DeleteModel()
		s.Next(&deleteModel)
		s.next.Process(p)
		return
	case "\\n":
		p.SetQuestion("")
		p.SetLastQA("")
		return
	default:
		p.std.Err("invalid command")
		return
	}
}

func requestLoader(ctx context.Context) {
	dotsCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print(hideCursor)
			fmt.Print(fmt.Sprintf("\r♻️ Thinking%s ", strings.Repeat(".", dotsCount)))

			time.Sleep(300 * time.Millisecond)

			if dotsCount == 3 {
				dotsCount = 0
			} else {
				dotsCount++
			}
		}
	}
}

func printChatResponse(p *Process, answer domain.Message) {
	unixDuration := int64(answer.TotalDuration())
	duration := time.Duration(unixDuration) * time.Nanosecond

	//

	p.Info().AddDuration(duration)
	p.Info().IncreaseQuestionsCount()

	//

	fmt.Printf(
		"\r✅ (%s | time: %s - duration: %.2f seconds)\n\n%s\n%s\n\n",
		p.Name(),
		answer.CreatedAt().Format("15-04-05"),
		duration.Seconds(),
		answer.Content(),
		strings.Repeat("-", 10),
	)
}

func printInsights(p *Process) string {
	timeUnit := "min(s)"
	duration := p.Info().TotalResponseDuration().Minutes()

	if duration < 1 {
		timeUnit = "second(s)"
		duration = p.Info().TotalResponseDuration().Seconds()
	}

	currentModel := fmt.Sprintf(
		"👀 Info(questions: %d | responses duration: %.2f %s)",
		p.Info().TotalQuestionsCount(),
		duration, timeUnit,
	)
	return currentModel
}
