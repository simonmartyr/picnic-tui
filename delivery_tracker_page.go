package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/simonmartyr/picnic-api"
	"time"
)

type DeliveryTrackPage struct {
	Deliveries       *tview.List
	DeliveryArticles *tview.List
	Flex             *tview.Flex
}

func ShowDeliveryTrackPage() {
	deliveryTrackPage := newDeliveryTrackPage()
	App.Pages.AddAndSwitchToPage(DeliveryTrackerId, deliveryTrackPage.Flex, true)
}

func newDeliveryTrackPage() *DeliveryTrackPage {
	deliveryTrackPage := &DeliveryTrackPage{}

	deliveryTrackPage.Flex = createDeliveryTrackFlex()
	deliveryTrackPage.Deliveries = createDeliveryList()
	deliveryTrackPage.DeliveryArticles = createDeliveryDetails()
	deliveryTrackPage.Flex.SetInputCapture(deliveryTrackPage.setupHotkeys)
	deliveryTrackPage.Flex.AddItem(deliveryTrackPage.Deliveries, 0, 5, true)
	deliveryTrackPage.Flex.AddItem(deliveryTrackPage.DeliveryArticles, 0, 5, false)

	deliveryTrackPage.Deliveries.SetSelectedFunc(deliveryTrackPage.setupDeliverySelected)

	go deliveryTrackPage.renderDeliveries()

	return deliveryTrackPage
}

func createDeliveryTrackFlex() *tview.Flex {
	deliveryFlex := tview.NewFlex()

	deliveryFlex.
		SetBorder(false).
		SetTitle(" Delivery Tracker ").
		SetBorderColor(tcell.ColorRed).
		SetTitleAlign(tview.AlignLeft)

	return deliveryFlex
}

func (d *DeliveryTrackPage) setupHotkeys(event *tcell.EventKey) *tcell.EventKey {
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
	if event.Key() == tcell.KeyEscape {
		SwitchToMainPage()
	}
	return event
}

func createDeliveryDetails() *tview.List {
	articleList := tview.NewList().
		SetSelectedTextColor(HighlightColor).
		SetSelectedFocusOnly(true)
	articleList.ShowSecondaryText(false).
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitle(" Article Details ").
		SetTitleAlign(tview.AlignLeft)
	return articleList
}

func createDeliveryList() *tview.List {
	deliveryList := tview.NewList().
		SetSelectedTextColor(HighlightColor).
		SetSelectedFocusOnly(true)
	deliveryList.ShowSecondaryText(false).
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitle(" Deliveries ").
		SetTitleAlign(tview.AlignLeft)
	return deliveryList
}

func (d *DeliveryTrackPage) renderDeliveries() {
	data, deliveryErr := App.Client.GetDeliveries([]picnic.DeliveryStatus{picnic.CURRENT, picnic.COMPLETED})
	if deliveryErr != nil {
		ShowErrorPage(deliveryErr.Error())
		return
	}
	App.Tview.QueueUpdateDraw(func() {
		d.Deliveries.Clear()
		for _, delivery := range *data {
			date, err := time.Parse(time.RFC3339, delivery.CreationTime)
			if err != nil {
				panic(err.Error())
			}
			d.Deliveries.AddItem(fmt.Sprintf("%s - %s - %s",
				string(delivery.Status), date.Format("Mon, 02 Jan"), calculateTotal(delivery.Orders)),
				delivery.DeliveryId, 0, nil)
		}
	})
}

func (d *DeliveryTrackPage) setupDeliverySelected(i int, s string, s2 string, r rune) {
	go d.renderArticleDetails(s2)
}

func (d *DeliveryTrackPage) renderArticleDetails(deliveryId string) {
	d.DeliveryArticles.Clear()
	qDelivery, err := App.Client.GetDelivery(deliveryId)
	if err != nil {
		App.Tview.QueueUpdateDraw(func() {
			ShowErrorPage(err.Error())
		})
		return
	}
	if err != nil {
		panic(err.Error())
	}
	App.Tview.QueueUpdateDraw(func() {
		for _, order := range qDelivery.Orders {
			for _, item := range order.Items {
				for _, article := range item.Items {
					basketItemText := fmt.Sprintf("%d: %s - %s", article.Quantity(), article.Name, FormatIntToPrice(item.PriceIncludingPromotions()))
					d.DeliveryArticles.AddItem(basketItemText, article.Id, 0, nil)
				}
			}
		}
	})
}

func calculateTotal(orders []picnic.Order) string {
	var total = 0
	for _, order := range orders {
		total += order.TotalPrice
	}
	return FormatIntToPrice(total)
}
