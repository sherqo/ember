# Attribution

Ember's audio backend (`bin/omarchy-*`, `data/`) is taken from
[Omarchy](https://github.com/omacom/omarchy) by David Heinemeier Hansson,
which is released under the MIT License (see `LICENSE.omarchy`).

The vendored scripts are lightly adapted: the `$OMARCHY_PATH`-rooted data
paths in `omarchy-audio-tuning` now resolve beside the scripts so the
backend runs standalone. Everything else — the Go TUI, theme, Waybar
integration — is original work.

# Omarchy license

Copied from the Omarchy repository at time of vendoring.
