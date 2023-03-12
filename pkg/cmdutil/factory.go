package cmdutil

import (
	"io"
	"os"
)

type Factory struct {
	IOStreams *IOStreams
}

func NewFactory() *Factory {
	f := &Factory{
		IOStreams: ioStreams(),
	}
	return f
}

type IOStreams struct {
	In  io.Reader
	Out io.Writer
}

func ioStreams() *IOStreams {
	return &IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}
