package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	picnic "github.com/simonmartyr/picnic-api"
	"strings"
	"time"
)

type MainPage struct {
	Search       *tview.InputField
	Basket       *tview.List
	Articles     *tview.List
	ArticleInfo  *tview.TextView
	ArticleImage *tview.Image
	Flex         *tview.Flex
	ShowBundles  bool
}

func ShowMainPage() {
	mainPage := newMainPage()
	if App.Cart == nil {
		go mainPage.renderBasket()
	}
	App.Main = mainPage
	App.Pages.AddAndSwitchToPage(MainId, mainPage.Flex, true)
}

func SwitchToMainPage() {
	if !App.Pages.HasPage(MainId) {
		ShowMainPage()
		return
	}
	App.Pages.SwitchToPage(MainId)
	go App.Main.renderBasket()
}

func newMainPage() *MainPage {
	mainPage := &MainPage{
		ShowBundles: true,
	}

	mainPage.Search = createSearch()
	mainPage.Basket = createBasket()
	mainPage.Articles = createArticleList()
	mainPage.ArticleInfo = createArticleInfo()
	mainPage.ArticleImage = createArticleImage()

	mainPage.Search.SetDoneFunc(mainPage.performSearch)
	mainPage.Basket.SetInputCapture(mainPage.setupBasketHotkeys)
	mainPage.Basket.SetSelectedFunc(mainPage.setupBasketSelected)
	mainPage.Articles.SetInputCapture(mainPage.setupArticleHotkeys)
	mainPage.Articles.SetSelectedFunc(mainPage.setupArticleSelected)

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainPage.Search, 3, 1, false).
		AddItem(mainPage.ArticleImage, 0, 4, false).
		AddItem(mainPage.ArticleInfo, 0, 4, false)

	flex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, false).
		AddItem(mainPage.Articles, 0, 2, false).
		AddItem(mainPage.Basket, 0, 1, false)

	mainPage.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flex, 0, 5, false)

	mainPage.Flex.SetInputCapture(mainPage.setupHotKeys)
	return mainPage
}

func createSearch() *tview.InputField {
	searchInput := tview.NewInputField().
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	searchInput.SetFieldWidth(30).
		SetAcceptanceFunc(tview.InputFieldMaxLength(30))

	searchInput.SetTitle(" (S)earch ").
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return searchInput
}

func createArticleList() *tview.List {
	articleList := tview.NewList().
		SetSelectedTextColor(HighlightColor).
		SetSelectedFocusOnly(true)
	articleList.ShowSecondaryText(false).
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitle(" (A)rticles ").
		SetTitleAlign(tview.AlignLeft)
	return articleList
}

func createBasket() *tview.List {
	basket := tview.NewList().
		SetSelectedTextColor(HighlightColor).
		SetSelectedFocusOnly(true)
	basket.ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(" (B)asket ").
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft)
	return basket
}

func createArticleImage() *tview.Image {
	return tview.NewImage().
		SetColors(tview.TrueColor)
}

func createArticleInfo() *tview.TextView {
	articleInfo := tview.NewTextView()
	articleInfo.SetTitle(" Article Info ").
		SetBorder(true).
		SetBorderColor(primaryColor).
		SetTitleAlign(tview.AlignLeft)
	articleInfo.SetText(" To see more information about an article press 'f' whilst highlighting in the basket or search list")
	articleInfo.SetRegions(true).
		SetWordWrap(true)
	return articleInfo
}

func (m *MainPage) performSearch(key tcell.Key) {
	if key == tcell.KeyEnter {
		searchTerm := m.Search.GetText()
		if len(searchTerm) == 0 {
			return
		}
		m.Articles.SetTitle(" (A)rticles [Searching...] ")
		go m.renderSearch(searchTerm)
		App.Tview.SetFocus(m.Articles)
	}
}

func (m *MainPage) renderSearch(term string) {
	data, searchErr := App.Client.SearchArticles(term)
	if searchErr != nil {
		ShowErrorPage(searchErr.Error())
		return
	}

	App.Tview.QueueUpdateDraw(func() {
		m.Articles.Clear()
		m.Articles.SetTitle(fmt.Sprintf("(A)rticles ([orange]term:[white] %s | [orange]bonus:[white] %t)", term, m.ShowBundles))

		var previousItem = ""
		for _, art := range data {
			if previousItem == art.Name && strings.Contains(art.UnitQuantity, "x") && !m.ShowBundles {
				continue
			}

			itemText := art.Name
			if art.IsOnPromotion() || (previousItem == art.Name && strings.Contains(art.UnitQuantity, "x")) {
				itemText += " [green]" + FormatIntToPrice(art.PriceIncludingPromotions())
				if previousItem == art.Name {
					itemText += " [white:red] B "
				}
			} else {
				itemText += " " + FormatIntToPrice(art.DisplayPrice)
			}
			previousItem = art.Name

			m.Articles.AddItem(itemText, art.Id, 0, nil)
		}
	})
}

func (m *MainPage) setupHotKeys(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEscape {
		App.Tview.SetFocus(m.Flex)
		return event
	}
	if event.Key() == tcell.KeyCtrlR {
		go m.renderBasket()
	}
	if event.Key() == tcell.KeyCtrlS {
		m.Search.SetText("")
		App.Tview.SetFocus(m.Search)
		return &tcell.EventKey{}
	}
	switch r := event.Rune(); {
	case m.Search.HasFocus():
		return event
	case r == '/' || r == 's':
		App.Tview.SetFocus(m.Search)
		return &tcell.EventKey{}
	case r == 'a' || r == 'A':
		App.Tview.SetFocus(m.Articles)
	case r == 'b' || r == 'B':
		App.Tview.SetFocus(m.Basket)
	case r == 'd' || r == 'D':
		ShowDeliveryPage()
	case r == 'c' || r == 'C':
		ShowCheckoutPage()
	case r == 't' || r == 'T':
		ShowDeliveryTrackPage()
	case r == 'h' || r == 'H':
		m.ShowBundles = !m.ShowBundles
		go m.renderSearch(m.Search.GetText())
	case r == 'f' || r == 'F':
		var itemId = ""
		if m.Articles.HasFocus() {
			_, itemId = m.Articles.GetItemText(m.Articles.GetCurrentItem())
		}
		if m.Basket.HasFocus() {
			_, itemId = m.Basket.GetItemText(m.Basket.GetCurrentItem())
		}
		if itemId == "" {
			return event
		}
		go m.renderArticleDetails(itemId)
	}
	return event
}

func (m *MainPage) renderArticleDetails(articleId string) {
	art, artErr := App.Client.GetArticleDetails(articleId)
	if artErr != nil {
		ShowErrorPage(artErr.Error())
		return
	}
	App.Tview.QueueUpdateDraw(func() {
		m.renderArticleImage(art)
		m.renderArticleInfo(art)
	})
}

func (m *MainPage) setupBasketSelected(i int, s string, s2 string, r rune) {
	cart, err := App.Client.AddToCart(s2, 1)
	if err != nil {
		ShowErrorPage(err.Error())
		return
	}
	App.Cart = cart
	go m.renderBasket()
}

func (m *MainPage) setupBasketHotkeys(event *tcell.EventKey) *tcell.EventKey {
	switch key := event; {
	case key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyDEL:
		_, s2 := m.Basket.GetItemText(m.Basket.GetCurrentItem())
		App.Cart, _ = App.Client.RemoveFromCart(s2, 1)
		go m.renderBasket()
	case '0' <= key.Rune() && event.Rune() <= '9':
		_, s2 := m.Basket.GetItemText(m.Basket.GetCurrentItem())
		App.Cart, _ = App.Client.AddToCart(s2, int(event.Rune()-'0'))
		go m.renderBasket()
	case key.Key() == tcell.KeyCtrlK:
		App.Cart, _ = App.Client.ClearCart()
		go m.renderBasket()
	}
	return event
}

func (m *MainPage) setupArticleSelected(i int, s string, s2 string, r rune) {
	App.Cart, _ = App.Client.AddToCart(s2, 1)
	go m.renderBasket()
}

func (m *MainPage) setupArticleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	switch key := event; {
	case '0' <= key.Rune() && event.Rune() <= '9':
		_, s2 := m.Articles.GetItemText(m.Articles.GetCurrentItem())
		App.Cart, _ = App.Client.AddToCart(s2, int(event.Rune()-'0'))
		go m.renderBasket()
	}
	return event
}

func (m *MainPage) renderArticleImage(art *picnic.ArticleDetails) {
	image, imgErr := App.Client.GetArticleImage(art.Images[0].ImageId, picnic.Medium)
	if imgErr != nil {
		ShowErrorPage(imgErr.Error())
		return
	}
	m.ArticleImage.SetImage(*image)
}

func (m *MainPage) renderArticleInfo(art *picnic.ArticleDetails) {
	m.ArticleInfo.Clear()
	var articleText string
	if art.GetPromotion() != "" {
		articleText = fmt.Sprintf("[orange]Name:[white] %s\n[orange]Quantity:[white] %s\n[orange]Price:[white] %s\n[green]Promotion:[white] %s\n[orange]Description:[white] %s",
			art.Name,
			art.UnitQuantity,
			FormatIntToPrice(art.PriceInfo.Price),
			art.GetPromotion(),
			art.Description.Main,
		)
	} else {
		articleText = fmt.Sprintf("[orange]Name:[white] %s\n[orange]Quantity:[white] %s\n[orange]Price:[white] %s\n[orange]Description:[white] %s",
			art.Name,
			art.UnitQuantity,
			FormatIntToPrice(art.PriceInfo.Price),
			art.Description.Main,
		)
	}

	m.ArticleInfo.SetText(articleText)
}

func (m *MainPage) renderBasket() {
	cart, cartErr := App.Client.GetCart()
	if cartErr != nil {
		ShowErrorPage(cartErr.Error())
		return
	}

	App.Cart = cart
	App.Tview.QueueUpdateDraw(func() {
		m.Basket.Clear()
		m.populateBasketItemList()

		deliveryTime := m.getSelectedDeliveryTime()

		basketTitle := "(B)asket (" + FormatIntToPrice(App.Cart.TotalPrice)
		if deliveryTime != "" {
			basketTitle += " | " + deliveryTime
		}
		basketTitle += ")"
		m.Basket.SetTitle(basketTitle)
	})
}

func (m *MainPage) populateBasketItemList() {
	for _, orderLines := range App.Cart.Items {
		for _, article := range orderLines.Items {
			prefix := ""
			switch {
			case !article.IsAvailable():
				prefix = "[red]"
			case orderLines.IsOnPromotion():
				prefix = "[green]"
			}
			basketItemText := fmt.Sprintf("%s%d: %s - %s", prefix, article.Quantity(), article.Name, FormatIntToPrice(orderLines.PriceIncludingPromotions()))
			m.Basket.AddItem(basketItemText, article.Id, 0, nil)
		}
	}
}

func (m *MainPage) getSelectedDeliveryTime() string {
	for _, slot := range App.Cart.DeliverySlots {
		if slot.Selected {
			start, _ := time.Parse(time.RFC3339, slot.WindowStart)
			end, _ := time.Parse(time.RFC3339, slot.WindowEnd)
			return start.Format("Mon, 02 Jan") + " " + start.Format("15:04") + "-" + end.Format("15:04")
		}
	}
	return ""
}
