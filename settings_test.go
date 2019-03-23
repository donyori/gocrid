package gocrid

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Fatal("Cannot get GOPATH")
	}
	t.Log("GOPATH:", gopath)
	testFilename := filepath.Join(gopath, "src", "github.com", "donyori",
		"gocrid", "test_settings_resource", "settings_sample.json")
	t.Log("testFilename:", testFilename)
	settings, err := LoadSettings(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", settings)
}
