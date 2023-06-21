package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	myCfg, err := NewConfig()
	if err != nil {
		t.Error(err)
	}
	fmt.Print(myCfg.Version)
}
