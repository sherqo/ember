import QtQuick
import Quickshell
import Quickshell.Wayland

// Ember host: EmberPanel fills this window. Esc quits. Zero idle RAM.
ShellRoot {
  id: shell

  PanelWindow {
    id: win
    color: "transparent"
    exclusionMode: ExclusionMode.Ignore
    WlrLayershell.namespace: "ember"
    WlrLayershell.layer: WlrLayer.Overlay

    anchors.top: true
    anchors.right: true
    margins.top: 34
    margins.right: 12
    implicitWidth: 400
    implicitHeight: 600

    FocusScope {
      anchors.fill: parent
      focus: true
      Keys.onEscapePressed: Qt.quit()

      Rectangle {
        anchors.fill: parent
        radius: 10
        color: "#202330"
        border.color: "#565970"
        border.width: 1
      }

      EmberPanel {
        anchors.fill: parent
        anchors.margins: 14
      }
    }
  }
}
