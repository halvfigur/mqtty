package view

import (
	"github.com/halvfigur/mqtty/network"
	"github.com/rivo/tview"
)

type (
	SubscribePageController interface {
		Subscribe(topic string, qos network.Qos)
		Cancel()
	}
)

func NewSubscribePage(ctrl SubscribePageController) *tview.Form {
	topic := ""
	qos := network.QosAtMostOnce

	return tview.NewForm().
		AddInputField("Topic", "", 0, nil, func(text string) {
			topic = text
		}).
		AddDropDown("Qos", []string{"At most once", "At least once", "Exatly once"}, 0, func(opt string, index int) {
			switch index {
			case 0:
				qos = network.QosAtMostOnce
			case 1:
				qos = network.QosAtLeastOnce
			case 2:
				qos = network.QosExatlyOnce
			}
		}).
		AddButton("OK", func() {
			ctrl.Subscribe(topic, qos)
		}).
		AddButton("Cancel", func() {
			ctrl.Cancel()
		})
}
