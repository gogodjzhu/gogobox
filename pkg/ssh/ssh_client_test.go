package ssh

import (
	"os"
	"testing"
)

func TestClient_Run(t *testing.T) {
	cli, err := NewClient("/Users/djzhu/.ssh/known_hosts", "gogodjzhu.com:2212", "chak", "djzhu1984")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	err = cli.Run("top", os.Stdout, os.Stderr)
}
