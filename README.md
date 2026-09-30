# Ember — warm little audio mixer

Output + mic sliders, sink/source picker, per-app volumes.
Pink Cat Boo, CaskaydiaMono Nerd Font, zero idle RAM (exits on quit).

## Run
`ember` (Go binary; Waybar audio click opens it floating). `ember-fzf`
is the older fzf fallback. `ember-launch` opens the dormant QML panel
(needs quickshell; kept for later).

## Build
`go build -o bin/ember ember.go` (stdlib only, no modules).

## Layout
- `ember.go` — the TUI (this is the app now).
- `bin/` — backend: vendored Omarchy audio helpers (`omarchy-audio-*`,
  `omarchy-hw-match`, `omarchy-restart-audio`, `omarchy-cmd-present`) + `data/`
  speaker-tuning data + `ember-fzf` fallback + `ember-launch` (QML).
- `qml/` + `app/` — dormant Quickshell frontend (needs quickshell installed).
- `waybar/` — module snippet + Hyprland float rule.

## Controls
`↑↓/jk` move · `←→` or `-/+` volume · `m` mute · `enter` select/mute ·
`r` refresh · `q/esc` quit.
