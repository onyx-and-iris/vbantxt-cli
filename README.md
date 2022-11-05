# VBAN Sendtext CLI Utility

VBAN sendtext cli utility for sending Voicemeeter string requests over a network.

## Tested against

-   Basic 1.0.8.4
-   Banana 2.0.6.4
-   Potato 3.0.2.4

## Requirements

-   [Voicemeeter](https://voicemeeter.com/)
-   Go 1.18 or greater

## `Use`

#### Command Line

Pass `host`, `port` and `streamname` as flags, for example:

`vbantxt-cli -h="gamepc.local" -p=6980 -s=Command1 "strip[0].mute=1 strip[1].mono=1"`

You may also store them in a `config.toml` located in `home directory / .vbantxt_cli /`

A valid `config.toml` might look like this:

```toml
[connection]
Host="gamepc.local"
Port=6990
Streamname="Command1"
```

#### Script files

The vbantxt-cli utility accepts a single string request or an array of string requests. This means you can pass scripts stored in files.

For example, in Windows with Powershell you could:

`vbantxt-cli $(Get-Content .\script.txt)`

to load commands from a file:

```
strip[0].mute=0;strip[0].mute=0
strip[1].mono=0;strip[1].mono=0
```
