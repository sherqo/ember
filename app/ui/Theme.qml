pragma Singleton
import QtQuick

QtObject {
  // Pink Cat Boo — Ember's whole theme lives here.
  readonly property color background: "#202330"
  readonly property color foreground: "#FFF0F5"
  readonly property color accent: "#FF4C7A"
  readonly property color urgent: "#FEC831"
  readonly property color muted: "#565970"
  readonly property color green: "#3BC089"
  readonly property color blue: "#6767CE"
  readonly property color magenta: "#C77DFF"
  readonly property color cyan: "#4CC9C0"

  readonly property string fontFamily: "CaskaydiaMono Nerd Font"
  readonly property int fontSmall: 11
  readonly property int fontBody: 13
  readonly property int fontTitle: 16
  readonly property int fontHero: 28

  readonly property int radius: 8
  readonly property int padding: 16
  readonly property int spacing: 10
}
