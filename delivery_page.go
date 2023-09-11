package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type DeliveryPage struct {
	Flex *tview.Flex
}

func ShowDeliveryPage() {
	deliveryPage := newDeliveryPage()
	App.Pages.AddAndSwitchToPage(DeliveryId, deliveryPage.Flex, true)
	App.Tview.SetFocus(deliveryPage.Flex.GetItem(0))
}

func newDeliveryPage() *DeliveryPage {
	deliveryPage := &DeliveryPage{}

	deliveryPage.Flex = createDeliveryFlex()
	deliveryPage.renderDeliverySlotsPage()
	deliveryPage.Flex.SetInputCapture(deliveryPage.setupHotkeys)

	return deliveryPage
}

func createDeliveryFlex() *tview.Flex {
	deliveryFlex := tview.NewFlex()

	deliveryFlex.
		SetBorder(true).
		SetTitle(" Choose Delivery Time (💚 green choice for your neighbourhood) ").
		SetBorderColor(tcell.ColorRed).
		SetTitleAlign(tview.AlignLeft)

	return deliveryFlex
}

func (d *DeliveryPage) setupHotkeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyRight {
		for i := 0; i < d.Flex.GetItemCount(); i++ {
			if !d.Flex.GetItem(i).HasFocus() {
				continue
			}
			i = i + 1
			i = i % d.Flex.GetItemCount()
			App.Tview.SetFocus(d.Flex.GetItem(i))
			return &tcell.EventKey{}
		}
	}
	if event.Key() == tcell.KeyLeft {
		for i := d.Flex.GetItemCount() - 1; i >= 0; i-- {
			if !d.Flex.GetItem(i).HasFocus() {
				continue
			}
			i = i - 1
			if i < 0 {
				App.Tview.SetFocus(d.Flex.GetItem(d.Flex.GetItemCount() - 1))
			} else {
				App.Tview.SetFocus(d.Flex.GetItem(i))
			}
			return &tcell.EventKey{}
		}
	}
	if event.Key() == tcell.KeyEscape {
		SwitchToMainPage()
	}
	return event
}

func (d *DeliveryPage) renderDeliverySlotsPage() {
	var previousDate time.Time
	var dayRef *tview.List
	d.Flex.Clear()
	for i, slot := range App.Cart.DeliverySlots {
		switch slot.IsAvailable {
		case true:
			var postFix = ""
			var preFix = ""
			if slot.Icon.PmlVersion != "" {
				postFix = "💚"
			}
			if slot.Selected {
				preFix = "[:red]"
			}
			date, err := time.Parse(time.RFC3339, slot.WindowStart)
			if err != nil {
				panic(err.Error())
			}
			endDate, endErr := time.Parse(time.RFC3339, slot.WindowEnd)
			if endErr != nil {
				panic(err.Error())
			}
			day := date.Truncate(24 * time.Hour)
			if i == 0 || !day.Equal(previousDate) {
				dayRef = tview.NewList().
					ShowSecondaryText(false).
					SetHighlightFullLine(true)
				dayRef.SetTitle(" " + date.Format("Mon, 02 Jan") + " ").
					SetBorderColor(tcell.ColorRed).
					SetBorder(true)

				dayRef.SetSelectedFocusOnly(true)
				dayRef.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
					App.Cart, err = App.Client.SetDeliverySlot(s2)
					if err != nil {
						ShowErrorPage(err.Error())
						return
					}
					SwitchToMainPage()
				})
				d.Flex.AddItem(dayRef, 0, 1, false)
				previousDate = day
			}
			dayRef.AddItem(preFix+date.Format("15:04")+"-"+endDate.Format("15:04")+" "+postFix, slot.SlotId, 0, nil)
		}
	}
}
