package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var help = cobra.Command{
	Use:   "help",
	Short: "Help",
	Run: func(cmd *cobra.Command, args []string) {
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
			"\n",
		}

		help := strings.Join(items, "\n")
		fmt.Print(help)
	},
}

func init() {
	root.AddCommand(&help)
}
