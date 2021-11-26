package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FollowTheProcess/goignore/cli"
	"github.com/fatih/color"
	"github.com/spf13/afero"
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
		title := color.New(color.FgRed).Add(color.Bold)
		msg := color.New(color.FgWhite).Add(color.Bold)
		fmt.Fprintf(os.Stderr, "%s: %s\n", title.Sprint("error"), msg.Sprintf("%s", err))
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
	if flag.NArg() < 1 && !(help || version || list) {
		title := color.New(color.FgRed).Add(color.Bold)
		msg := color.New(color.FgWhite).Add(color.Bold)
		fmt.Fprintf(os.Stderr, "%s: %s\n", title.Sprint("error"), msg.Sprint("must pass at least 1 argument"))
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
			title := color.New(color.FgRed).Add(color.Bold)
			msg := color.New(color.FgWhite).Add(color.Bold)
			fmt.Fprintf(os.Stderr, "%s: %s\n", title.Sprint("error"), msg.Sprint(err))
			os.Exit(1)
		}
	}
}
