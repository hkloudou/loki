# B1 backend: `sysinput` (LocalSystem service + SendInput)

The `github.com/hkloudou/loki/sysinput` package is the **driverless** input
backend. It is intentionally **separate from the driver backend** (root `loki`
package + `loki.dll`): different package, different mechanism, no shared runtime
requirement. Pick per deployment.

| | Driver backend (`loki`) | sysinput backend (`sysinput`) |
|---|---|---|
| Mechanism | virtual HID driver + `loki.dll` | `SendInput` from an in-session agent |
| Needs kernel driver | yes (signing) | **no** |
| Physical console | yes | yes |
| Secure desktop (UAC/logon) | yes | yes (agent runs as SYSTEM-in-session) |
| Games / anti-cheat | often accepted (real HID) | often rejected (`LLMHF_INJECTED`) |
| Ship today, no signing | no | **yes** |

## Why an agent process is needed

A service running as **LocalSystem lives in Session 0**, which is isolated from
the interactive user session. `SendInput` called from Session 0 does not reach
the user's desktop, and you cannot `SetThreadDesktop` across the session
boundary. So the service must place a small **agent** *inside* the user session
and let that agent inject.

To also cover the **secure desktop** (UAC elevation, the logon screen,
Ctrl+Alt+Del), the agent is launched **as SYSTEM within the session** (the
service duplicates its own SYSTEM token and retargets it with
`SetTokenInformation(TokenSessionId)`). A user-token agent cannot access the
`Winlogon` desktop; a SYSTEM-in-session agent can, and it *follows* the active
input desktop with `OpenInputDesktop` + `SetThreadDesktop` before every event.

## What your LocalSystem process actually does

Your process runs in Session 0 and must **not** call `SendInput` itself. Its job:

1. Launch the agent into the active session (as SYSTEM-in-session).
2. Relaunch it when the session changes (logon / unlock / RDP / fast-user-switch)
   or when the agent dies.
3. Forward the input events it receives from the remote peer to the current agent.
4. Close the agent on shutdown.

Steps 1, 2 and 4 are handled for you by **`Manager`** — this is the recommended
entry point. Link the package, construct a `Manager`, `Start()` it, then just call
input methods from your transport loop:

```go
import "github.com/hkloudou/loki/sysinput"

mgr := sysinput.NewManager(`C:\Program Files\loki\loki-inputagent.exe`)
mgr.SetLogger(func(s string){ log.Println("[sysinput]", s) })
mgr.Start()
defer mgr.Stop()

// in your remote-input handler:
mgr.MoveTo(400, 300)      // absolute, virtual-screen pixels
mgr.Click(sysinput.Left)
mgr.KeyTap(sysinput.VKReturn)
mgr.Type("你好, world")    // Unicode, layout-independent
mgr.Wheel(-3)             // scroll down 3 notches
// methods return sysinput.ErrNoSession while no user is logged on.
```

`Manager` polls `WTSGetActiveConsoleSessionId` (default every 1s) and checks agent
liveness, relaunching as needed — so you don't have to wire an SCM
`SERVICE_CONTROL_SESSIONCHANGE` handler. If you *are* a real SCM service and want
instant reaction, you can additionally call `LaunchAgentInActiveSession` yourself
from a session-change handler; the low-level API is exported for that. See
[examples/localsystem](../examples/localsystem).

Under the hood the process talks to the agent over the agent's **stdin pipe**
(newline-delimited JSON) — no named pipe, no extra ACLs. Commands are applied in
order on the active desktop.

Build the agent binary:

```bash
GOOS=windows GOARCH=amd64 go build -o loki-inputagent.exe ./cmd/loki-inputagent
```

## Session lifecycle (important)

`LaunchAgentInActiveSession` targets the **current** active console session. In a
real service, relaunch the agent on session changes so it stays valid:

- Handle `SERVICE_CONTROL_SESSIONCHANGE` (or `WTSRegisterSessionNotification`):
  on `WTS_CONSOLE_CONNECT` / `WTS_SESSION_LOGON` / `WTS_SESSION_UNLOCK`, call
  `LaunchAgentInActiveSession` again and `Close()` the old one.
- On fast-user-switching / RDP, the active session id changes; re-launch.
- If `ActiveConsoleSessionID()` returns `0xFFFFFFFF`, no one is logged on.

## Known limitations / notes

- **Anti-cheat:** injected input carries `LLMHF_INJECTED`; protected games may
  ignore it. Use the driver backend there.
- **Environment block:** the agent is launched with the service's environment
  (nil env). If you need the user's environment, extend the launcher to call
  `CreateEnvironmentBlock`.
- **UIPI/integrity:** running the agent as SYSTEM (high integrity) is what lets it
  drive elevated windows and the secure desktop.
- Not safe to inject from multiple services at once into one session; keep one
  agent per session.
