# Ember — warm little audio mixer

Output + mic levels, sink/source picker, per-app volumes.
Pink Cat Boo, CaskaydiaMono Nerd Font, zero idle RAM (exits on quit).

## Run
`ember` (Go binary; Waybar audio click opens it floating).

## Build
`go build -o bin/ember ember.go` (stdlib + Bubble Tea/Lipgloss, see `go.mod`).

## Layout
- `ember.go` — the TUI.
- `bin/omarchy-*` — audio backend vendored from Omarchy (see `ATTRIBUTION.md`).
- `data/` — speaker-tuning data for the backend.
- `waybar/` — module snippet + Hyprland float rule.

## Controls
`↑↓/jk` move · `←→` or `-/+` volume (to 150%) · `m` mute · `enter` select
default · `r` refresh · `q/esc` quit. Mouse: click select, click bar to
set volume, wheel adjusts, right-click mutes.

## License
MIT — see `LICENSE`. Omarchy-derived files keep their original terms;
see `ATTRIBUTION.md` and `LICENSE.omarchy`.
