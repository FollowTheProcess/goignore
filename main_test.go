package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

	want := []byte("macos,vscode,python,go,bash")

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", string(want))
	}))
	defer fakeServer.Close()

	testURL := fakeServer.URL

	got, err := GetList(testURL)
	if err != nil {
		t.Fatalf("did not expect an error but got one: %s\n", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q expected %q", got, want)
	}

}

func TestWriteToIgnoreFile(t *testing.T) {

	fakeIgnorePath := "./testdata/.fakeignore"
	nonExistIgnorePath := "./testdata/.nothereignore"
	fakeWriteData := []byte("this is some fake stuff")

	t.Run("should return error if file exists", func(t *testing.T) {

		want := ErrIgnoreFileExists

		err := WriteToIgnoreFile(fakeWriteData, fakeIgnorePath)

		if err == nil {
			t.Errorf("should have raised an error but didnt")
		} else if err != ErrIgnoreFileExists {
			t.Errorf("expected %q, got %q", want, err)
		}
	})

	t.Run("should write data to file if doesn't exist", func(t *testing.T) {

		want := fakeWriteData

		err := WriteToIgnoreFile(fakeWriteData, nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}

		// Now read from the newly created file
		file, err := os.Open(nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}
		defer file.Close()

		got, err := os.ReadFile(nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, expected %q", got, want)
		}

		// Remove the file again
		err = os.Remove(nonExistIgnorePath)
		if err != nil {
			t.Fatalf("did not expect an error but got one: %s", err)
		}

	})
}
