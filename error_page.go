package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ErrorPage struct {
	Text *tview.TextView
}

func ShowErrorPage(error string) {
	errorPage := newErrorPage(error)
	App.Pages.AddAndSwitchToPage(ErrorId, errorPage.Text, true)
}

func newErrorPage(errorMessage string) *ErrorPage {
	errorPage := &ErrorPage{}

	errorText := tview.NewTextView()
	errorText.SetTitle(" Error 🫣 ").
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft)
	errorText.SetText("An Error Occurred \n\n\n" + errorMessage + "\n press any key to quit")
	errorText.SetTextAlign(tview.AlignCenter).SetRegions(true)
	errorText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		App.Tview.Stop()
		return event
	})
	errorPage.Text = errorText
	return errorPage
}
