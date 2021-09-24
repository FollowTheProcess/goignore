package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetIgnore(t *testing.T) {
	want := []byte("fake gitignore stuff here")

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", string(want))
	}))
	defer fakeServer.Close()

	testURL := fakeServer.URL

	got, err := GetIgnore([]string{"fake", "targets"}, testURL)
	if err != nil {
		t.Fatalf("did not expect an error but got one: %s\n", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q expected %q", got, want)
	}
}

func TestGetList(t *testing.T) {
	listData := []byte("macos,vscode,python,go,bash")

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", string(listData))
	}))
	defer fakeServer.Close()

	testURL := fakeServer.URL

	t.Run("test fetches correct data", func(t *testing.T) {
		got, err := GetList(testURL)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s\n", err)
		}

		if !reflect.DeepEqual(got, listData) {
			t.Errorf("got %q expected %q", got, listData)
		}
	})

	t.Run("test prints correctly", func(t *testing.T) {
		want := fmt.Sprintf("%s\n", listData)

		got := bytes.Buffer{}

		printList(&got, testURL)

		assertString(t, got, want)
	})
}

func TestWriteToIgnoreFile(t *testing.T) {
	fakeIgnorePath := filepath.Join("testdata", ".fakeignore")
	nonExistIgnorePath := filepath.Join("testdata", ".nothereignore")
	fakeWriteData := []byte("this is some fake stuff")

	t.Run("should return error if file exists", func(t *testing.T) {
		want := errIgnoreFileExists

		err := WriteToIgnoreFile(fakeWriteData, fakeIgnorePath)

		if err == nil {
			t.Errorf("should have raised an error but didnt")
		} else if err != errIgnoreFileExists {
			t.Errorf("expected %q, got %q", want, err)
		}
	})

	t.Run("should write data to file if doesn't exist", func(t *testing.T) {
		want := fakeWriteData

		err := WriteToIgnoreFile(fakeWriteData, nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}
		// Remove the file again
		defer cleanUpFile(t, nonExistIgnorePath)

		// Now read from the newly created file
		got, err := os.ReadFile(nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, expected %q", got, want)
		}
	})
}

func TestPrintVersion(t *testing.T) {
	want := fmt.Sprintf("goignore version: %s\n", version)

	got := bytes.Buffer{}
	printVersion(&got)

	assertString(t, got, want)
}

func TestPrintUsage(t *testing.T) {
	want := fmt.Sprintf("%s\n", helpMessage)

	got := bytes.Buffer{}
	printUsage(&got)

	assertString(t, got, want)
}

func assertString(t testing.TB, got bytes.Buffer, want string) {
	t.Helper()

	if got.String() != want {
		t.Errorf("got %s, wanted %s", got.String(), want)
	}
}

func cleanUpFile(t testing.TB, filepath string) {
	t.Helper()
	err := os.Remove(filepath)
	if err != nil {
		t.Fatalf("did not expect an error but got one: %s", err)
	}
}
