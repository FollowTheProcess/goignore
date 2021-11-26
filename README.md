# goignore

[![License](https://img.shields.io/github/license/FollowTheProcess/goignore)](https://github.com/FollowTheProcess/goignore)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/goignore)](https://goreportcard.com/report/github.com/FollowTheProcess/goignore)
[![CI](https://github.com/FollowTheProcess/goignore/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/goignore/actions?query=workflow%3ACI)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/goignore?logo=github&sort=semver)](https://github.com/FollowTheProcess/goignore)

An extremely simple go CLI to hit the [gitignore API] with whatever you pass as command line arguments. The list of things you can pass here are documented on [gitignore.io].

You'll get back a .gitignore file saved to `$CWD/.gitignore` with the contents generated from the API.

## Installation

```shell
brew tap FollowTheProcess/homebrew-tap
brew install FollowTheProcess/homebrew-tap/goignore
```

Or download the correct binary for your OS and ARCH from the GitHub Releases :grin:

## Usage

Inside the folder you want the `.gitignore` to live in, run:

```shell
goignore macos vscode go
```

This will get you a `.gitignore` file that looks like...

```plaintext
# Created by https://www.toptal.com/developers/gitignore/api/macos,vscode,go
# Edit at https://www.toptal.com/developers/gitignore?templates=macos,vscode,go

### Go ###
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

### Go Patch ###
/vendor/
/Godeps/

### macOS ###
# General
.DS_Store
.AppleDouble
.LSOverride

# Icon must end with two \r
Icon


# Thumbnails
._*

# Files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

### vscode ###
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
*.code-workspace

# End of https://www.toptal.com/developers/gitignore/api/macos,vscode,go
```

## List Options

If you're not sure what you can type in check on [gitignore.io] or run:

```shell
goignore --list
```

And you'll see something like...

```shell
1c,
1c-bitrix,
a-frame,
actionscript,
ada
adobe,
advancedinstaller,
adventuregamestudio,
agda,
al

# etc.
```

If you have a particular thing in mind:

```shell
goignore --list | grep visualstudiocode

visualstudiocode
```

## Help

When in doubt, run:

```shell
goignore --help
```

[gitignore API]: https://www.toptal.com/developers/gitignore
[gitignore.io]: https://www.toptal.com/developers/gitignore
