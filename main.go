package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("MyApp")

	label := widget.NewLabel("Test1")
	label2 := widget.NewLabel("Test2")

	w.SetContent(container.NewVBox(
		label,
		label2,
	))

	w.ShowAndRun()
}
