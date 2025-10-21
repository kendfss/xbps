xbps
---

Are you tired of typing those *pesky hyphens* when using Void Linux's package manager? Do you believe that `xbps install` is objectively better than `xbps-install`? Welcome to the resistance.


## What Even Is This?
This is a **zero-config, zero-deps, maximum-petty** wrapper that fixes xbps's unfortunate hyphen addiction. We scan your PATH, find all those `xbps-*` commands, and let you use them like God intended: WITH SPACES.


## Installation (The Part Where We Pretend This Is Serious)
```bash
go install github.com/kendfss/xbps@latest
sudo ln -s ~/go/bin/xbps /usr/bin/xbps # this is so that it can be found by sudo
```

*or*

```bash
curl -sL "https://github.com/kendfss/xbps/releases/latest/download/xbps_linux_$(uname -m).tar.gz" | tar -xz -O xbps | sudo tee /usr/bin/xbps >/dev/null && sudo chmod +x /usr/bin/xbps
```

## Usage (Finally, Sanity)

### Before (The Dark Times)
```bash
xbps-install -S helix
xbps-query -Rs nasm
xbps-remove -O
```

### After (Enlightenment)
```bash
xbps install -S neovim
xbps query -Rs uv  
xbps remove -Ooy
```

*See how much better that looks? Your pinky finger thanks you for not making it hunt hyphens all day.*


## What Magic Is This?
We do the thing you wish xbps did out of the box:

1. **Raid your PATH** for any `xbps-*` executables
2. **Chop off the hyphen** because it had no business being there
3. **Route your commands** like a proper adult package manager
4. **Pass everything through** - flags, args, your existential dread


## FAQ (Frequently Annoyed Questions)
**Q: Why does this exist?**
A: Because someone had to fix the hyphen crime.

**Q: Is this production-ready?**
A: Is anything, really?

**Q: Won't this break?**
A: Only if xbps suddenly becomes sensible and adds proper subcommands.

**Q: Why not just alias everything?**
A: Where's the fun in that? Also, this automatically discovers new commands.


## The Fine Print (We Have Lawyers?)
This is essentially a shitpost with compile steps. It works surprisingly well, but we take no responsibility if it somehow installs TempleOS on your system.

Run `xbps -h` to see all the commands we've liberated from hyphen hell on your system.


## Build
```bash
git clone https://github.com/kendfss/xbps
cd xbps
go build
```

*you can also use `make build` or `make install` to include debugging information in your local builds*


## Requirements
- Void Linux (or system with xbps-* commands available)
- Go 1.24+ (for building from source)


---

*Join the movement. End hyphen tyranny. `xbps install freedom`.*
