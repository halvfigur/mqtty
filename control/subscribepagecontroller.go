package control

import (
	"github.com/halvfigur/mqtty/network"
)

const subscribePageLabel = "subscribepage"

type (
	SubscribePageController struct {
		ctrl Control
	}
)

func NewSubscribePageController(ctrl Control) *SubscribePageController {
	return &SubscribePageController{
		ctrl: ctrl,
	}
}

func (c *SubscribePageController) OnSubscribe(topic string, qos network.Qos) {
	c.ctrl.Subscribe(topic, qos)
}

func (c *SubscribePageController) OnCancel() {
	c.ctrl.Cancel()
}
