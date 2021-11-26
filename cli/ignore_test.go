package cli

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestWriteToIgnoreFile(t *testing.T) {
	t.Run("write if not exists", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		want := []byte("fake gitignore data here")

		app := New(stdout, stderr, afero.NewMemMapFs())

		if err := app.writeToIgnoreFile(".", want); err != nil {
			t.Errorf("writeToIgnoreFile returned an error: %v", err)
		}

		// Read the data back
		data, err := app.fs.ReadFile("./.gitignore")
		if err != nil {
			t.Errorf("could not read file: %v", err)
		}

		if !reflect.DeepEqual(data, want) {
			t.Errorf("written data mismatch: got %v, wanted %v", data, want)
		}
	})

	t.Run("error if exists", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		want := []byte("fake gitignore data here")

		app := New(stdout, stderr, afero.NewMemMapFs())

		// Create the file to trigger the "already exists" error
		err := app.fs.WriteFile("./.gitignore", []byte("I'm already here"), 0o755)
		if err != nil {
			t.Fatalf("error writing test file: %v", err)
		}

		if err := app.writeToIgnoreFile(".", want); err == nil {
			t.Errorf("expected an error but did not get one")
		}
	})
}
