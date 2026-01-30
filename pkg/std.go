package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Std struct {
	reader *bufio.Reader
	emoji  string
}

func NewStd() *Std {
	return &Std{reader: bufio.NewReader(os.Stdin)}
}

// ClearScreen goes to new line and previous lines will be scrolled to top
func (s *Std) ClearScreen() *Std {
	fmt.Print("\033[2J\033[H")
	return s
}

// JumpBreak goes to new line and shows the last upper line too
func (s *Std) JumpBreak() *Std {
	fmt.Print("\x1b[2J\x1b[H")
	return s
}

func (s *Std) SetEmoji(emoji string) *Std {
	s.emoji = emoji
	return s
}

func (s *Std) Ask(question string, jump bool, params ...any) string {
	emoji := "🎯"

	if s.emoji != "" {
		emoji = s.emoji
	}

	q := fmt.Sprintf("%s %s ", emoji, question)

	if len(params) > 0 {
		fmt.Printf(q, params...)
	} else {
		fmt.Print(q)
	}

	s.emoji = ""

	var input string

	std, _ := s.reader.ReadString('\n')
	val := strings.TrimSpace(std)

	// the jump-break used here to be triggered if any question asked
	if jump == true {
		s.JumpBreak()
	}

	if val == ".." {
		for {
			std, _ = s.reader.ReadString('\n')
			val = strings.TrimSpace(std)

			if val == ".." {
				s.JumpBreak()
				break
			}

			input += "\n" + val
		}
	} else {
		input = val
	}

	return input
}

func (s *Std) Err(err string) {
	emoji := "⚠️"

	if s.emoji != "" {
		emoji = s.emoji
	}

	s.ClearScreen()

	fmt.Print(fmt.Sprintf("%s %s\n\n", emoji, err))
	s.emoji = ""

	time.Sleep(1500 * time.Millisecond)
}
