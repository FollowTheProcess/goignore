package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetIgnoreData(t *testing.T) {
	want := []byte("fake gitignore stuff here")

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", string(want))
	}))
	defer fakeServer.Close()

	testURL := fakeServer.URL

	got, err := getIgnoreData(testURL, []string{"some", "targets"})
	if err != nil {
		t.Fatalf("getIgnoreData returned an error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
