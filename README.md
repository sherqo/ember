# Ember — warm little audio mixer

Output + mic sliders, sink/source picker, per-app volumes.
Pink Cat Boo, CaskaydiaMono Nerd Font, zero idle RAM (exits on quit).

## Run
`ember` (floating terminal via Waybar audio module click, or any terminal).

## Layout
- `bin/` — backend: vendored Omarchy audio helpers (`omarchy-audio-*`,
  `omarchy-hw-match`, `omarchy-restart-audio`, `omarchy-cmd-present`) + `data/`
  speaker-tuning data. `ember` itself drives PipeWire directly via `wpctl`.
- `qml/` — dormant Quickshell frontend (needs quickshell installed;
  kept for later, same look, see PLAN in /tmp/sherqo-apps-plan.md).
- `waybar/` — module snippet + Hyprland float rule (added at integration).

## Controls
- `+`/`-` volume ±5% on the focused section, `m` mute, `enter` set default,
  `tab` next section, `r` refresh, `q` quit.
