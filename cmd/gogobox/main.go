package main

import (
	"fmt"
	"gogobox/pkg/cmd/root"
	"gogobox/pkg/cmdutil"
	"os"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmdFactory := cmdutil.NewFactory()

	mainCmd, err := root.NewCmdRoot(cmdFactory)
	if err != nil {
		fmt.Println(err)
		return exitError
	}
	if _, err := mainCmd.ExecuteC(); err != nil {
		return exitError
	}
	return exitOK
}
