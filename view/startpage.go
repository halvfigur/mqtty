package view

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/rivo/tview"
)

type (
	StartPageController interface {
		OnConnect(host string, port int, username, password string)
		OnExit()
	}

	StartPage struct {
		*tview.Form
		ctrl StartPageController
	}
)

func NewStartPage(ctrl StartPageController) *StartPage {
	const (
		hostLabel     = "Host"
		portLabel     = "Port"
		usernameLabel = "Username"
		passwordLabel = "Password"

		textFieldWidth   = 32
		numberFieldWidth = 5
		defaultPort      = 1883
	)

	p := &StartPage{
		Form: tview.NewForm(),
		ctrl: ctrl,
	}

	host := "test.mosquitto.org"
	port := defaultPort
	username := ""
	password := ""

	p.AddInputField(hostLabel, host, textFieldWidth, nil, func(text string) {
		host = text
	})

	validatePort := func(text string, lastChar rune) bool {
		var re = regexp.MustCompile(`^\d+$`)
		m := !re.Match([]byte(text))
		if !m {
			port = -1
		}

		return m
	}
	p.AddInputField(portLabel, fmt.Sprintf("%d", defaultPort), numberFieldWidth, validatePort, func(text string) {
		port, _ = strconv.Atoi(text)
	})

	p.AddInputField(usernameLabel, "", textFieldWidth, nil, func(text string) {
		username = text
	})

	p.AddPasswordField(passwordLabel, "", textFieldWidth, '*', func(text string) {
		password = text
	})

	p.AddButton("Connect", func() {
		if host == "" {
			p.SetFocus(p.Form.GetFormItemIndex(hostLabel))
			return
		}
		if port == -1 {
			p.SetFocus(p.Form.GetFormItemIndex(portLabel))
			return
		}
		if username != "" && password == "" {
			p.SetFocus(p.Form.GetFormItemIndex(passwordLabel))
			return
		}
		if username == "" && password != "" {
			p.SetFocus(p.Form.GetFormItemIndex(usernameLabel))
			return
		}

		ctrl.OnConnect(host, port, username, password)
	})

	p.AddButton("Exit", func() {
		ctrl.OnExit()
	})

	return p
}
