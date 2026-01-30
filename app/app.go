package app

import (
	"github.com/spf13/cobra"
)

type App struct {
	cli *cobra.Command
}

func New() *App {
	return &App{}
}

func (a *App) Cli() *cobra.Command {
	return a.cli
}

func (a *App) SetCli(cli *cobra.Command) {
	a.cli = cli
}
