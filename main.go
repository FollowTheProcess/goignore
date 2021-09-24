package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/FollowTheProcess/msg"
)

// Version is set by ldflags at compile time
var Version = "dev"

const (
	// ignoreURL is the base url for the gitignore API
	ignoreURL   string = "https://www.toptal.com/developers/gitignore/api"
	listMessage string = "To get a list of valid targets, run goignore --list"

	helpMessage string = `
Usage: goignore [FLAGS] [TARGETS]...

Handy CLI to generate great gitignore files.

Flags:
	--version: Display goignore's version.
	--help: Show this help message and exit.
	--list: Show the valid gitignore.io targets.

Arguments:
	TARGETS: Space separated list of gitignore.io targets.	[required]

Examples:

$ goignore macos vscode go

$ goignore --list
`
)

var (
	errIgnoreFileExists = errors.New("'.gitignore' already exists. Not doing anything")
	helpFlag            bool
	versionFlag         bool
	listFlag            bool
)

func main() {
	flag.BoolVar(&helpFlag, "help", false, "--help")
	flag.BoolVar(&versionFlag, "version", false, "--version")
	flag.BoolVar(&listFlag, "list", false, "--list")

	flag.Usage = func() {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	flag.Parse()

	run()
}

func run() {
	if flag.NArg() < 2 && !(helpFlag || listFlag || versionFlag) {
		printUsage(os.Stdout)
		os.Exit(1)
	}

	switch {
	case helpFlag:
		printUsage(os.Stdout)
		os.Exit(0)

	case versionFlag:
		printVersion(os.Stdout)
		os.Exit(0)

	case listFlag:
		printList(os.Stdout, ignoreURL)
		os.Exit(0)

	case os.Args[1] == "list":
		fmt.Println(listMessage)
		os.Exit(1)

	default:
		ignoreList := os.Args[1:]
		makeIgnoreFile(ignoreList, ignoreURL)
		os.Exit(0)
	}
}

// GetIgnore hits the gitignore.io API with targets and returns the response
func GetIgnore(targets []string, url string) ([]byte, error) {
	constructedString := strings.Join(targets, ",")

	targetURL := strings.Join([]string{url, constructedString}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response: %w", err)
	}

	return data, nil
}

// GetList returns the gitignore API response for 'list'
func GetList(url string) ([]byte, error) {
	targetURL := strings.Join([]string{url, "list"}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response: %w", err)
	}

	return data, nil
}

// WriteToIgnoreFile takes data in and writes it to cwd/.gitignore
func WriteToIgnoreFile(data []byte, filename string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working dir: %w", err)
	}

	ignoreFilePath := filepath.Join(cwd, filename)

	if _, err := os.Stat(ignoreFilePath); errors.Is(err, fs.ErrNotExist) {
		// No .gitignore, we're good to go
		_, err := os.Create(ignoreFilePath)
		if err != nil {
			return fmt.Errorf("could not create ignore file: %w", err)
		}

		file, err := os.OpenFile(ignoreFilePath, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("could not open ignore file: %w", err)
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("could not write to ignore file: %w", err)
		}
		if err := file.Sync(); err != nil {
			return fmt.Errorf("could not save ignore file: %w", err)
		}
	} else {
		return errIgnoreFileExists
	}

	return nil
}

func printUsage(where io.Writer) {
	fmt.Fprintln(where, helpMessage)
}

func printVersion(where io.Writer) {
	fmt.Fprintf(where, "goignore version: %s\n", Version)
}

func printList(where io.Writer, url string) {
	data, err := GetList(url)
	if err != nil {
		msg.Failf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Fprintln(where, string(data))
}

func makeIgnoreFile(targets []string, url string) {
	msg.Infof("Creating a gitignore for: %v", strings.Join(targets, ", "))

	data, err := GetIgnore(targets, url)
	if err != nil {
		msg.Failf("Error: %s", err)
		os.Exit(1)
	}

	// By default doesn't ignore all of .vscode/
	// personal preference
	vscode := []byte("\n.vscode/\n")
	data = append(data, vscode...)

	err = WriteToIgnoreFile(data, ".gitignore")
	if err != nil {
		msg.Failf("Error: %s", err)
		os.Exit(1)
	} else {
		msg.Good("Done!")
	}
}
