# xbps

Are you tired of typing those _pesky hyphens_ when using Void Linux's package
manager? Do you believe that `xbps i` is objectively better than `xbps-install`?
Welcome to the resistance.

## What Even Is This?

This is a **zero-config, zero-deps, maximum-petty** dispatcher that provides an
alternative interface to xbps's hyphenated commands. We scan your PATH, find all
those `xbps-*` executables, and let you use them with spaces and smart aliases.

> **Note**: This is a loving tribute to xbps, not a criticism! The hyphenated
> design works perfectly fine - we're just having fun with alternatives.

## Installation

```bash
go install github.com/kendfss/xbps@latest
sudo ln -s ~/go/bin/xbps /usr/bin/xbps # this is so that it can be found by sudo
```

_or one-liner:_

```bash
curl -sL "https://github.com/kendfss/xbps/releases/latest/download/xbps_linux_$(uname -m).tar.gz" | tar -xz -O xbps | sudo tee /usr/bin/xbps >/dev/null && sudo chmod +x /usr/bin/xbps
```

## Usage

### Before (The Dark Times)

```bash
xbps-install -S helix
xbps-query -Rs nasm
xbps-remove -O
```

### After (Enlightenment)

```bash
xbps i -S neovim        # 'i' for install
xbps q -Rs uv           # 'q' for query  
xbps rem -Ooy           # 'rem' for remove (avoids clash with 'reconfigure')
xbps -h                 # see all commands and their generated aliases
```

_Your pinky finger will thank you for the reduced hyphen hunting._

## How It Works

We automatically discover and create intelligent shortcuts:

1. **Scan PATH** for all `xbps-*` executables
2. **Generate unique aliases** using a prefix trie (so `install` → `i`, but
   `remove` → `rem` to avoid conflicts)
3. **Route commands transparently** - all flags and arguments pass through
   unchanged
4. **Discover new commands automatically** - no configuration/updates needed

## FAQ

**Q: Why does this exist?**\
A: Mostly for fun! Also because implementing a trie-based alias system was more
interesting than writing shell aliases.

**Q: Is this production-ready?**\
A: It's surprisingly robust! The error handling and test suite are quite
thorough for a "fun" project.

**Q: Won't this break?** A: It's implemented so as to detect new hyphenated
commands as and when they appear and forget them once they've reached the end of
their lifespan.

**Q: Why not just use shell aliases?**\
A: Where's the fun in that? This automatically discovers new commands and
generates optimal aliases using proper data structures. If you do prefer that
though you might like [vpm][vpm]! It might be a little out of date but it should
work for the basics! More recently, and also not a shell script, there's
[xbps-tui][xbps-tui] which looks gorgeous!

<!-- **Q: Couldn't this just be a shell script?** A: at one point, certainly! but I -->
<!-- for one am glad I didn't subject myself to starting down a path that would lead -->
<!-- to me trying to implement a trie-based alias system in a posix compliant shell -->
<!-- scripting language! -->

## The Fine Print (We Have Lawyers?)

This is essentially a **loving shitpost with compile steps** (that simply
couldn't wait until April 1, 2026). It works surprisingly well, but we take no
responsibility if it somehow installs TempleOS on your system.

Run `xbps -h` to see all the commands and their smart aliases on your system.

## Build

```bash
git clone https://github.com/kendfss/xbps
cd xbps
go build
```

_Use `make build` or `make install` to include debugging information in local
builds._

## Requirements

- Void Linux (or system with xbps-* commands available)
- Go 1.24+ (for building from source)

---

_Join the movement. End hyphen tyranny. `xbps install freedom`._

[vpm]: https://github.com/bahamas10/vpm
[xbps-tui]: https://codeberg.org/lukeflo/xbps-tui
