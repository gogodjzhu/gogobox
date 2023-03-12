package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

type Config interface {
	Get(key string) (interface{}, error)
	GetOrDefault(key string, value string) (string, error)
	Set(key string, value string)
	Write() error
}

func NewConfig(filePath string) (Config, error) {
	root := make(map[string]interface{})
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &root)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(root) == 0 {
		err := yaml.Unmarshal([]byte(defaultGeneralEntries), &root)
		if err != nil {
			return nil, errors.New("invalid default config")
		}
	}
	return &cfg{filePath, root, sync.RWMutex{}}, nil
}

// Implements Config interface
type cfg struct {
	filePath string
	root     map[string]interface{}
	lock     sync.RWMutex
}

func (c *cfg) Get(key string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.root == nil {
		return nil, errors.New("config is empty")
	}
	return c.root[key], nil
}

func (c *cfg) GetOrDefault(s string, s2 string) (string, error) {
	panic("implement me")
}

func (c *cfg) Set(s2 string, s3 string) {
	panic("implement me")
}

func (c *cfg) Write() error {
	panic("implement me")
}

var defaultGeneralEntries = `
# instance name
node_name: gogobox-node1
# sub-config
hosts:
	  # host name
	  - name: gogobox-host1
`
