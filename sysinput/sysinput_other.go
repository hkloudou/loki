//go:build !windows

package sysinput

import "io"

// AgentProcess is unavailable on non-Windows platforms.
type AgentProcess struct{}

// ActiveConsoleSessionID returns 0xFFFFFFFF on non-Windows platforms.
func ActiveConsoleSessionID() uint32 { return 0xFFFFFFFF }

// LaunchAgentInActiveSession returns ErrUnsupported on non-Windows platforms.
func LaunchAgentInActiveSession(agentPath string, extraArgs ...string) (*AgentProcess, error) {
	return nil, ErrUnsupported
}

// RunAgent returns ErrUnsupported on non-Windows platforms.
func RunAgent(r io.Reader) error { return ErrUnsupported }

func (p *AgentProcess) PID() uint32                   { return 0 }
func (p *AgentProcess) MoveTo(x, y int) error         { return ErrUnsupported }
func (p *AgentProcess) Move(dx, dy int) error         { return ErrUnsupported }
func (p *AgentProcess) ButtonDown(b Button) error     { return ErrUnsupported }
func (p *AgentProcess) ButtonUp(b Button) error       { return ErrUnsupported }
func (p *AgentProcess) Click(b Button) error          { return ErrUnsupported }
func (p *AgentProcess) Wheel(notches int) error       { return ErrUnsupported }
func (p *AgentProcess) KeyDown(vk, scan uint16) error { return ErrUnsupported }
func (p *AgentProcess) KeyUp(vk, scan uint16) error   { return ErrUnsupported }
func (p *AgentProcess) KeyTap(vk uint16) error        { return ErrUnsupported }
func (p *AgentProcess) Type(s string) error           { return ErrUnsupported }
func (p *AgentProcess) Close() error                  { return nil }
