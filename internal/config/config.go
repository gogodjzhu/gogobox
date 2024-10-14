package config

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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

func InitConfig(configFilename string) error {
	_, err := os.Stat(configFilename)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("config file already exists: %s", configFilename))
	}
	// create config dir
	conf := defaultConfig
	conf.Common.ConfigFilename = configFilename
	return conf.Save()
}

func ReadConfig() (*Config, error) {
	return ReadConfigSpecified("")
}

func ReadConfigSpecified(configFilename string) (*Config, error) {
	var root Config
	if configFilename == "" {
		configFilename = defaultConfig.Common.ConfigFilename
	}
	bytes, err := os.ReadFile(configFilename)
	if err != nil {
		if os.IsNotExist(err) {
			err = InitConfig(configFilename)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed init config file: %s", configFilename))
			}
			// read again
			bytes, err = os.ReadFile(configFilename)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed read config file: %s", configFilename))
			}
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf("failed read config file: %s", configFilename))
		}
	}
	// overwrite default config
	err = yaml.Unmarshal(bytes, &root)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed unmarshal config file: %s", configFilename))
	}
	return &root, nil
}

func (c *Config) ToString() (string, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return "", errors.Wrap(err, "failed marshal config")
	}
	return string(bytes), nil
}

func (c *Config) Save() error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed marshal config")
	}
	return os.WriteFile(c.Common.ConfigFilename, bytes, 0644)
}

type Config struct {
	Version  string          `yaml:"version"`
	Common   *CommonConfig   `yaml:"common"`
	Dict     *DictConfig     `yaml:"dict"`
	Notebook *NotebookConfig `yaml:"notebook"`
	Server   *ServerConfig   `yaml:"server"`
}

type CommonConfig struct {
	HomeDir        string `yaml:"homeDir"`
	ConfigFilename string `yaml:"configFilename"`
}

type DictConfig struct {
	Endpoint        string               `yaml:"endpoint"`
	EcdictConfig    *DictEcdictConfig    `yaml:"ecdictConfig"`
	YoudaoConfig    *DictYoudaoConfig    `yaml:"youdaoConfig"`
	EtymonineConfig *DictEtymonineConfig `yaml:"etymonineConfig"`
	ChatgptConfig   *DictChatgptConfig   `yaml:"chatgptConfig"`
	MWebsterConfig  *DictMWebsterConfig  `yaml:"mwebsterConfig"`
}

type DictEcdictConfig struct {
	DBFilename string `yaml:"dbFilename"`
}

type DictYoudaoConfig struct {
}

type DictEtymonineConfig struct {
}

type DictMWebsterConfig struct {
	Key string `yaml:"key"`
}

type DictChatgptConfig struct {
	ResourceName string `yaml:"resourceName"`
	DeploymentId string `yaml:"deploymentId"`
	ApiVersion   string `yaml:"apiVersion"`
	Key          string `yaml:"key"`
}

type NotebookConfig struct {
	CurrentChapter     string                      `yaml:"currentChapter"`
	FileNotebookConfig *NotebookFileNotebookConfig `yaml:"fileNotebookConfig"`
}

type NotebookFileNotebookConfig struct {
	Directory string `yaml:"directory"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Root string `yaml:"root"`
}

var defaultConfig = Config{
	Version: "0.1",
	Common: &CommonConfig{
		HomeDir:        configDir(),
		ConfigFilename: filepath.Join(configDir(), "config.yaml"),
	},
	Dict: &DictConfig{
		Endpoint: "youdao",
		EcdictConfig: &DictEcdictConfig{
			DBFilename: filepath.Join(configDir(), "stardict.db"),
		},
		YoudaoConfig:    &DictYoudaoConfig{},
		EtymonineConfig: &DictEtymonineConfig{},
		MWebsterConfig:  &DictMWebsterConfig{},
	},
	Notebook: &NotebookConfig{
		CurrentChapter: "default",
		FileNotebookConfig: &NotebookFileNotebookConfig{
			Directory: filepath.Join(configDir(), "notebooks"),
		},
	},
	Server: &ServerConfig{
		Port: 8080,
		Root: "/",
	},
}
