package main

import (
	"flag"
	"os"

	"github.com/spf13/afero"
	"go.followtheprocess.codes/goignore/cli"
	"go.followtheprocess.codes/msg"
)

var (
	help    bool // The --help flag
	version bool // The --version flag
	list    bool // The --list flag
)

func main() {
	// Get cwd upfront
	cwd, err := os.Getwd()
	if err != nil {
		msg.Error("%s", err)
		os.Exit(1)
	}
	// Set up flags
	flag.BoolVar(&help, "help", false, "--help")
	flag.BoolVar(&version, "version", false, "--version")
	flag.BoolVar(&list, "list", false, "--list")

	app := cli.New(os.Stdout, os.Stderr, afero.NewOsFs())

	flag.Usage = app.Help

	flag.Parse()

	// Must pass at least 1 argument
	if flag.NArg() < 1 && (!help && !version && !list) {
		msg.Error("%s", err)
		os.Exit(1)
	}

	switch {
	case help:
		app.Help()

	case version:
		app.Version()

	case list:
		app.List()

	default:
		// Run the actual program
		if err := app.Run(cwd, flag.Args()); err != nil {
			msg.Error("%s", err)
			os.Exit(1)
		}
	}
}
