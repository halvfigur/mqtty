package control

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
	"github.com/rivo/tview"
)

const subscriptionFiltersViewLabel = "subscriptionsfilterview"

type SubscriptionFiltersViewController struct {
	ctrl    Control
	view    *view.SubscriptionFiltersView
	filters *model.SubscriptionFilters
}

func NewSubscriptionFiltersViewController(ctrl Control) *SubscriptionFiltersViewController {
	return &SubscriptionFiltersViewController{
		ctrl:    ctrl,
		filters: model.NewSubscriptionFilters(),
	}
}

func (c *SubscriptionFiltersViewController) SetView(v *view.SubscriptionFiltersView) {
	c.view = v
}

func (c *SubscriptionFiltersViewController) Subscribe(topic string, qos network.Qos) {
	c.ctrl.Subscribe(topic, qos)
	c.filters.Add(model.NewSubscriptionFilter(topic, qos))
	c.view.SetSubscriptionFilters(c.filters)
}

func (c *SubscriptionFiltersViewController) Cancel() {
	c.ctrl.Cancel()
}

func (c *SubscriptionFiltersViewController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}
