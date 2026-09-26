# loki virtual HID driver

A KMDF kernel-mode driver that registers a virtual HID device with a **mouse**
collection and a **keyboard** collection. User mode writes HID *output* reports
(via `loki.dll`); the driver forwards them as HID *input* reports, so Windows
delivers real input events.

The C sources are unchanged from
[hedgar2017/loki-hidriver](https://github.com/hedgar2017/loki-hidriver) /
[dengqizhou30/HIDDriver](https://github.com/dengqizhou30/HIDDriver); see
[../../docs/comparison.md](../../docs/comparison.md).

## Build

Requires **Visual Studio** + the **Windows Driver Kit (WDK)** matching your VS
version. Build `loki-driver.vcxproj` as `Release|x64`. Output: `hidriver.sys`,
`hidriver.inf`, and (when signed) `hidriver.cat`.

## Generating a branded INF

The install strings (`DeviceName`, `DiskName`, `ProviderName`,
`ManufacturerName`, `ServiceName`) are brandable and support **Chinese**. Use the
`mkinf` tool to fill them in and emit a correctly-encoded `hidriver.inf`:

```bash
# from repo root; flags:
go run ./cmd/mkinf -out windows/driver/hidriver.inf \
  -device "Loki 虚拟键鼠" -disk "Loki 安装盘" \
  -provider hkloudou -manufacturer "杭州云侯科技" -service "Loki 虚拟键鼠服务"

# or from a JSON config (see branding.example.json), with a proper zh-CN section:
go run ./cmd/mkinf -config windows/driver/branding.example.json \
  -out windows/driver/hidriver.inf
```

Notes:
- Non-ASCII values make the tool write **UTF-16LE + BOM** (required for Unicode
  INFs). A `zhCN` block in the JSON emits a localized `[Strings.0804]` section
  with an ASCII neutral `[Strings]` fallback.
- **`ClassName` is coupled to the user-mode lib**: the device is found via
  `\\?\HID#<ClassName>&Col02#1` / `&Col04#1`. Keep it `HIDRIVER` unless you also
  edit `windows/lib/mouse.cpp` / `keyboard.cpp`; `mkinf` warns if you change it.
- Any INF change invalidates the catalog — **regenerate and re-sign
  `hidriver.cat`** afterwards (`inf2cat` + `signtool`), or for a shipped product
  choose branding before the one-time Microsoft attestation submission.

## Install (developer / test mode)

The default build is **test-signed only**, so it needs signature enforcement
relaxed. In an **elevated** command prompt:

```bat
bcdedit /set testsigning on
:: reboot
devcon install hidriver.inf root\hidriver
devcon status root\hidriver
```

Remove / disable:

```bat
devcon remove root\hidriver
bcdedit /set testsigning off
```

Install logs: `C:\Windows\INF\setupapi.dev.log`.

## Shipping to real users

Test mode is **not** shippable. To load on normal machines the driver must be
properly signed — for a software driver like this that means **Microsoft
attestation signing** via Partner Center, using an **EV code-signing
certificate**. See [../../docs/signing.md](../../docs/signing.md), which also
covers the trade-offs vs. going driverless (SYSTEM service + `SendInput`, or an
external USB-HID microcontroller) the way ToDesk / 影刀 / UiPath do.

## Notes

- `VENDOR_ID` / `PRODUCT_ID` / `VERSION_NUMBER` are all `0x00` in `device.h`.
- The device is discovered from user mode by the string `HID#HIDRIVER&Col02`
  (mouse) / `&Col04` (keyboard), which is derived from `[Strings] ClassName` in
  `hidriver.inf`. Keep the INF `ClassName` and the lib's device paths in sync.
