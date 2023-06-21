package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func configDir() string {
	var path string
	if a := os.Getenv("GOGOBOX_HOME"); a != "" {
		path = a
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "gogobox")
	}
	return path
}

func NewConfig() (*Config, error) {
	filePath := filepath.Join(configDir(), "config.yaml")
	root := defaultConfig
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// overwrite default config
	err = yaml.Unmarshal(bytes, &root)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return &root, nil
}

func (c *Config) ToString() (string, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type Config struct {
	Version string      `yaml:"version"`
	Dict    *DictConfig `yaml:"dict"`
}

type DictConfig struct {
	NotebookPath string `yaml:"notebook"`
}

var defaultConfig = Config{
	Version: "0.1",
	Dict: &DictConfig{
		NotebookPath: filepath.Join(configDir(), "notebook.json"),
	},
}
