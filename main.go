package main

import (
	"fmt"
	"ollama-cli/app"
	"ollama-cli/internal/driven/cli"
)

func main() {
	service := app.New()
	service.SetCli(cli.New())

	cmd := service.Cli()
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("❌ Err: %s", err.Error())
		return
	}
}
