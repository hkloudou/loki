// Command loki-inputagent is the in-session input agent for the sysinput (B1)
// backend. It is launched by the LocalSystem service via
// sysinput.LaunchAgentInActiveSession, runs inside the interactive session, and
// injects the commands it reads from stdin with SendInput on the active desktop
// (following the user onto the secure desktop).
//
// It is not meant to be run by hand; the service spawns it with its stdin wired
// to the command pipe. Build:
//
//	GOOS=windows go build -o loki-inputagent.exe ./cmd/loki-inputagent
package main

import (
	"os"

	"github.com/hkloudou/loki/sysinput"
)

func main() {
	// Blocks until stdin (the service's command pipe) reaches EOF.
	if err := sysinput.RunAgent(os.Stdin); err != nil {
		os.Exit(1)
	}
}
