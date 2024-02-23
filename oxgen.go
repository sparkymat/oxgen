package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/sparkymat/oxgen/internal/generator"
)

func main() {
	app := OxgenApp{}
	ctx := kong.Parse(
		&app,
		kong.Name("oxgen"),
		kong.Description("A web-app generator"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}

type OxgenApp struct {
	Setup SetupCommand `cmd:"" help:"Initialize a new project"`
}

type SetupCommand struct {
	Name  string `required:"" help:"Name of the project"`
	Force bool   `help:"Forcibly initialize even if the folder is not empty."`
}

func (i *SetupCommand) Run(ctx *kong.Context) error {
	s := generator.New()
	return s.Setup(context.Background(), i.Name, i.Force)
}
