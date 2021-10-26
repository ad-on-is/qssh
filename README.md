# qSSH

A simple wrapper for `ssh` to select profiles from your `~/.ssh/config` file. Can be aliased to `ssh`

![qssh](https://i.imgur.com/ZNekyDI.png)

## Installation

- Download binary from releases, or build iti yourself.
- Rename binary to qssh.
- (Optional) add alias to replace ssh

## Usage

- `qssh` Shows a list of profiles specified in the config-file
- `qssh <param>` Checks whether a profile with the name of `<param>` exists, if so opens that profile, else executes `ssh <param>`
- `qssh <param> <param> ...` Passes all params directly to `ssh <param> <param> ...`

## Build yourself

Simply run `./build`. Optionally you can comment out not needed `PLATFORMS=` in `build`.
