package util

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	err := os.Setenv("HOME", "/tmp/gogobox")
	if err != nil {
		t.Fatal(err)
	}

	Init()
	defer os.RemoveAll("/tmp/gogobox")

	// parse config
	conf, err := ParseConfig(GlobalConfig.file)
	if err != nil {
		t.Fatal(err)
	}
	if conf.Version != GlobalConfig.Version {
		t.Fatalf("expected version %s, got %s", GlobalConfig.Version, conf.Version)
	}
	if conf.DictionaryConfig.Current != GlobalConfig.DictionaryConfig.Current {
		t.Fatalf("expected dictionary current %s, got %s", GlobalConfig.DictionaryConfig.Current, conf.DictionaryConfig.Current)
	}
	if conf.NotebookConfig.Path != GlobalConfig.NotebookConfig.Path {
		t.Fatalf("expected notebook path %s, got %s", GlobalConfig.NotebookConfig.Path, conf.NotebookConfig.Path)
	}
	// update config
	GlobalConfig.Version = "1.0.1"
	err = GlobalConfig.Sync()
	if err != nil {
		t.Fatal(err)
	}
	conf, err = ParseConfig(GlobalConfig.file)
	if err != nil {
		t.Fatal(err)
	}
	if conf.Version != "1.0.1" {
		t.Fatalf("expected version 1.0.1, got %s", conf.Version)
	}

}
