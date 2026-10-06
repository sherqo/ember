# Ember — warm little audio mixer

Output + mic levels, sink/source picker, per-app volumes.
Pink Cat Boo, CaskaydiaMono Nerd Font, zero idle RAM (exits on quit).

## Run
`ember` (Go binary; Waybar audio click opens it floating).

## Build
`go build -o ember .` (Bubble Tea/Lipgloss, see `go.mod`),
or straight from the network with no clone:

```bash
go install github.com/sherqo/ember@latest
```

## Install (Arch Linux)

Runtime deps (the backend shells out to them):

```bash
sudo pacman -S --needed wireplumber libpulse
sudo pacman -S --needed playerctl # optional: active-player highlight in APPS
```

A Nerd Font is recommended for the speaker/mic glyphs, but not required —
without one Ember falls back to plain ASCII marks automatically
(`EMBER_ASCII=1` forces the fallback).

The backend ships inside the binary: on first run Ember extracts its
helpers to `~/.local/share/ember` (or `$XDG_DATA_HOME/ember`) and
re-syncs them whenever they change. No Makefile, no manual copying.

## Layout
- `ember.go` — the TUI.
- `bin/omarchy-*` — audio backend vendored from Omarchy (see `ATTRIBUTION.md`),
  embedded into the binary via `go:embed`.
- `data/` — speaker-tuning data for the backend.
- `waybar/` — module snippet + Hyprland float rule.

## Controls
`↑↓/jk` move · `←→` or `-/+` volume (to 150%) · `m` mute · `enter` select
default · `r` refresh · `q/esc` quit. Mouse: click select, click bar to
set volume, wheel adjusts, right-click mutes.

## License
MIT — see `LICENSE`. Omarchy-derived files keep their original terms;
see `ATTRIBUTION.md` and `LICENSE.omarchy`.
