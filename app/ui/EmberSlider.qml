import QtQuick
import "ui" as Ui

// Drag + wheel slider with % readout. Binds to any 0..1 volume + mute.
Item {
  id: root
  property real value: 0
  property bool muted: false
  property string icon: ""
  signal setVolume(real v)
  signal toggleMute()

  implicitHeight: 44
  clip: true

  function shift(d) {
    root.setVolume(Math.max(0, Math.min(1.5, root.value + d)))
  }

  Row {
    anchors.fill: parent
    spacing: Ui.Theme.spacing

    Text {
      id: glyph
      width: 30
      verticalAlignment: Text.AlignVCenter
      height: parent.height
      text: root.icon
      color: root.muted ? Ui.Theme.muted : Ui.Theme.foreground
      font.family: Ui.Theme.fontFamily
      font.pixelSize: Ui.Theme.fontHero - 6
      MouseArea {
        anchors.fill: parent
        cursorShape: Qt.PointingHandCursor
        onClicked: root.toggleMute()
      }
    }

    Column {
      width: parent.width - glyph.width - pct.width - parent.spacing * 2
      anchors.verticalCenter: parent.verticalCenter
      spacing: 2

      Rectangle {
        width: parent.width
        height: 8
        radius: 4
        color: Ui.Theme.muted

        Rectangle {
          width: parent.width * Math.min(1, root.muted ? 0 : root.value)
          height: parent.height
          radius: 4
          color: root.muted ? Ui.Theme.muted : Ui.Theme.accent
        }

        MouseArea {
          anchors.fill: parent
          cursorShape: Qt.PointingHandCursor
          onPressed: mouse => root.setVolume(mouse.x / width * 1.0)
          onPositionChanged: mouse => { if (mouse.buttons) root.setVolume(mouse.x / width * 1.0) }
          onWheel: wheel => root.shift(wheel.angleDelta.y > 0 ? 0.05 : -0.05)
        }
      }
    }

    Text {
      id: pct
      width: 52
      height: parent.height
      verticalAlignment: Text.AlignVCenter
      horizontalAlignment: Text.AlignRight
      text: Math.round(root.value * 100) + "%"
      color: root.muted ? Ui.Theme.muted : Ui.Theme.foreground
      font.family: Ui.Theme.fontFamily
      font.pixelSize: Ui.Theme.fontBody
    }
  }
}
