# VIIPER backend for PureDS4

Project repository: [meiameiameia/PureDS4-VIIPER](https://github.com/meiameiameia/PureDS4-VIIPER).

This is a source fork of [hbashton/VIIPER](https://github.com/hbashton/VIIPER),
itself based on [Alia5/VIIPER](https://github.com/Alia5/VIIPER). The Git history
and [GPL-3.0 license](LICENSE.txt) are retained. The Go module path stays
`github.com/Alia5/VIIPER` to avoid changing internal imports and the client
protocol while this fork is validated.

The fork is limited to PureDS4's current virtual-controller contract: Xbox 360
and DualShock 4, including the DS4 v3 audio device variants. DualSense, Switch
2 Pro, keyboard, mouse, the shared-library build, and the built-in update
checker and inherited JSON/YAML/TOML configuration loading have been removed.
The command line exposes only the server; it does not install drivers, register
startup, or manage system configuration. PureDS4 owns backend installation,
version pinning, and replacement. Future DualShock 3 support is not implied by
this fork.

This repository is source code, not a standalone Windows installer or an
endorsement to replace a running backend. USB-IP driver installation and
machine-wide ownership remain PureDS4 installation concerns.

## Local verification

On Windows, with the Go version specified in `go.mod`, first run the source
checks and then generate the Windows version resource, executable, and current
dependency notices together:

```powershell
go test -count=1 ./...
go vet ./...
.\scripts\build-pureds4-backend.ps1 -OutputDir .\dist -RequireClean
```

The build alone does not establish compatibility with a physical DS4, USB-IP,
or a particular PureDS4 package. Release use requires pinning the exact source
revision and binary hash, license notices, installer verification, and the
PureDS4 USB/Bluetooth hardware gates.

The fork starts from hbashton's `v0.1.0` source revision
`fd298a04d7d229293be15b2af664405c9e68114c`, which matches the backend
previously bundled with PureDS4. Later upstream fixes are evaluated and
backported selectively; this is not a blanket tracking branch.
