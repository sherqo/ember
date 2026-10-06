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

## Install (Arch Linux)

Prerequisites (runtime):

```bash
sudo pacman -S --needed wireplumber libpulse playerctl ttf-cascadia-mono-nerd
```

- `wireplumber` provides `wpctl`, `libpulse` provides `pactl`.
- `playerctl` is optional — without it the APPS section loses its
  active-player highlight, everything else works.
- A Nerd Font is needed for the speaker/mic glyphs.

Build requirements: `go` (`sudo pacman -S go`).

Install from source:

```bash
git clone https://github.com/sherqo/ember.git ~/ember
cd ~/ember
make install   # builds + links ~/.local/bin/ember
```

`make install` symlinks the binary (rather than copying it) because
`ember` resolves its `omarchy-*` backend helpers relative to its own
directory. Keep the clone around. `make uninstall` removes the link.

Waybar (optional) — pulseaudio click opens Ember floating:

```json
"on-click": "alacritty --class org.sherqo.ember -e ember",
```

Hyprland (optional) — float rule, e.g. in `windows.lua`:

```lua
hl.window_rule({ match = { class = "org.sherqo.ember" }, float = true, center = true, size = { 620, 460 } })
```

## Controls
`↑↓/jk` move · `←→` or `-/+` volume (to 150%) · `m` mute · `enter` select
default · `r` refresh · `q/esc` quit. Mouse: click select, click bar to
set volume, wheel adjusts, right-click mutes.

## License
MIT — see `LICENSE`. Omarchy-derived files keep their original terms;
see `ATTRIBUTION.md` and `LICENSE.omarchy`.
