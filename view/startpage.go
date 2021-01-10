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

	errorMsgView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetDynamicColors(true)

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

		errorMsgView.Clear()
		if err := ctrl.OnConnect(host, port, username, password); err != nil {
			errorMsgView.SetText(fmt.Sprint("[red]Failed to connect:[-] ", err.Error()))
		}
	})

	form.AddButton("Quit", func() {
		ctrl.Stop()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Connection").SetBorder(true)
	flex.AddItem(form, 0, 1, true).
		AddItem(errorMsgView, 1, 0, false)

	p.Flex = Center(flex, 1, 1)
	return p
}
