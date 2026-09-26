# loki

A Windows **virtual HID (Human Interface Device) mouse + keyboard** that can be
driven from Go. It is the input-injection core for building a remote-desktop /
RPA system (think ToDesk / 影刀 / UiPath): a kernel driver presents a fake HID
mouse and keyboard to Windows, and writing HID output reports to it produces
input events that the OS — and most games / anti-cheat / anti-bot layers — treat
as coming from real hardware.

This is fundamentally more robust than `SendInput` / `mouse_event` /
`keybd_event`, which are user-mode APIs that many protected apps can detect
(`LLMHF_INJECTED`) or that fail across UAC / secure-desktop boundaries.

## Layout

```
loki/
├── loki.go, loki_windows.go, loki_other.go   Go library (public API)
├── keys.go                                    HID key codes (generated)
├── gen_keys.sh                                regenerates keys.go from keys.h
├── examples/main.go                           usage example
└── windows/
    ├── driver/   KMDF kernel-mode virtual HID driver  -> hidriver.sys
    └── lib/      user-mode C++ lib + C-ABI wrapper     -> loki.dll
```

There are **two independent input backends**, chosen per deployment:

- **Driver backend** (root `loki` package + `loki.dll`) — a signed virtual HID
  driver; input looks like real hardware, so games / anti-cheat tend to accept
  it. Requires a signed kernel driver.
- **`sysinput` backend** ([`sysinput/`](sysinput), [docs/sysinput.md](docs/sysinput.md))
  — driverless: a LocalSystem service spawns an in-session agent that injects via
  `SendInput`, covering the physical console **and** the secure desktop (UAC /
  logon). No driver signing; but injected input is detectable by anti-cheat.

The driver stack, bottom to top:

1. **`windows/driver`** — a KMDF kernel driver that registers a virtual HID
   device exposing a mouse collection and a keyboard collection. User mode sends
   it HID *output* reports; the driver turns them into HID *input* reports that
   Windows delivers as genuine input. Produces `hidriver.sys` + `hidriver.inf`.
2. **`windows/lib`** — a user-mode C++ library (`Device` / `Mouse` / `Keyboard`
   / `RegistryService`) that opens the driver's HID interface and writes reports,
   plus **`loki_c_api.*`**, a flat `extern "C"` wrapper compiled into
   **`loki.dll`** so any language with a C FFI can use it.
3. **root Go package** — loads `loki.dll` with `syscall.NewLazyDLL` (no cgo, no
   third-party dependencies) and exposes idiomatic `Mouse` / `Keyboard` types.

## Go usage

```go
import "github.com/hkloudou/loki"

m, _ := loki.NewMouse()
defer m.Close()
m.Move(400, 300)      // absolute move (walks the cursor)
m.LeftClick()
m.MoveRelative(10, 0) // one relative report, deltas clamped to [-127,127]

k, _ := loki.NewKeyboard()
defer k.Close()
k.Type(loki.KeyA)                    // press+release 'a'
k.Send(loki.ModLshift, loki.KeyH)    // 'H' (shift held)
k.Send(0)                            // release all keys
```

At runtime `loki.dll` must be loadable: next to your `.exe`, on the DLL search
path, or via the `LOKI_DLL` environment variable (absolute path). The virtual
driver must be installed first (below). On non-Windows platforms the package
still compiles; every constructor returns `loki.ErrUnsupported`.

## Building

### Driver (`hidriver.sys`)
Requires Visual Studio + the Windows Driver Kit (WDK). Open
`windows/driver/loki-driver.vcxproj` (or add it to a solution) and build
`Release|x64`. Driver development and signing details are in
[`windows/driver/README.md`](windows/driver/README.md).

### `loki.dll`
From an *x64 Native Tools Command Prompt for VS*:

```bat
cd windows\lib
build.bat Release
```

This produces `build\x64\Release\loki.dll`. Or build
`windows/lib/loki-lib.vcxproj` in Visual Studio.

## Installing the driver (test mode)

The driver currently uses a **test certificate**, so it only loads with driver
signature enforcement relaxed. In an **elevated** prompt:

```bat
bcdedit /set testsigning on
:: reboot, then:
devcon install hidriver.inf root\hidriver
```

Turn it off again with `bcdedit /set testsigning off`. Getting this driver to
load on normal, non-test machines requires proper signing — see
[docs/signing.md](docs/signing.md) for the real options (EV + Microsoft
attestation signing, and why that matters for a ToDesk-style product).

## Relationship to prior work / licenses

Built on two upstream projects; see [docs/comparison.md](docs/comparison.md) for
a file-by-file diff.

- [hedgar2017/loki-hidriver](https://github.com/hedgar2017/loki-hidriver) — MIT,
  © 2018 Alex Zarudnyy (original driver + example lib).
- [dengqizhou30/HIDDriver](https://github.com/dengqizhou30/HIDDriver) — Apache
  2.0 (Windows 10 compatibility fork, static-lib split, extra mouse moves).

The driver C sources are unchanged from upstream; loki adds the C-ABI DLL
wrapper, the Go bindings, and the build tooling.
