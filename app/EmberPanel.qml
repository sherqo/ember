import QtQuick
import Quickshell
import Quickshell.Services.Pipewire
import Quickshell.Io
import "ui" as Ui
import "ui" as Ui

// Ember's own mixer panel. Simple Column, no cursors, no popouts.
Column {
  id: root
  spacing: 10

  property var sink: Pipewire.defaultAudioSink
  property var source: Pipewire.defaultAudioSource
  property real outVol: sink && sink.audio ? sink.audio.volume : 0
  property bool outMuted: sink && sink.audio ? sink.audio.muted : false
  property real inVol: source && source.audio ? source.audio.volume : 0
  property bool inMuted: source && source.audio ? source.audio.muted : false

  function nodes() { return Pipewire.nodes ? Pipewire.nodes.values : [] }
  function sinks() { return nodes().filter(n => n && n.isSink && !n.isStream && n.audio) }
  function sources() {
    return nodes().filter(n => n && !n.isSink && !n.isStream && n.audio && (n.name || "").indexOf("quickshell") !== 0)
  }
  function streams() { return nodes().filter(n => n && n.isStream && n.audio) }

  function label(n) {
    if (!n) return "—"
    var p = n.properties || {}
    var t = n.description || p["node.description"] || n.nickname || n.name || "Unknown"
    return String(t).length > 34 ? String(t).slice(0, 33) + "…" : String(t)
  }

  Process { id: runProc }
  function shQuote(s) { return "'" + String(s == null ? "" : s).replace(/'/g, "'\\''") + "'" }
  function setDefault(kind, node) {
    var cmd = kind === "SOURCE" ? "omarchy-audio-input-set-default" : "omarchy-audio-output-set-default"
    runProc.command = ["sh", "-c", "export PATH=\"$HOME/ember/bin:$PATH\"; " + cmd + " " + shQuote(node.id) + " " + shQuote(node.name)]
    runProc.running = true
  }

  Ui.EmberSlider {
    width: parent.width
    value: root.outVol
    muted: root.outMuted
    icon: root.outMuted ? "" : ""
    onSetVolume: v => { if (root.sink && root.sink.audio) root.sink.audio.volume = Math.max(0, Math.min(1.5, v)) }
    onToggleMute: { if (root.sink && root.sink.audio) root.sink.audio.muted = !root.sink.audio.muted }
  }

  Ui.EmberSlider {
    width: parent.width
    value: root.inVol
    muted: root.inMuted
    icon: root.inMuted ? "" : ""
    onSetVolume: v => { if (root.source && root.source.audio) root.source.audio.volume = Math.max(0, Math.min(1, v)) }
    onToggleMute: { if (root.source && root.source.audio) root.source.audio.muted = !root.source.audio.muted }
  }

  Ui.SectionHeader { title: "Outputs"; visible: sinkRep.count > 0 }
  Repeater {
    id: sinkRep
    model: []
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: ""
      label: root.label(modelData)
      detail: Math.round((modelData.audio ? modelData.audio.volume : 0) * 100) + "%"
      active: root.sink && modelData.id === root.sink.id
      dimmed: modelData.audio && modelData.audio.muted
      onClicked: root.setDefault("SINK", modelData)
      onSecondary: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }

  Ui.SectionHeader { title: "Inputs"; visible: sourceRep.count > 0 }
  Repeater {
    id: sourceRep
    model: []
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: ""
      label: root.label(modelData)
      active: root.source && modelData.id === root.source.id
      onClicked: root.setDefault("SOURCE", modelData)
      onSecondary: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }

  Ui.SectionHeader { title: "Apps"; visible: streamRep.count > 0 }
  Repeater {
    id: streamRep
    model: []
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: modelData.audio && modelData.audio.muted ? "" : ""
      label: root.label(modelData)
      detail: Math.round((modelData.audio ? modelData.audio.volume : 0) * 100) + "%"
      dimmed: modelData.audio && modelData.audio.muted
      onClicked: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }

  // PipeWire enumeration lands late; re-resolve section models on a tick.
  function refresh() {
    sinkRep.model = root.sinks()
    sourceRep.model = root.sources()
    streamRep.model = root.streams()
  }
  Timer {
    interval: 2000
    running: true
    repeat: true
    triggeredOnStart: true
    onTriggered: root.refresh()
  }
}
