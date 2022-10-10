// Package cli implements the CLI functionality
// main defers execution to the exported methods in this package
package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/FollowTheProcess/msg"
	"github.com/fatih/color"
	"github.com/spf13/afero"
)

var (
	version = "dev" // goignore version, set at compile time with ldflags
	commit  = ""    // goignore commit hash, set at compile time with ldflags
)

const (
	// The base URL for gitignore.io.
	ignoreURL = "https://www.toptal.com/developers/gitignore/api"
	helpText  = `
Generate great gitignore files, straight from the command line! 🛠

Usage:
	goignore [targets] [flags]

Examples:

# Add gitignore targets as arguments
$ goignore macos visualstudiocode go

# See a list of allowed targets
$ goignore --list

Flags:
	-h, --help      Help for venv
	-v, --version   Show venv's version info
	-l, --list      Show a list of all allowed targets`
)

// App represents the goignore CLI program.
type App struct {
	stdout  io.Writer
	stderr  io.Writer
	fs      afero.Afero
	printer *msg.Printer
}

// New creates and returns a new App configured with an afero file system
// and IO streams.
func New(stdout, stderr io.Writer, fs afero.Fs) *App {
	af := afero.Afero{Fs: fs}
	printer := msg.Default()
	printer.Stdout = stdout
	printer.Stderr = stderr
	return &App{stdout: stdout, stderr: stderr, fs: af, printer: printer}
}

// Help prints the CLI help text.
func (a *App) Help() {
	fmt.Fprintln(a.stdout, helpText)
}

// Version prints CLI version info.
func (a *App) Version() {
	ver := color.CyanString("goignore version")
	sha := color.CyanString("commit")

	fmt.Fprintf(a.stdout, "%s: %s\n", ver, version)
	fmt.Fprintf(a.stdout, "%s: %s\n", sha, commit)
}

// List prints the list of valid gitignore targets.
func (a *App) List() {
	for _, target := range targets {
		fmt.Fprintln(a.stdout, target)
	}
}

// Run represents the entry point to the CLI.
func (a *App) Run(cwd string, args []string) error {
	for _, arg := range args {
		if !IsValidTarget(arg) {
			return fmt.Errorf("%q is not a valid gitignore target", arg)
		}
	}

	a.printer.Infof("Generating gitignore for %v", strings.Join(args, ", "))

	data, err := getIgnoreData(ignoreURL, args)
	if err != nil {
		return err
	}

	if err := a.writeToIgnoreFile(cwd, data); err != nil {
		return err
	}

	a.printer.Good("Done")

	return nil
}
