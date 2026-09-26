// Example of what your LocalSystem process does with the sysinput (B1) backend.
//
// This process runs as LocalSystem (Session 0). It does NOT call SendInput
// itself — it drives a Manager, which keeps one agent alive in the interactive
// session and forwards input there (including the secure desktop).
//
// Your real process would replace the demo loop at the bottom with your
// remote-desktop transport: decode input events from the peer and call
// mgr.MoveTo / mgr.Click / mgr.KeyTap / mgr.Type.
//
//	GOOS=windows GOARCH=amd64 go build -o loki-controller.exe ./examples/localsystem
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hkloudou/loki/sysinput"
)

func main() {
	// The agent exe is normally deployed next to this controller.
	exeDir, _ := os.Executable()
	agentPath := filepath.Join(filepath.Dir(exeDir), "loki-inputagent.exe")

	mgr := sysinput.NewManager(agentPath)
	mgr.SetLogger(func(s string) { log.Println("[sysinput]", s) })
	mgr.Start()      // launches the agent and keeps it alive across session changes
	defer mgr.Stop() // on service shutdown

	// --- your transport loop goes here; demo below just waits for a session ---
	for !mgr.Available() {
		time.Sleep(200 * time.Millisecond)
	}

	// Example: click at (500,400) and type some text on the user's desktop.
	if err := mgr.MoveTo(500, 400); err != nil {
		log.Println("moveto:", err) // ErrNoSession if the user logged off
	}
	_ = mgr.Click(sysinput.Left)
	_ = mgr.Type("你好, loki")
	_ = mgr.KeyTap(sysinput.VKReturn)

	select {} // a real service blocks until its stop signal
}
