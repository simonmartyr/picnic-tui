package main

import (
	"errors"
	"flag"
	"github.com/rivo/tview"
	picnic "github.com/simonmartyr/picnic-api"
	"net/http"
)

var (
	App *PicnicTui
)

type PicnicTui struct {
	Client *picnic.Client
	Tview  *tview.Application
	Pages  *tview.Pages
	Main   *MainPage
	Cart   *picnic.Order
}

func Start() {
	c, clientErr := configureClient()

	App = &PicnicTui{
		Client: c,
		Tview:  tview.NewApplication(),
		Pages:  tview.NewPages(),
	}

	//err := App.Client.Authenticate()
	App.Tview.SetRoot(App.Pages, true).SetFocus(App.Pages)

	if clientErr == nil {
		ShowWelcomePage()
	} else {
		ShowErrorPage(clientErr.Error())
	}

	if err := App.Tview.Run(); err != nil {
		panic(err)
	}
}

func configureClient() (*picnic.Client, error) {
	token := flag.String("t", "", "access token for authentication")
	username := flag.String("u", "", "username for authentication (required if token not set)")
	password := flag.String("p", "", "password for authentication")
	hashedPassword := flag.String("hp", "", "md5 hashed password for authentication")
	flag.Parse()

	if *token != "" {
		return picnic.New(&http.Client{},
			picnic.WithToken(*token),
		), nil
	}

	if *username == "" {
		return nil, errors.New("client could not be configured: username (-u) not provided, alternatively provider an access token (-t)")
	}

	var c *picnic.Client
	switch {
	case *hashedPassword != "":
		c = picnic.New(&http.Client{},
			picnic.WithUsername(*username),
			picnic.WithHashedPassword(*hashedPassword),
		)
	case *password != "":
		c = picnic.New(&http.Client{},
			picnic.WithUsername(*username),
			picnic.WithPassword(*password),
		)
	default:
		return nil, errors.New("client could not be configured: password (-p) or hashedPassword (-hp) not provided")
	}
	err := c.Authenticate()
	if err != nil {
		return nil, errors.New("client could not be configured: authentication failed - please verify your username and password")
	}

	return c, nil
}
