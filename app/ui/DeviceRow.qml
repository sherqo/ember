import QtQuick
import "ui" as Ui

// Generic row: glyph + label + right-side text. Click = primary action,
// right-click = secondary (mute) when provided.
Item {
  id: root
  property string glyph: ""
  property string label: ""
  property string detail: ""
  property bool active: false
  property bool dimmed: false
  signal clicked()
  signal secondary()

  implicitHeight: 34

  Rectangle {
    anchors.fill: parent
    radius: 6
    color: Ui.Theme.accent
    opacity: 0.14
    visible: root.active
  }

  Row {
    anchors.fill: parent
    anchors.leftMargin: 8
    anchors.rightMargin: 8
    spacing: Ui.Theme.spacing

    Text {
      width: 26
      height: parent.height
      verticalAlignment: Text.AlignVCenter
      text: root.glyph
      color: root.dimmed ? Ui.Theme.muted : Ui.Theme.foreground
      font.family: Ui.Theme.fontFamily
      font.pixelSize: Ui.Theme.fontBody + 2
    }

    Text {
      width: parent.width - 26 - detailLabel.width - parent.spacing * 2
      height: parent.height
      verticalAlignment: Text.AlignVCenter
      elide: Text.ElideRight
      text: root.label
      color: root.dimmed ? Ui.Theme.muted : Ui.Theme.foreground
      font.family: Ui.Theme.fontFamily
      font.pixelSize: Ui.Theme.fontBody
    }

    Text {
      id: detailLabel
      height: parent.height
      verticalAlignment: Text.AlignVCenter
      text: root.detail
      color: Ui.Theme.muted
      font.family: Ui.Theme.fontFamily
      font.pixelSize: Ui.Theme.fontSmall
    }
  }

  MouseArea {
    anchors.fill: parent
    hoverEnabled: true
    cursorShape: Qt.PointingHandCursor
    acceptedButtons: Qt.LeftButton | Qt.RightButton
    onClicked: mouse => {
      if (mouse.button === Qt.RightButton) root.secondary()
      else root.clicked()
    }
  }
}
