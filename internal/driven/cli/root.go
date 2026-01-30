package cli

import (
	"fmt"
	svc "ollama-cli/internal/core/service"
	"ollama-cli/internal/driven/cli/flow"

	"github.com/spf13/cobra"
)

var root = cobra.Command{
	Short: "Ollama Chat Cli",
	Run: func(cmd *cobra.Command, args []string) {
		api := svc.NewOllamaChat()
		chat := flow.New(api)

		var (
			process     = chat.Flows()
			general     = process.General()
			initCheck   = process.InitCheck()
			setModel    = process.GetModels()
			sendMessage = process.SendMessage()
		)

		{
			initCheck.Next(&setModel)
			setModel.Next(&sendMessage)
			sendMessage.Next(&general)
		}

		initCheck.Process(chat)

		fmt.Printf("⚠️ Stopped\n")
	},
}

func New() *cobra.Command {
	return &root
}
