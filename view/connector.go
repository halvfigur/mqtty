package view

import (
	"regexp"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	ConnectorController interface {
		OnChangeFocus(p tview.Primitive)
		OnConnect(host string, port int, username, password string)
		Cancel()
	}

	Connector struct {
		*tview.Flex
	}
)

func NewConnector(ctrl ConnectorController) *Connector {
	const (
		textFieldWidth   = 32
		numberFieldWidth = 5
		defaultPort      = 1883
	)

	validatePort := func(text string) bool {
		return regexp.MustCompile(`^[1-9]\d*$`).Match([]byte(text))
	}

	hostField := tview.NewInputField().
		SetLabel("    Host: ").
		SetText("test.mosquitto.org")
	portField := tview.NewInputField().
		SetLabel("    Port: ").
		SetFieldWidth(numberFieldWidth).
		SetText("1883")
	usernameField := tview.NewInputField().
		SetLabel("Username: ")
	passwordField := tview.NewInputField().
		SetLabel("Password: ").
		SetMaskCharacter('*')

	fc := NewFocusChain(hostField, portField, usernameField, passwordField)

	connectButton := tview.NewButton("Connect").
		SetSelectedFunc(func() {
			host := hostField.GetText()
			if host == "" {
				ctrl.OnChangeFocus(fc.SetFocus(0))
				return
			}

			sport := portField.GetText()
			if !validatePort(sport) {
				ctrl.OnChangeFocus(fc.SetFocus(1))
				return
			}
			port, _ := strconv.Atoi(sport)

			username := usernameField.GetText()
			password := passwordField.GetText()
			if usernameField.GetText() == "" && passwordField.GetText() != "" {
				ctrl.OnChangeFocus(fc.SetFocus(2))
				return
			}

			ctrl.OnConnect(host, port, username, password)
			fc.Reset()
		})

	cancelButton := tview.NewButton("Cancel").
		SetSelectedFunc(func() {
			ctrl.Cancel()
		})

	fc.Add(connectButton, cancelButton)

	buttonFlex := Space(tview.FlexColumn, connectButton, cancelButton)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(hostField, 1, 0, true).
		AddItem(portField, 1, 0, false).
		AddItem(usernameField, 1, 0, false).
		AddItem(passwordField, 1, 0, false).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(buttonFlex, 1, 0, false)

	flex.SetTitle("Connect").SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				ctrl.OnChangeFocus(fc.Next())
			case tcell.KeyBacktab:
				ctrl.OnChangeFocus(fc.Prev())
			}

			return event
		})

	return &Connector{
		Flex: Center(flex, 1, 1),
	}
}
