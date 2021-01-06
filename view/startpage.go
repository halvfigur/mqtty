package view

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/rivo/tview"
)

type (
	StartPageController interface {
		OnConnect(host string, port int, username, password string) error
		Stop()
	}

	StartPage struct {
		*tview.Flex
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
	host := "test.mosquitto.org"
	port := defaultPort
	username := ""
	password := ""

	p := &StartPage{
		ctrl: ctrl,
	}

	form := tview.NewForm()
	form.AddInputField(hostLabel, host, textFieldWidth, nil, func(text string) {
		host = text
	})

	validatePort := func(text string, lastChar rune) bool {
		var re = regexp.MustCompile(`^[1-9]\d*$`)
		m := !re.Match([]byte(text))
		if !m {
			port = -1
		}

		return m
	}
	form.AddInputField(portLabel, fmt.Sprintf("%d", defaultPort), numberFieldWidth, validatePort, func(text string) {
		port, _ = strconv.Atoi(text)
	})

	form.AddInputField(usernameLabel, "", textFieldWidth, nil, func(text string) {
		username = text
	})

	form.AddPasswordField(passwordLabel, "", textFieldWidth, '*', func(text string) {
		password = text
	})

	form.AddButton("Connect", func() {
		if host == "" {
			form.SetFocus(form.GetFormItemIndex(hostLabel))
			return
		}
		if port == -1 {
			form.SetFocus(form.GetFormItemIndex(portLabel))
			return
		}
		if username != "" && password == "" {
			form.SetFocus(form.GetFormItemIndex(passwordLabel))
			return
		}
		if username == "" && password != "" {
			form.SetFocus(form.GetFormItemIndex(usernameLabel))
			return
		}

		ctrl.OnConnect(host, port, username, password)
	})

	form.AddButton("Exit", func() {
		ctrl.Stop()
	})
	form.SetTitle("Connection").SetBorder(true)

	p.Flex = center(form, 1, 1)
	return p
}
