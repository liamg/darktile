package platform

import (
	"os/exec"
)

type cmdProc struct {
	cmd *exec.Cmd
}

func newCmdProc(c *exec.Cmd) *cmdProc {
	return &cmdProc{
		cmd: c,
	}
}

func (p *cmdProc) Wait() error {
	return p.cmd.Wait()
}

func (p *cmdProc) Close() error {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	ret := p.cmd.Process.Kill()
	p.cmd = nil
	return ret
}
