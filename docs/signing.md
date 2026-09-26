# Driver signing & better alternatives to a test certificate

> Direct answer to: *"besides a test certificate, is there a better way?"* — and
> how ToDesk / 影刀 (YingDao) / UiPath actually inject input.

## Why the test certificate is not an option for a product

A self-made test certificate only works when the machine is in **test signing
mode** (`bcdedit /set testsigning on` + reboot). That means:

- Every end user must disable a core security feature and reboot. Non-starter
  for consumer software.
- A **"Test Mode" watermark** sits on the desktop.
- It lowers the machine's security posture and is often flagged by AV/EDR and
  blocked by enterprise policy.
- On Secure-Boot machines it interacts badly with policy; many managed devices
  simply refuse.

So for anything shipped, you need a driver that loads on a **stock** Windows 10
1607+ / Windows 11 machine. There are two broad strategies: **sign the kernel
driver properly**, or **don't ship a kernel driver at all**.

---

## Option A — Sign the driver the way Windows now requires (recommended for a product)

Since Windows 10 version 1607, new kernel-mode drivers are no longer loaded by
old-style cross-signing. The supported path is **Microsoft signing through
Partner Center (Windows Hardware Dev Center)**:

1. **Get an EV (Extended Validation) code-signing certificate.**
   - From a CA such as DigiCert / Sectigo / GlobalSign. Ships on a hardware
     token (or cloud HSM). Requires organization identity validation (a legal
     entity + something like a D-U-N-S number). Roughly a few hundred USD/year.
   - The EV cert is what lets you **register/validate a Partner Center hardware
     account** and sign your submission bundle.
2. **Choose the Microsoft signing tier:**
   - **Attestation signing** — you upload the driver package (a `.cab` you EV-sign)
     and Microsoft signs it. **No HLK hardware tests required.** This is the right
     tier for a pure *software* HID driver like loki. You select which Windows
     versions to target; the returned `.cat` chains to the Microsoft root and
     loads on those versions with **no test mode**.
   - **WHQL / HLK certification** — full Hardware Lab Kit test passes, gets the
     "Certified for Windows" logo and Windows Update distribution. Overkill for a
     virtual HID driver, but it's the same portal if you ever want it.
3. **Ship** the Microsoft-signed `.sys` + `.cat` + `.inf`. Installs on normal
   machines via `devcon`, `pnputil /add-driver ... /install`, or your installer.

This is what a commercial remote-desktop product does when it ships a driver.
Cost/effort: an EV cert, a Partner Center account, and the submission workflow —
but then it "just works" for every user.

> Note on the loki sources: `VENDOR_ID`/`PRODUCT_ID` are `0x00`. Before
> submitting for signing, set real IDs and unique hardware IDs / class strings so
> your device is distinct.

---

## Option B — Avoid the kernel driver entirely

The signing problem exists *because* you shipped a kernel driver. For many use
cases you don't need one.

### B1. SYSTEM service + `SendInput` (what remote-desktop tools mostly do)

For controlling **normal applications** (the ToDesk / AnyDesk / TeamViewer /
RustDesk case), input is usually injected in **user mode** with `SendInput`,
with the twist that a **Windows service running as `LocalSystem`** does the
injecting and attaches to the correct desktop:

- Enumerate/attach to the active input desktop with `OpenInputDesktop` /
  `SetThreadDesktop`, so you can drive the **secure desktop** (UAC elevation
  prompts, the logon screen, Ctrl+Alt+Del) that a normal user process can't touch.
- The service survives session switches and the Winlogon secure desktop.

RustDesk (open source) is a clean reference for exactly this model. No driver, no
signing. **Limitation:** injected input carries the `LLMHF_INJECTED` flag, so
some games and aggressive anti-cheat / anti-bot systems ignore or detect it.

### B2. External USB-HID microcontroller ("hardware injector" / KMBox-style)

Put the HID device *outside* the PC: a small microcontroller (RP2040 / Arduino
Leonardo / STM32, or a ready-made "KMBox"-type board) enumerates as a **real USB
keyboard + mouse**. Your host software sends it commands over USB-serial /
network, and it emits genuine HID reports.

- **Windows uses its in-box HID driver** → *zero* driver-signing problem.
- The input is real hardware input → effectively **undetectable**, works even
  against strict anti-cheat, and works before the OS/agent is up.
- **Cost:** you need the physical dongle on each controlled machine — fine for a
  gaming/anti-detection rig or a lab, impractical for mass consumer remote
  desktop.

### B3. Reuse the RDP terminal input bus (`TERMINPUT_BUS`)

Windows' own Remote Desktop stack creates virtual HID keyboard/mouse devices per
remote session, with instance IDs like
`TERMINPUT_BUS\UMB\...&Session1Keyboard0`. They are produced by Microsoft's
in-box, signed `rdpbus`/UMBus drivers — so the appeal is "already signed, already
present."

Two caveats decide how you can use them:

- **You cannot open that devnode and write reports to inject input.** It is an
  input *source* fed by the RDP wire protocol (the remote client's keystrokes
  arrive over the RDP virtual channel and the terminal stack emits them through
  this device). There is no public "write HID output report → input" sink like
  loki's own driver exposes.
- **It only exists inside an active terminal session** (`Session1`+), not on the
  physical console, and Home SKUs have no RDP host.

The usable form of this idea is the **loopback / headless RDP-session
technique**: start a real terminal session (mstsc to a loopback / second
session, or a headless RDP session) and inject with `SendInput` *inside that
session*. The OS then attributes the input to `TERMINPUT_BUS` — a legitimately
signed HID device — which many games / anti-bot layers accept as real hardware,
with **no driver signing at all**. Trade-offs: needs an RDP-capable SKU; the
automated work runs in a session separate from the console (use a virtual display
or the console locks); and some anti-cheat explicitly flags remote sessions.
This is a strong signing-free option when a headless/second session is
acceptable; the custom driver (Option A) remains the choice when you must drive
the **physical console** with no RDP.

### B4. Existing signed input drivers

If you want kernel-level injection without running your own signing pipeline, you
can build on a driver someone else already ships signed:

- **Interception** — a signed keyboard/mouse filter driver with a user-mode API;
  installs on stock machines. It filters/injects at the device stack above real
  devices. Good for input tooling; it's a filter model, not a fully virtual
  device, and adds its own install footprint.
- (**ViGEmBus** is signed and virtual, but it emulates **game controllers**
  (Xbox/DS4), not a keyboard/mouse, so it doesn't fit this project.)

---

## How 影刀 / UiPath / ToDesk map onto this

- **RPA tools (影刀 / UiPath / Power Automate Desktop):** primarily **do not ship
  kernel drivers**. They automate through **UI Automation (UIA) / MSAA**, Win32
  **window messages** (`PostMessage`/`SendMessage` to specific controls — works in
  the background, no real cursor movement), and **`SendInput`** ("hardware
  events" mode, which moves the real cursor). UiPath explicitly exposes these as
  selectable input methods: *SendWindowMessages* (background), *Simulate*
  (control API), and *Hardware Events* (`SendInput`). This covers business apps;
  it does not defeat games/anti-cheat, which those products don't target.

- **Remote desktop (ToDesk / AnyDesk / TeamViewer / RustDesk):** mostly the **B1
  model** — a privileged **SYSTEM service + `SendInput` + desktop switching** to
  handle UAC/secure desktop/login. Some products add an **optional signed kernel
  driver** (Option A) for the minority of apps/games that reject injected input.

**Recommendation for loki:**

- If the goal is a shippable remote-desktop / RPA core for **normal apps** →
  start with **B1** (SYSTEM service + `SendInput` + `SetThreadDesktop`); no
  signing needed, ships today.
- If you specifically need input that **games / anti-bot layers accept**: **A**
  (EV cert + Microsoft attestation-sign this driver) for a pure-software solution
  that drives the physical console; **B3** (loopback/headless RDP session +
  `SendInput`) when a separate session is acceptable and you want zero signing;
  or **B2** (external USB-HID board) for the strongest, hardware-genuine path.
- The **test certificate** is fine only for **your own dev/CI machines**.
