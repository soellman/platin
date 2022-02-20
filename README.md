# platin

Controls your platin hub

## What can it do?

- Turn speakers on/off
- Show and select sources
- Show current status of Platin hub
- Show and set the volume

## Installing

Requires an installation of golang.

Run `go install github.com/soellman/platin/cmd/platin@latest` to install the `platin` commandline tool.

Run `platin` for instructions.

Example:

```
> platin -a 10.0.0.10 status
Power: on
Volume: 50
Sources:
 - Line
 - AUX
 - HDMI
 * OPT1
 - OPT2
 - OPT3
 - COAX
 - USB
 - Bluetooth
 - Streaming
>
```
