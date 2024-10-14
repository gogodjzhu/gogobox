package util

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	gopkgyaml "gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	file             string
	Version          string            `mapstructure:"version" yaml:"version"`
	DictionaryConfig *DictionaryConfig `mapstructure:"dictionary" yaml:"dictionary"`
	NotebookConfig   *NotebookConfig   `mapstructure:"notebook" yaml:"notebook"`
}

func (c *Config) Sync() error {
	bs, err := gopkgyaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(GlobalConfig.file, bs, 0755)
}

type DictionaryConfig struct {
	Current string `mapstructure:"current" yaml:"current"`
}

type NotebookConfig struct {
	Path string `mapstructure:"path" yaml:"path"`
}

var GlobalConfig *Config

func Init() {
	config.AddDriver(yaml.Driver)

	var configHome string
	if os.IsPathSeparator('\\') {
		configHome = os.Getenv("APPDATA")
	} else {
		configHome = os.Getenv("HOME")
	}
	configHome = configHome + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "gogobox"
	if _, err := os.Stat(configHome); os.IsNotExist(err) {
		_ = os.MkdirAll(configHome, os.ModePerm)
	}
	configFile := configHome + string(os.PathSeparator) + "config.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		GlobalConfig = defaultConfig(configHome)
		GlobalConfig.file = configFile
		if err := GlobalConfig.Sync(); err != nil {
			panic(err)
		}
	} else {
		GlobalConfig, err = ParseConfig(configFile)
		if err != nil {
			panic(err)
		}
		GlobalConfig.file = configFile
	}
}

func defaultConfig(home string) *Config {
	var conf = Config{
		Version: "1.0.0",
		DictionaryConfig: &DictionaryConfig{
			Current: "youdao",
		},
		NotebookConfig: &NotebookConfig{
			Path: home + string(os.PathSeparator) + "notebook",
		},
	}
	return &conf
}

func ParseConfig(file string) (*Config, error) {
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(file)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = config.Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
