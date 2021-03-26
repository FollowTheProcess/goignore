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

// URL is the base url for the gitignore API
const URL string = "https://www.toptal.com/developers/gitignore/api"
const versionMessage string = "goignore version: 0.1.0"

const invalidArgsMessage string = `Error: Expected list of gitignore.io targets or a CLI option.
run '$ goignore --help' for help.`

const helpMessage string = `
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

// ErrIgnoreFileExists when there is already a gitignore file
var ErrIgnoreFileExists = errors.New("'.gitignore' already exists. Not doing anything")

func main() {

	helpFlag := flag.Bool("help", false, "--help")
	versionFlag := flag.Bool("version", false, "--version")
	listFlag := flag.Bool("list", false, "--list")

	flag.Parse()

	if *helpFlag {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println(versionMessage)
		os.Exit(0)
	}

	if *listFlag {
		// Return the list of valid targets
		data, err := GetList(URL)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("\n", string(data))
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println(invalidArgsMessage)
		os.Exit(1)
	}

	// Only other usage is for valid gitignore targets
	ignoreList := os.Args[1:]

	data, err := GetIgnore(ignoreList, URL)
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
		return ErrIgnoreFileExists
	}

	return nil
}
