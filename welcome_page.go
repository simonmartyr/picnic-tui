package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type WelcomePage struct {
	Text *tview.TextView
}

func ShowWelcomePage() {
	welcomePage := newWelcomePage()
	if welcomePage == nil {
		return
	}
	App.Pages.AddAndSwitchToPage(WelcomeId, welcomePage.Text, true)
}

func newWelcomePage() *WelcomePage {
	welcome := &WelcomePage{}
	welcomeText := tview.NewTextView()
	welcomeText.SetTitle(" Welcome! ").
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft)
	welcomeText.SetText(fmt.Sprintf("%s! \n Welcome to Picnic-tui \n\n\n[red]%s", TimeOfDay(), Logo))
	welcomeText.SetTextAlign(tview.AlignCenter).SetRegions(true)
	welcomeText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		App.Pages.RemovePage(WelcomeId)
		ShowMainPage()
		return event
	})
	welcome.Text = welcomeText
	go welcome.displayUser()
	return welcome
}

func (w *WelcomePage) displayUser() {
	user, err := App.Client.GetUser()
	if err != nil {
		ShowErrorPage(err.Error())
	}
	App.Tview.QueueUpdateDraw(func() {
		w.Text.SetText(fmt.Sprintf("%s %s! \n Welcome to Picnic-tui \n\n\n[red]%s\n\nversion: %s", TimeOfDay(), user.Firstname, Logo, Version))
	})
}
