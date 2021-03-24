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

// URL is the base url for the gitignore API
const URL string = "https://www.toptal.com/developers/gitignore/api"

// ErrIgnoreFileExists when there is already a gitignore file
var ErrIgnoreFileExists = errors.New(".gitignore already exists, not doing anything")

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Expected either list of valid gitignore.io targets or 'list'.")
		os.Exit(1)
	}

	if os.Args[1] == "list" {
		// Return the list of valid targets
		fmt.Println("Valid gitignore targets...")
		data, err := GetList(URL)
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		// Only other usage is for valid gitignore targets
		ignoreList := os.Args[1:]
		fmt.Printf("Writing gitignore for %v...\n", ignoreList)

		data, err := GetIgnore(ignoreList, URL)
		if err != nil {
			fmt.Printf("Error occurred: %s\n", err)
			os.Exit(1)
		}

		// By default doesn't ignore all of .vscode/
		// personal preference
		vscode := []byte("\n.vscode/\n")
		data = append(data, vscode...)

		err = WriteToIgnoreFile(data, ".gitignore")
		if err != nil {
			fmt.Printf("Error occurredL %s\n", err)
			os.Exit(1)
		} else {
			fmt.Println("Done!")
		}
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
