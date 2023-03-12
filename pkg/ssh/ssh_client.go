package ssh

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
)

type Client struct {
	cli *ssh.Client
}

func NewClient(knownHostsPath string, addr, username, password string) (*Client, error) {
	hostkeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: hostkeyCallback,
	}
	cli, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}

func (c *Client) Close() {
	_ = c.cli.Close()
}

func (c *Client) Run(cmd string, stdOut, errOut io.Writer) error {
	session, err := c.cli.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = stdOut
	session.Stderr = errOut
	return session.Run(cmd)
}
