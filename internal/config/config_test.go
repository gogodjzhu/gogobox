package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	myCfg, err := NewConfig("test")
	if err != nil {
		t.Error(err)
	}
	myCfg.Set("test", "test")
}
