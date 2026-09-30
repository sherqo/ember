// Ember — warm little audio mixer (Go TUI, stdlib only).
// Mirrors the Omarchy audio panel order: hero, mic, outputs, inputs, apps.
// Keys: up/down or j/k move · left/right or -/+ volume · m mute ·
// enter select/mute · r refresh · q/esc quit.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Pink Cat Boo
const (
	cReset   = "\x1b[0m"
	cFg      = "\x1b[38;2;255;240;245m"
	cAccent  = "\x1b[38;2;255;76;122m"
	cMuted   = "\x1b[38;2;86;89;112m"
	cDim     = "\x1b[38;2;120;124;150m"
	cGreen   = "\x1b[38;2;59;192;137m"
	cBgBar   = "\x1b[48;2;53;55;70m"
	cSelBg   = "\x1b[48;2;64;50;70m"
	cBold    = "\x1b[1m"
	altOn    = "\x1b[?1049h\x1b[H"
	altOff   = "\x1b[?1049l"
	hideCur  = "\x1b[?25l"
	showCur  = "\x1b[?25h"
	clearAll = "\x1b[H\x1b[2J"
)

type Node struct {
	ID     string
	Name   string
	Kind   string // OUTPUT, INPUT, SINK, SOURCE, STREAM
	Pct    int
	Muted  bool
	Active bool
}

func sh(args ...string) string {
	out, _ := exec.Command(args[0], args[1:]...).Output()
	return string(out)
}

func volOf(id string) (int, bool) {
	out := sh("wpctl", "get-volume", id)
	f := strings.Fields(out)
	if len(f) < 2 {
		return 0, false
	}
	v, _ := strconv.ParseFloat(f[1], 64)
	muted := strings.Contains(out, "MUTED")
	return int(v*100 + 0.5), muted
}

func sectionNodes(status, section string) [][2]string {
	var res [][2]string
	seen := map[string]bool{}
	inSec := false
	for _, line := range strings.Split(status, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "├─") {
			inSec = strings.Contains(t, section)
			continue
		}
		if !inSec || t == "" || !strings.HasPrefix(t, "│") {
			continue
		}
		// match "│  60. Name [vol: x]"
		s := strings.TrimLeft(t, "│ ")
		dot := strings.Index(s, ".")
		if dot < 0 {
			continue
		}
		id := strings.TrimSpace(s[:dot])
		if _, err := strconv.Atoi(id); err != nil {
			continue
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		name := strings.TrimSpace(s[dot+1:])
		if i := strings.Index(name, " [vol:"); i >= 0 {
			name = name[:i]
		}
		res = append(res, [2]string{id, name})
	}
	return res
}

func bar(pct int) string {
	const w = 16
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	n := pct * w / 100
	return strings.Repeat("█", n) + strings.Repeat("░", w-n)
}

func snapshot() []Node {
	status := sh("wpctl", "status")
	defSink := strings.TrimSpace(sh("pactl", "get-default-sink"))
	defSrc := strings.TrimSpace(sh("pactl", "get-default-source"))
	var rows []Node

	p, m := volOf("@DEFAULT_AUDIO_SINK@")
	rows = append(rows, Node{"@DEFAULT_AUDIO_SINK@", defSink, "OUTPUT", p, m, true})
	p, m = volOf("@DEFAULT_AUDIO_SOURCE@")
	rows = append(rows, Node{"@DEFAULT_AUDIO_SOURCE@", defSrc, "INPUT", p, m, true})

	defID := nodeIDByName(status, defSink)
	defSrcID := nodeIDByName(status, defSrc)
	_ = defSrcID
	for _, nn := range sectionNodes(status, "Sinks") {
		p, m := volOf(nn[0])
		rows = append(rows, Node{nn[0], nn[1], "SINK", p, m, nn[0] == defID})
	}
	for _, nn := range sectionNodes(status, "Sources") {
		p, m := volOf(nn[0])
		rows = append(rows, Node{nn[0], nn[1], "SOURCE", p, m, false})
	}
	for _, nn := range sectionNodes(status, "Streams") {
		p, m := volOf(nn[0])
		rows = append(rows, Node{nn[0], nn[1], "STREAM", p, m, false})
	}
	return rows
}

func nodeIDByName(status, name string) string {
	for _, line := range strings.Split(status, "\n") {
		if strings.Contains(line, name) {
			s := strings.TrimLeft(strings.TrimSpace(line), "│ ")
			if dot := strings.Index(s, "."); dot > 0 {
				if _, err := strconv.Atoi(strings.TrimSpace(s[:dot])); err == nil {
					return strings.TrimSpace(s[:dot])
				}
			}
		}
	}
	return ""
}

func shortName(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

func render(rows []Node, cur int) string {
	var b strings.Builder
	b.WriteString(clearAll)
	b.WriteString(cBold + cFg + "  Ember" + cReset + cDim + "  ·  +- volume · m mute · enter select/mute · r refresh · q quit" + cReset + "\n\n")
	lastKind := ""
	for i, r := range rows {
		if r.Kind != lastKind && r.Kind != "OUTPUT" && r.Kind != "INPUT" {
			if r.Kind == "SINK" {
				b.WriteString("\n" + cDim + "  OUTPUTS" + cReset + "\n")
			} else if r.Kind == "SOURCE" {
				b.WriteString("\n" + cDim + "  INPUTS" + cReset + "\n")
			} else if r.Kind == "STREAM" {
				b.WriteString("\n" + cDim + "  APPS" + cReset + "\n")
			}
			lastKind = r.Kind
		}
		sel := i == cur
		name := shortName(r.Name, 34)
		vol := cGreen + bar(r.Pct) + cReset + cFg + fmt.Sprintf(" %3d%%", r.Pct) + cReset
		if r.Muted {
			vol = cMuted + bar(r.Pct) + fmt.Sprintf(" %3d%% MUTED", r.Pct) + cReset
		}
		active := ""
		if r.Active {
			active = cAccent + " ●" + cReset
		}
		prefix := "  "
		if sel {
			prefix = cAccent + "▸ " + cReset
		}
		fmt.Fprintf(&b, "%s%s%-36s %s%s\n", prefix, cFg, name, vol, active)
		if sel {
			b.WriteString(cReset)
		}
	}
	return b.String()
}

var tty *os.File

func rawOn() {
	tty, _ = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1", "-echo").Run()
	fmt.Print(altOn + hideCur)
}

func rawOff() {
	fmt.Print(showCur + altOff)
	exec.Command("stty", "-F", "/dev/tty", "sane").Run()
	if tty != nil {
		tty.Close()
	}
}

func readKey() string {
	buf := make([]byte, 8)
	n, _ := tty.Read(buf)
	if n == 0 {
		return ""
	}
	if buf[0] == 0x1b {
		if n == 1 {
			return "esc"
		}
		switch string(buf[1:n]) {
		case "[A":
			return "up"
		case "[B":
			return "down"
		case "[C":
			return "right"
		case "[D":
			return "left"
		}
		return "esc"
	}
	switch buf[0] {
	case 'q', 'Q':
		return "quit"
	case 'r', 'R':
		return "refresh"
	case 'm', 'M':
		return "mute"
	case '+', '=':
		return "up-vol"
	case '-', '_':
		return "down-vol"
	case '\r', '\n':
		return "enter"
	case 'j':
		return "down"
	case 'k':
		return "up"
	case 'h':
		return "left"
	case 'l':
		return "right"
	}
	return ""
}

func adjust(id string, delta string) {
	exec.Command("wpctl", "set-volume", id, delta).Run()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--dump" {
		for _, r := range snapshot() {
			fmt.Printf("%s\t%s\t%s\t%d\t%v\n", r.Kind, r.ID, r.Name, r.Pct, r.Muted)
		}
		return
	}
	rows := snapshot()
	cur := 0
	rawOn()
	defer rawOff()
	fmt.Print(render(rows, cur))

	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	keych := make(chan string, 8)
	go func() {
		for {
			keych <- readKey()
		}
	}()

	for {
		select {
		case k := <-keych:
			switch k {
			case "quit", "esc":
				return
			case "refresh":
				rows = snapshot()
			case "up":
				if cur > 0 {
					cur--
				}
			case "down":
				if cur < len(rows)-1 {
					cur++
				}
			case "up-vol", "right":
				if cur < len(rows) {
					adjust(rows[cur].ID, "5%+")
					rows = snapshot()
				}
			case "down-vol", "left":
				if cur < len(rows) {
					adjust(rows[cur].ID, "5%-")
					rows = snapshot()
				}
			case "mute":
				if cur < len(rows) {
					exec.Command("wpctl", "set-mute", rows[cur].ID, "toggle").Run()
					rows = snapshot()
				}
			case "enter":
				if cur < len(rows) {
					r := rows[cur]
					switch r.Kind {
					case "SINK", "SOURCE":
						exec.Command("wpctl", "set-default", r.ID).Run()
					default:
						exec.Command("wpctl", "set-mute", r.ID, "toggle").Run()
					}
					rows = snapshot()
				}
			}
			fmt.Print(render(rows, cur))
		case <-tick.C:
			rows = snapshot()
			if cur >= len(rows) && len(rows) > 0 {
				cur = len(rows) - 1
			}
			fmt.Print(render(rows, cur))
		}
	}
}
