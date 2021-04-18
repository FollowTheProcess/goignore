package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ignoreURL is the base url for the gitignore API
	ignoreURL      string = "https://www.toptal.com/developers/gitignore/api"
	versionMessage string = "goignore version: 0.2.1\n"
	listMessage    string = "To get a list of valid targets, run goignore --list"

	helpMessage string = `
Usage: goignore [OPTIONS] [ARGS]...

Handy CLI to generate great gitignore files.

Options:
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
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return data, nil

}

// GetList returns the gitignore API response for 'list'
func GetList(url string) ([]byte, error) {

	targetURL := strings.Join([]string{url, "list"}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return data, nil
}

// WriteToIgnoreFile takes data in and writes it to cwd/.gitignore
func WriteToIgnoreFile(data []byte, filename string) error {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ignoreFilePath := filepath.Join(cwd, filename)

	if _, err := os.Stat(ignoreFilePath); os.IsNotExist(err) {
		// No .gitignore, we're good to go
		_, err := os.Create(ignoreFilePath)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(ignoreFilePath, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			return err
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
	fmt.Fprintln(where, versionMessage)
}

func printList(where io.Writer, url string) {
	data, err := GetList(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(where, string(data))
}

func makeIgnoreFile(targets []string, url string) {
	data, err := GetIgnore(targets, url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// By default doesn't ignore all of .vscode/
	// personal preference
	vscode := []byte("\n.vscode/\n")
	data = append(data, vscode...)

	err = WriteToIgnoreFile(data, ".gitignore")
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Done!")
	}
}
