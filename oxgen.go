package main

import (
	"fmt"

	"github.com/alecthomas/kong"
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
	Init InitCommand `cmd:"" help:"Initialize a new project"`
}

type InitCommand struct {
	Force bool `help:"Forcibly initialize even if the folder is not empty."`
}

func (i *InitCommand) Run(ctx *kong.Context) error {
	fmt.Println("I won't do that")
	return nil
}
