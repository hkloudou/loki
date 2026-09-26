# Upstream comparison: `hedgar2017/loki-hidriver` vs `dengqizhou30/HIDDriver`

Both descend from the same code. `HIDDriver` is a Windows-10 compatibility fork
of the original `loki-hidriver` (whose user-mode code lived in the separate
`loki-example` repo). This is the file-by-file result of diffing them, which
informed what loki adopted.

## 1. Kernel driver (C sources) — effectively identical

`device.c/.h`, `driver.c/.h`, `hid.h`, `memory.c/.h`, `queue_default.c/.h`,
`queue_manual.c/.h` are **byte-for-byte identical** (ignoring CRLF/BOM). The HID
report descriptor, the report structs, and the IOCTL handling
(`IOCTL_HID_WRITE_REPORT` / `SET_OUTPUT_REPORT` → forwarded to the manual queue →
delivered as an input report) are unchanged.

The report protocol (from `hid.h`), which loki's Go and C APIs ultimately speak:

| Report ID | Name            | Payload                                   |
|-----------|-----------------|-------------------------------------------|
| `0x01`    | mouse input     | buttons, x, y (signed) — delivered to OS  |
| `0x02`    | mouse output    | what user mode writes; driver maps →0x01  |
| `0x03`    | keyboard input  | modifiers, reserved, 6 key codes → OS     |
| `0x04`    | keyboard output | what user mode writes; driver maps →0x03  |

### Differences that do exist in the driver

- **`hidriver.inf`**
  - `DriverVer` bumped `01/07/2018` → `10/10/2021`.
  - Install section renamed `DefaultInstall*` → `MyInstall*`.
  - `[Strings]` placeholders (`VARIABLE_1..6`) replaced with concrete names:
    `ClassName="HIDRIVER"`, `ServiceName="HIDRIVERSVC"`, etc. **This matters**:
    the user-mode lib finds the device by the string `HID#HIDRIVER&Col02` /
    `&Col04`, so the INF `ClassName` and the lib's device path are coupled.
- **`KMDFDriver.vcxproj`** — `HIDDriver` targets the latest platform
  (`$(LatestTargetPlatformVersion)`, empty `TargetVersion`/`DriverTargetPlatform`
  so it builds on a modern WDK) and sets `LanguageStandard=stdcpp17`. The
  original pinned `Windows10` / `Universal` / `Windows7`+`KMDF 1.9`.

## 2. User-mode lib — `HIDDriver/HIDDriverLib` vs `loki-example`

`loki-example` is a Qt console **app** (`main.cpp`, `.pro`). `HIDDriverLib` is
that same code repackaged as a **static library** (VS project, precompiled
headers). `device.cpp/.h`, `registry.cpp/.h`, `keys.h` are the same logic.

Meaningful differences (all in `HIDDriver`, all adopted by loki):

1. **Device interface path** — the original hard-codes the placeholder
   `\\?\HID#VARIABLE_6&Col02#1` (mouse) / `&Col04#1` (keyboard). `HIDDriver`
   changes it to `\\?\HID#HIDRIVER&Col02#1` to match the filled-in INF
   `ClassName`. Without this the lib can't find the device.
2. **New mouse motion functions** in `mouse.cpp/.h`:
   - `moveCursor(x1,y1,x2,y2, z, mouseMoveSlow)` — an aim-assist style move
     designed for 3D games: computes speed from the delta between a "center"
     point and a detected target, scaled by a `z` distance factor.
   - `moveCursorEx(x, y)` — a simpler relative/stepped move.
   - Several helpers (`getSpeedByRange`, `sendMouseReport`) were promoted to
     `public` so callers can send raw reports.
3. **VS scaffolding** — `#include "pch.h"` / `framework.h` added to each `.cpp`;
   comments are GBK-encoded Chinese (loki converts them to UTF-8).
4. **A test project** (`HIDDriverLibTest`) exercising the static lib.

## 3. What loki changes on top of `HIDDriver`

- **`loki_c_api.h/.cpp`** — a new flat `extern "C"` surface (opaque handles,
  status codes, per-thread `loki_last_error`) built as **`loki.dll`**. Neither
  upstream had a stable C ABI; both exposed C++ classes (or a static lib) only.
- **`Keyboard::setModifiers` / `Keyboard::send`** — added because upstream had no
  public way to hold modifiers or send multi-key reports; `type()` did one key
  with no modifier control.
- **Go bindings** — the whole point of this repo; upstream is C++ only.
- Encoding normalized to UTF-8; a `cl.exe` `build.bat` and a DLL `.vcxproj`
  added; driver/lib renamed to the `loki-*` projects.

### One coupling to remember
If you ever rebrand the device from `HIDRIVER` to, say, `LOKI`, you must change
**both** the INF `[Strings] ClassName` **and** the two device-path strings in
`windows/lib/mouse.cpp` / `keyboard.cpp` together, then reinstall the driver.
