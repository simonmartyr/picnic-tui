package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/simonmartyr/picnic"
	qrcode "github.com/skip2/go-qrcode"
	"time"
)

type CheckoutPage struct {
	Instructions *tview.TextView
	PaymentLink  *tview.TextView
	Basket       *tview.List
	Flex         *tview.Flex
	Checkout     *picnic.Checkout
	Payment      *picnic.Payment
}

func ShowCheckoutPage() {
	checkoutPage := newCheckoutPage()
	App.Pages.AddAndSwitchToPage(CheckoutId, checkoutPage.Flex, true)
	go checkoutPage.renderBasket()
}

func newCheckoutPage() *CheckoutPage {
	checkoutPage := &CheckoutPage{}

	checkoutPage.Instructions = createInstructions()
	checkoutPage.Basket = createCheckoutBasket()
	checkoutPage.PaymentLink = createPaymentLink()

	checkout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(checkoutPage.Instructions, 3, 1, false).
		AddItem(checkoutPage.Basket, 0, 1, false)

	checkoutPage.Flex = tview.NewFlex().
		AddItem(checkout, 0, 2, false).
		AddItem(checkoutPage.PaymentLink, 0, 1, false)
	checkoutPage.Flex.SetInputCapture(checkoutPage.setupHotKeys)
	return checkoutPage
}

func createInstructions() *tview.TextView {
	instructions := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	instructions.SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitle(" Checkout ")
	instructions.SetText("To begin checking out press 'c' ")
	return instructions
}

func createPaymentLink() *tview.TextView {
	paymentLink := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	paymentLink.SetTitle(" Payment Link ").
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignCenter)
	paymentLink.SetText("Awaiting Payment Link...")
	return paymentLink
}

func (c *CheckoutPage) renderPaymentLink(url string) {
	q, _ := qrcode.New(url, qrcode.Highest)
	c.PaymentLink.SetText(q.ToSmallString(true))
	c.Instructions.SetText("Scan QR code to complete or ESC to cancel")
}

func createCheckoutBasket() *tview.List {
	basket := tview.NewList().
		SetSelectedFocusOnly(true)
	basket.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(" Basket ").
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft)
	return basket
}

func (c *CheckoutPage) setupHotKeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEscape {
		if c.Checkout != nil {
			App.Client.CancelCheckout(c.Checkout.OrderId)
		}
		SwitchToMainPage()
		return event
	}
	if event.Rune() == 'c' {
		c.Instructions.SetText("Verifying cart...")
		go c.beginCheckout("")
	}
	return event
}

func (c *CheckoutPage) beginCheckout(verification string) {
	var result *picnic.Checkout
	var err *picnic.CheckoutError
	if verification != "" {
		result, err = App.Client.CheckoutWithResolveKey(App.Cart.Mts, verification)
	} else {
		result, err = App.Client.StartCheckout(App.Cart.Mts)
	}
	if err != nil {
		if err.Blocking {
			App.Tview.QueueUpdateDraw(func() {
				c.PresentBlockingCheckoutError(err)
			})
		} else {
			App.Tview.QueueUpdateDraw(func() {
				c.PresentNonBlockingError(err)
			})
		}
		return
	}
	c.Checkout = result
	c.Instructions.SetText("[green]Cart verified")
	go c.beginPayment()
}

func (c *CheckoutPage) PresentBlockingCheckoutError(error *picnic.CheckoutError) {
	model := tview.NewModal().
		SetText(fmt.Sprintf("%s %s", error.Title, error.Message)).
		SetBackgroundColor(primaryColor).
		AddButtons([]string{"ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			App.Pages.RemovePage(CheckoutErrorId)
		})
	App.Pages.AddPage(CheckoutErrorId, model, true, true)
}

func (c *CheckoutPage) PresentNonBlockingError(error *picnic.CheckoutError) {
	model := tview.NewModal().
		SetText(fmt.Sprintf("%s %s", error.Title, error.Message)).
		SetBackgroundColor(primaryColor).
		AddButtons([]string{"Ja", "Nee"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			App.Pages.RemovePage(AgeVerificationId)
			if buttonLabel == "Ja" {
				go c.beginCheckout(error.ResolveKey)
			}
		})
	App.Pages.AddPage(AgeVerificationId, model, true, true)
}

func (c *CheckoutPage) beginPayment() {
	App.Tview.QueueUpdateDraw(func() {
		c.Instructions.SetText("Requesting payment url...")
	})
	res, err := App.Client.InitiatePayment(c.Checkout.OrderId)
	if err != nil {
		ShowErrorPage(err.Error())
		return
	}
	App.Tview.QueueUpdateDraw(func() {
		c.renderPaymentLink(res.IssuerAuthenticationUrl)
	})
}

func (c *CheckoutPage) renderBasket() {
	cart, cartErr := App.Client.GetCart()
	if cartErr != nil {
		ShowErrorPage(cartErr.Error())
		return
	}
	App.Cart = cart
	App.Tview.QueueUpdateDraw(func() {
		ind := c.Basket.GetCurrentItem()
		c.Basket.Clear()
		for _, orderLines := range App.Cart.Items {
			for _, article := range orderLines.Items {
				var prefix = ""
				switch {
				case !article.IsAvailable():
					prefix = "[red]"
				case orderLines.IsOnPromotion():
					prefix = "[green]"
				}
				basketItemText := fmt.Sprintf("%s%d: %s - %s", prefix, article.Quantity(), article.Name, FormatIntToPrice(orderLines.PriceIncludingPromotions()))
				c.Basket.AddItem(basketItemText, article.Id, 0, nil)
			}
		}
		var deliveryTime = ""
		for _, slot := range App.Cart.DeliverySlots {
			if !slot.Selected {
				continue
			}
			start, _ := time.Parse(time.RFC3339, slot.WindowStart)
			end, _ := time.Parse(time.RFC3339, slot.WindowEnd)
			deliveryTime = start.Format("Mon, 02 Jan") + " " + start.Format("15:04") + "-" + end.Format("15:04")
			break
		}
		if deliveryTime == "" {
			c.Basket.SetTitle(" (B)asket (" + FormatIntToPrice(App.Cart.TotalPrice) + ")")
			return
		}
		c.Basket.SetTitle(" (B)asket (" + FormatIntToPrice(App.Cart.TotalPrice) + " | " + deliveryTime + ")")
		c.Basket.SetCurrentItem(ind)
	})
}
