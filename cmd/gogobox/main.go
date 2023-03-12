package main

import (
	"gogobox/pkg/cmd/root"
	"gogobox/pkg/cmdutil"
	"os"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmdFactory := cmdutil.NewFactory()

	mainCmd := root.NewCmdRoot(cmdFactory)
	if cmd, err := mainCmd.ExecuteC(); err != nil {
		if cmd != nil {
			cmd.PrintErrln(err)
		}
		return exitError
	}
	return exitOK
}
