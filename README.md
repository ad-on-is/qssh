# qSSH

A simple wrapper for `ssh` to select profiles from your `~/.ssh/config` file. Can be aliased to `ssh`

![qssh](https://i.imgur.com/ZNekyDI.png)

## Installation

- Download binary from [Releases](https://github.com/ad-on-is/qssh/releases), or build it yourself.
- Rename binary to qssh.
- (Optional) add alias to replace ssh

## Usage

- `qssh` Shows a list of profiles specified in the config-file
- `qssh <param>` Checks whether a profile with the name of `<param>` exists, if so opens that profile, else executes `ssh <param>`
- `qssh <param> <param> ...` Passes all params directly to `ssh <param> <param> ...`

## Enhancements

Add #Title and/or #Description to hosts

```yaml
Host Example
  #Title Some title with spaces
  #Description Some description, like, this host needs VPN prior to connecting, etc...
  HostName example.com
  User foo
```

## Build yourself

Simply run `go build && go install`. Optionally you can comment out not needed `PLATFORMS=` and run `./build`.
