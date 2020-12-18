package exec

import (
	e "os/exec"
)

type CmdService interface {
	Command(name string, arg ...string) RunService
}

type RunService interface {
	Run() error
}

type Cmd struct{}

type CmdRun struct {
	cmd *e.Cmd
}

var CmdClient CmdService = &Cmd{}
var RunClient RunService = &CmdRun{}

func (c *Cmd) Command(name string, arg ...string) RunService {
	return &CmdRun{cmd: e.Command(name, arg...)}
}

func (c *CmdRun) Run() error {
	return c.cmd.Run()
}
