package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pink Cat Boo
var (
	cBg     = lipgloss.Color("#202330")
	cFg     = lipgloss.Color("#FFF0F5")
	cAccent = lipgloss.Color("#FF4C7A")
	cMuted  = lipgloss.Color("#565970")
	cDim    = lipgloss.Color("#8A8DA3")
	cGreen  = lipgloss.Color("#3BC089")
	cBlue   = lipgloss.Color("#6767CE")
	cYellow = lipgloss.Color("#FEC831")
	titleSt = lipgloss.NewStyle().Bold(true).Foreground(cAccent)
	subSt   = lipgloss.NewStyle().Foreground(cDim)
	nameSt  = lipgloss.NewStyle().Foreground(cFg)
	dimSt   = lipgloss.NewStyle().Foreground(cMuted)
	selSt   = lipgloss.NewStyle().Foreground(cFg).Background(lipgloss.Color("#3A3048")).Bold(true)
	headSt  = lipgloss.NewStyle().Foreground(cFg).Bold(true).MarginTop(1)
	boxSt   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(cMuted).Padding(1, 2).Background(cBg)
	helpSt  = lipgloss.NewStyle().Foreground(cMuted)
	pctSt   = lipgloss.NewStyle().Foreground(cFg).Width(5).Align(lipgloss.Right)
	mutSt   = lipgloss.NewStyle().Foreground(cYellow).Bold(true)
)

type Node struct {
	ID     string
	Raw    string // PipeWire node.name (for set-default scripts)
	Name   string // friendly display label
	Kind   string
	Pct    int
	Muted  bool
	Active bool
}

func runCtx(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, _ := exec.CommandContext(ctx, name, args...).Output()
	return string(out)
}

func sh(args ...string) string { return runCtx(2*time.Second, args[0], args[1:]...) }

func volOf(id string) (int, bool) {
	out := sh("wpctl", "get-volume", id)
	f := strings.Fields(out)
	if len(f) < 2 {
		return 0, false
	}
	v, _ := strconv.ParseFloat(f[1], 64)
	return int(v*100 + 0.5), strings.Contains(out, "MUTED")
}

func sectionNodes(status, section string) [][2]string {
	var res [][2]string
	seen := map[string]bool{}
	inSec := false
	for _, line := range strings.Split(status, "\n") {
		t := strings.TrimSpace(line)
		// section headers use either ├─ or └─ tree branches
		if strings.HasPrefix(t, "├─") || strings.HasPrefix(t, "└─") {
			inSec = strings.Contains(t, section)
			continue
		}
		if !inSec || t == "" {
			continue
		}
		// node lines: "│  60. Name [vol: x]" or bare "140. App" (streams);
		// a leading "*" marks the default node — strip it before parsing.
		s := strings.TrimLeft(t, "│ ")
		s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "*"))
		dot := strings.Index(s, ".")
		if dot < 0 {
			continue
		}
		id := strings.TrimSpace(s[:dot])
		if _, err := strconv.Atoi(id); err != nil || seen[id] {
			continue
		}
		seen[id] = true
		name := strings.TrimSpace(s[dot+1:])
		if i := strings.Index(name, " [vol:"); i >= 0 {
			name = name[:i]
		}
		if i := strings.Index(name, "\t"); i >= 0 {
			name = strings.TrimSpace(name[:i])
		}
		if name == "" {
			continue
		}
		res = append(res, [2]string{id, name})
	}
	return res
}

func friendlyDeviceLabel(text string) string {
	label := strings.TrimSpace(text)
	stripPrefix := []string{"sof-soundwire ", "built-in audio ", "built in audio "}
	lower := strings.ToLower(label)
	for _, p := range stripPrefix {
		if strings.HasPrefix(lower, p) {
			label = strings.TrimSpace(label[len(p):])
			lower = strings.ToLower(label)
		}
	}
	for _, suf := range []string{" Output", " Input", " output", " input"} {
		if strings.HasSuffix(label, suf) {
			label = strings.TrimSpace(label[:len(label)-len(suf)])
		}
	}
	label = strings.ReplaceAll(label, "Microphones", "Microphone")
	return label
}

func nodeIDByName(status, name string) string {
	for _, line := range strings.Split(status, "\n") {
		if name != "" && strings.Contains(line, name) {
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

func inspectField(id, field string) string {
	out := sh("wpctl", "inspect", id)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "* ")
		if strings.HasPrefix(line, field+" = ") {
			v := strings.TrimSpace(strings.TrimPrefix(line, field+" = "))
			return strings.Trim(v, "\"")
		}
	}
	return ""
}

func snapshot() []Node {
	status := sh("wpctl", "status")
	defSink := strings.TrimSpace(sh("pactl", "get-default-sink"))
	defSrc := strings.TrimSpace(sh("pactl", "get-default-source"))
	// Omarchy abstraction: only available outputs are selectable.
	avail := map[string]bool{}
	for _, line := range strings.Split(sh("omarchy-audio-sink-availability"), "\n") {
		f := strings.Split(line, "\t")
		if len(f) >= 2 {
			avail[strings.TrimSpace(f[0])] = strings.TrimSpace(f[1]) != "0"
		}
	}
	nodeName := map[string]string{}
	labelOf := func(id, fallback string) string {
		nname := inspectField(id, "node.name")
		nodeName[id] = nname
		if nick := inspectField(id, "node.nick"); nick != "" {
			return nick
		}
		if desc := friendlyDeviceLabel(inspectField(id, "node.description")); desc != "" {
			return desc
		}
		return friendlyDeviceLabel(fallback)
	}
	var rows []Node
	p, m := volOf("@DEFAULT_AUDIO_SINK@")
	rows = append(rows, Node{"@DEFAULT_AUDIO_SINK@", defSink, defSink, "OUTPUT", p, m, true})
	p, m = volOf("@DEFAULT_AUDIO_SOURCE@")
	rows = append(rows, Node{"@DEFAULT_AUDIO_SOURCE@", defSrc, defSrc, "INPUT", p, m, true})
	defID := nodeIDByName(status, defSink)
	// Omarchy parity: mark streams belonging to the active MPRIS player.
	activePlayer := strings.ToLower(strings.TrimSpace(sh("playerctl", "metadata", "--format", "{{playerName}}")))
	for _, nn := range sectionNodes(status, "Sinks") {
		nname := nodeName[nn[0]]
		if nname == "" {
			nname = inspectField(nn[0], "node.name")
			nodeName[nn[0]] = nname
		}
		if nname == defSink {
			defID = nn[0]
		}
		if !avail[nname] && nn[0] != defID {
			lower := strings.ToLower(nname + " " + nn[1])
			if !strings.Contains(lower, "bluez") && !strings.Contains(lower, "bluetooth") {
				continue
			}
		}
		p, m := volOf(nn[0])
		rows = append(rows, Node{nn[0], nname, labelOf(nn[0], nn[1]), "SINK", p, m, nn[0] == defID})
	}
	for _, nn := range sectionNodes(status, "Sources") {
		p, m := volOf(nn[0])
		nname := nodeName[nn[0]]
		if nname == "" {
			nname = inspectField(nn[0], "node.name")
			nodeName[nn[0]] = nname
		}
		rows = append(rows, Node{nn[0], nname, labelOf(nn[0], nn[1]), "SOURCE", p, m, false})
	}
	for _, nn := range sectionNodes(status, "Streams") {
		// skip endpoint sub-rows ("out_FL > sink:playback_FL"); keep app rows
		if strings.Contains(nn[1], ">") {
			continue
		}
		p, m := volOf(nn[0])
		active := activePlayer != "" && strings.Contains(strings.ToLower(nn[1]), activePlayer)
		rows = append(rows, Node{nn[0], "", friendlyDeviceLabel(nn[1]), "STREAM", p, m, active})
	}
	return rows
}

func glyphFor(n Node) string {
	switch n.Kind {
	case "OUTPUT":
		if n.Muted {
			return ""
		}
		return ""
	case "INPUT":
		if n.Muted {
			return ""
		}
		return ""
	case "SINK":
		return ""
	case "SOURCE":
		return ""
	case "STREAM":
		return ""
	}
	return "•"
}

func friendlyHero(name string) string {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "bluez") || strings.Contains(n, "bluetooth"):
		return "Bluetooth"
	case strings.Contains(n, "hdmi") || strings.Contains(n, "displayport"):
		return "HDMI"
	case strings.Contains(n, "speaker"):
		return "Speaker"
	case strings.Contains(n, "headphone"):
		return "Headphones"
	case strings.Contains(n, "mic"):
		return "Microphone"
	}
	parts := strings.Split(name, ".")
	last := parts[len(parts)-1]
	if len(last) > 26 {
		last = last[:26]
	}
	if last == "" {
		return name
	}
	return last
}

func shortName(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

// rowLineY maps a data-row index to its terminal line in View.
// Layout: border(0) pad(1) title(2) hint(3) flash(4, always present) blank(5),
// heroes from 6; each new section adds blank + header (+2).
func rowLineY(rows []Node, idx int) int {
	y := 6
	lastKind := ""
	for i, r := range rows {
		if r.Kind != lastKind && r.Kind != "OUTPUT" && r.Kind != "INPUT" {
			y += 2
			lastKind = r.Kind
		}
		if i == idx {
			return y
		}
		y++
	}
	return -1
}

func rowAtY(rows []Node, y int) int {
	for i := range rows {
		if rowLineY(rows, i) == y {
			return i
		}
	}
	return -1
}

const barX0 = 37 // approx column where the 14-cell volume bar starts
const barW = 14

func volBar(pct, width int, muted bool) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 150 {
		pct = 150
	}
	fill := pct * width / 100
	if fill > width {
		fill = width
	}
	fg := cGreen
	if muted {
		fg = cMuted
	}
	f := lipgloss.NewStyle().Foreground(fg).Render(strings.Repeat("━", fill))
	e := lipgloss.NewStyle().Foreground(cMuted).Render(strings.Repeat("━", width-fill))
	return f + e
}

type model struct {
	rows   []Node
	cursor int
	width  int
	height int
	err    string
	flash  string
}

type refreshMsg []Node
type errMsg string

func doSnapshot() tea.Msg {
	// never block UI more than ~3s; snapshot uses 2s timeouts per call
	return refreshMsg(snapshot())
}

func tickRefresh() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return doSnapshot() })
}

func (m model) Init() tea.Cmd { return tickRefresh() }

func clampCursor(m *model) {
	if len(m.rows) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case refreshMsg:
		m.rows = []Node(msg)
		clampCursor(&m)
		return m, tickRefresh()
	case errMsg:
		m.err = string(msg)
		return m, tickRefresh()
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if i := rowAtY(m.rows, msg.Y); i >= 0 {
				m.cursor = i
			}
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				if n.Kind == "OUTPUT" {
					go sh("omarchy-audio-output-volume", "+5")
				} else {
					go sh("wpctl", "set-volume", n.ID, "5%+")
					m.rows[m.cursor].Pct += 5
				}
			}
		case tea.MouseButtonWheelDown:
			if i := rowAtY(m.rows, msg.Y); i >= 0 {
				m.cursor = i
			}
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				if n.Kind == "OUTPUT" {
					go sh("omarchy-audio-output-volume", "-5")
				} else {
					go sh("wpctl", "set-volume", n.ID, "5%-")
					m.rows[m.cursor].Pct -= 5
				}
			}
		case tea.MouseButtonRight:
			if msg.Action == tea.MouseActionPress {
				if i := rowAtY(m.rows, msg.Y); i >= 0 {
					m.cursor = i
					n := m.rows[i]
					switch n.Kind {
					case "OUTPUT":
						go sh("omarchy-audio-output-volume", "mute-toggle")
					case "INPUT":
						go sh("omarchy-audio-input-mute")
					default:
						go sh("wpctl", "set-mute", n.ID, "toggle")
					}
					m.flash = "toggled mute"
					return m, func() tea.Msg {
						time.Sleep(250 * time.Millisecond)
						return doSnapshot()
					}
				}
			}
		case tea.MouseButtonLeft:
			if msg.Action == tea.MouseActionPress {
				if i := rowAtY(m.rows, msg.Y); i >= 0 {
					m.cursor = i
					// click on the bar sets absolute volume (panel slider behavior)
					if msg.X >= barX0 && msg.X < barX0+barW && len(m.rows) > 0 {
						v := (msg.X - barX0 + 1) * 100 / barW
						if v < 0 {
							v = 0
						}
						if v > 100 {
							v = 100
						}
						n := m.rows[i]
						go sh("wpctl", "set-volume", n.ID, strconv.Itoa(v)+"%")
						m.rows[i].Pct = v
					}
				}
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		case "r":
			return m, func() tea.Msg { return doSnapshot() }
		case "m":
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				switch n.Kind {
				case "OUTPUT":
					go sh("omarchy-audio-output-volume", "mute-toggle")
					m.flash = "toggled mute"
				case "INPUT":
					go sh("omarchy-audio-input-mute")
					m.flash = "toggled mic"
				default:
					go sh("wpctl", "set-mute", n.ID, "toggle")
					m.flash = "toggled mute"
				}
			}
			return m, func() tea.Msg {
				time.Sleep(250 * time.Millisecond)
				return doSnapshot()
			}
		case "+", "=", "right", "l":
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				if n.Kind == "OUTPUT" {
					go sh("omarchy-audio-output-volume", "+5")
				} else {
					go sh("wpctl", "set-volume", n.ID, "5%+")
					m.rows[m.cursor].Pct += 5
				}
			}
		case "-", "_", "left", "h":
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				if n.Kind == "OUTPUT" {
					go sh("omarchy-audio-output-volume", "-5")
				} else {
					go sh("wpctl", "set-volume", n.ID, "5%-")
					m.rows[m.cursor].Pct -= 5
				}
			}
		case "enter":
			if len(m.rows) > 0 {
				n := m.rows[m.cursor]
				switch n.Kind {
				case "SINK":
					go sh("omarchy-audio-output-set-default", n.ID, n.Raw)
					m.flash = "default → " + shortName(n.Name, 30)
				case "SOURCE":
					go sh("omarchy-audio-input-set-default", n.ID, n.Raw)
					m.flash = "input → " + shortName(n.Name, 30)
				default:
					m.flash = "m mutes · enter selects outputs"
				}
				return m, func() tea.Msg {
					time.Sleep(300 * time.Millisecond)
					return doSnapshot()
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if len(m.rows) == 0 {
		return boxSt.Render(titleSt.Render("Ember") + "\n\n" + subSt.Render("scanning PipeWire…"))
	}
	var b strings.Builder
	b.WriteString(titleSt.Render(" Ember ") + " " + helpSt.Render("↑↓move · ←→vol · m mute · enter select · q quit") + "\n")
	b.WriteString(helpSt.Render("  click select · wheel volume · right-click mute") + "\n")
	flash := " "
	if m.flash != "" {
		flash = "  " + lipgloss.NewStyle().Foreground(cBlue).Render(m.flash)
	}
	b.WriteString(flash + "\n\n")
	lastKind := ""
	for i, r := range m.rows {
		if r.Kind != lastKind && r.Kind != "OUTPUT" && r.Kind != "INPUT" {
			switch r.Kind {
			case "SINK":
				b.WriteString("\n" + headSt.Render("OUTPUTS") + "\n")
			case "SOURCE":
				b.WriteString("\n" + headSt.Render("INPUTS") + "\n")
			case "STREAM":
				b.WriteString("\n" + headSt.Render("APPS") + "\n")
			}
			lastKind = r.Kind
		}
		name := shortName(r.Name, 28)
		if r.Kind == "OUTPUT" || r.Kind == "INPUT" {
			name = friendlyHero(r.Name)
		}
		glyph := glyphFor(r)
		vol := volBar(r.Pct, 14, r.Muted)
		pct := fmt.Sprintf("%3d%%", r.Pct)
		active := ""
		if r.Active {
			active = lipgloss.NewStyle().Foreground(cAccent).Render(" ●")
		}
		muteTag := ""
		if r.Muted {
			muteTag = " " + mutSt.Render("MUTED")
		}
		line := fmt.Sprintf("%s %-28s  %s %s%s%s", glyph, nameSt.Render(name), vol, pctSt.Render(pct), active, muteTag)
		if r.Kind == "OUTPUT" || r.Kind == "INPUT" {
			line = fmt.Sprintf("%s %-28s  %s %s%s%s", glyph, lipgloss.NewStyle().Foreground(cFg).Bold(true).Render(name), vol, pctSt.Render(pct), active, muteTag)
		}
		if i == m.cursor {
			b.WriteString(selSt.Render("▸ "+line) + "\n")
		} else {
			b.WriteString(dimSt.Render("  ") + line + "\n")
		}
	}
	if m.err != "" {
		b.WriteString("\n" + lipgloss.NewStyle().Foreground(cYellow).Render(m.err) + "\n")
	}
	w := m.width - 4
	if w < 52 {
		w = 64
	}
	if w > 78 {
		w = 78
	}
	_ = dimSt
	return boxSt.Width(w).Render(b.String())
}

func main() {
	// backend helpers live beside this binary; ensure they resolve anywhere
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		os.Setenv("PATH", dir+":"+os.Getenv("PATH"))
	}
	if len(os.Args) > 1 && os.Args[1] == "--dump" {
		for _, r := range snapshot() {
			fmt.Printf("%s\t%s\t%s\t%d\t%v\n", r.Kind, r.ID, r.Name, r.Pct, r.Muted)
		}
		return
	}
	m := model{rows: snapshot(), cursor: 0}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ember:", err)
		os.Exit(1)
	}
}
