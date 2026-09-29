import QtQuick
import QtQuick.Layouts
import Quickshell
import Quickshell.Services.Pipewire
import Quickshell.Io
import "Model.js" as Model
import "../ui" as Ui

// Ember's mixer: hero output + mic + sinks + sources + per-app streams.
Column {
  id: root
  spacing: Ui.Theme.spacing

  property var sink: Pipewire.defaultAudioSink
  property var source: Pipewire.defaultAudioSource
  property real outVol: sink && sink.audio ? sink.audio.volume : 0
  property bool outMuted: sink && sink.audio ? sink.audio.muted : false
  property real inVol: source && source.audio ? source.audio.volume : 0
  property bool inMuted: source && source.audio ? source.audio.muted : false

  function allNodes() { return Pipewire.nodes ? Pipewire.nodes.values : [] }
  function sinks() { return allNodes().filter(n => n && n.isSink && !n.isStream) }
  function sources() { return allNodes().filter(n => n && !n.isSink && !n.isStream && n.audio) }
  function streams() { return allNodes().filter(n => n && n.isStream && n.audio) }

  Process { id: runProc }

  function shQuote(s) { return "'" + String(s == null ? "" : s).replace(/'/g, "'\\''") + "'" }
  function setDefaultSink(node) {
    runProc.command = ["sh", "-c", "export PATH=\"$HOME/ember/bin:$PATH\"; omarchy-audio-output-set-default " + shQuote(node.id) + " " + shQuote(node.name)]
    runProc.running = true
  }
  function setDefaultSource(node) {
    runProc.command = ["sh", "-c", "export PATH=\"$HOME/ember/bin:$PATH\"; omarchy-audio-input-set-default " + shQuote(node.id) + " " + shQuote(node.name)]
    runProc.running = true
  }

  PwObjectTracker {
    objects: [...root.sinks(), ...root.sources(), ...root.streams(), root.sink, root.source]
  }

  Ui.SectionHeader { title: "Output · " + (root.sink ? Model.nodeLabel(root.sink) : "—") }

  Ui.EmberSlider {
    width: parent.width
    value: root.outVol
    muted: root.outMuted
    icon: root.outMuted ? "" : ""
    onSetVolume: v => { if (root.sink && root.sink.audio) root.sink.audio.volume = Math.max(0, Math.min(1.5, v)) }
    onToggleMute: { if (root.sink && root.sink.audio) root.sink.audio.muted = !root.sink.audio.muted }
  }

  Ui.SectionHeader { title: "Input" }

  Ui.EmberSlider {
    width: parent.width
    value: root.inVol
    muted: root.inMuted
    icon: root.inMuted ? "" : ""
    onSetVolume: v => { if (root.source && root.source.audio) root.source.audio.volume = Math.max(0, Math.min(1, v)) }
    onToggleMute: { if (root.source && root.source.audio) root.source.audio.muted = !root.source.audio.muted }
  }

  Ui.SectionHeader { title: "Outputs"; visible: sinkRep.count > 1 }

  Repeater {
    id: sinkRep
    model: root.sinks()
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: Model.sinkGlyph(modelData)
      label: Model.nodeLabel(modelData)
      detail: Math.round((modelData.audio ? modelData.audio.volume : 0) * 100) + "%"
      active: root.sink && modelData.id === root.sink.id
      dimmed: modelData.audio && modelData.audio.muted
      onClicked: root.setDefaultSink(modelData)
      onSecondary: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }

  Ui.SectionHeader { title: "Inputs"; visible: sourceRep.count > 1 }

  Repeater {
    id: sourceRep
    model: root.sources()
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: Model.sourceGlyph(modelData)
      label: Model.nodeLabel(modelData)
      active: root.source && modelData.id === root.source.id
      onClicked: root.setDefaultSource(modelData)
      onSecondary: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }

  Ui.SectionHeader { title: "Apps"; visible: streamRep.count > 0 }

  Repeater {
    id: streamRep
    model: root.streams()
    delegate: Ui.DeviceRow {
      required property var modelData
      width: root.width
      glyph: modelData.audio && modelData.audio.muted ? "" : ""
      label: Model.streamLabel(modelData)
      detail: Math.round((modelData.audio ? modelData.audio.volume : 0) * 100) + "%"
      dimmed: modelData.audio && modelData.audio.muted
      onClicked: { if (modelData.audio) modelData.audio.muted = !modelData.audio.muted }
    }
  }
}
