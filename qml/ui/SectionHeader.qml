import QtQuick
import "ui" as Ui

Text {
  property string title: ""
  text: title.toUpperCase()
  color: Ui.Theme.muted
  font.family: Ui.Theme.fontFamily
  font.pixelSize: Ui.Theme.fontSmall
  font.bold: true
}
