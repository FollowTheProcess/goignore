package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const URL string = "https://www.toptal.com/developers/gitignore/api"

var ErrIgnoreFileExists = errors.New(".gitignore already exists, not doing anything")

func main() {

	ignoreList := os.Args[1:]

	data, err := GetIgnore(ignoreList)
	if err != nil {
		fmt.Printf("Error occured %s\n", err)
		os.Exit(1)
	}

	// By default doesn't ignore all of .vscode/
	// Personal preference
	vscode := []byte("\n.vscode/\n")
	data = append(data, vscode...)

	err = WriteToIgnoreFile(data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Done!")
	}

}

// GetIgnore hits the gitignore.io API with targets and returns the response
func GetIgnore(targets []string) ([]byte, error) {

	constructedString := strings.Join(targets, ",")

	targetURL := strings.Join([]string{URL, constructedString}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil

}

// WriteToIgnoreFile takes data in and writes it to cwd/.gitignore
func WriteToIgnoreFile(data []byte) error {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ignoreFilePath := filepath.Join(cwd, ".gitignore")

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
