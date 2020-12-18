package exectest

import (
	"github.com/betasve/mstd/exec"
	e "os/exec"
)

type CmdMock struct{}

type CmdRunMock struct {
	cmd *e.Cmd
}

var CommandMockFunc = func(name string, arg ...string) exec.RunService {
	return &CmdRunMock{cmd: e.Command(name, arg...)}
}
var CommandRunMockFunc = func(c *CmdRunMock) error { return c.cmd.Run() }

func InitCmdMock(command *e.Cmd) *CmdRunMock {
	return &CmdRunMock{cmd: command}
}

func (c *CmdMock) Command(name string, arg ...string) exec.RunService {
	return CommandMockFunc(name, arg...)
}

func (c *CmdRunMock) Run() error {
	return CommandRunMockFunc(c)
}
