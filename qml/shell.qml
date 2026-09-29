import QtQuick
import Quickshell
import Quickshell.Wayland
import "ui" as Ui
import "audio" as Audio

// Ember — warm little audio mixer. One window, quits on close.
ShellRoot {
  id: shell

  PanelWindow {
    id: win
    color: "transparent"
    exclusionMode: ExclusionMode.Ignore
    WlrLayershell.namespace: "ember"
    WlrLayershell.layer: WlrLayer.Overlay
    WlrLayershell.keyboardFocus: WlrKeyboardFocus.Exclusive

    anchors.top: true
    anchors.right: true
    margins.top: 34
    margins.right: 12
    implicitWidth: 380
    implicitHeight: 540

    FocusScope {
      anchors.fill: parent
      focus: true
      Keys.onEscapePressed: Qt.quit()

      Rectangle {
        anchors.fill: parent
        radius: Ui.Theme.radius
        color: Ui.Theme.background
        border.color: Ui.Theme.muted
        border.width: 1
      }

      Audio.Panel {
        anchors.fill: parent
        anchors.margins: Ui.Theme.padding
      }
    }
  }
}
