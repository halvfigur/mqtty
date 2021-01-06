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
	c := &SubscriptionFiltersViewController{
		ctrl:    ctrl,
		filters: model.NewSubscriptionFilters(),
	}

	c.view = view.NewSubscriptionFiltersView(c)

	ctrl.Register(subscriptionFiltersViewLabel, c.view, false)

	return c
}

func (c *SubscriptionFiltersViewController) GetView() *view.SubscriptionFiltersView {
	return c.view
}

func (c *SubscriptionFiltersViewController) Subscribe(topic string, qos network.Qos) error {
	if err := c.ctrl.Subscribe(topic, qos); err != nil {
		return err
	}
	c.filters.Add(model.NewSubscriptionFilter(topic, qos))
	c.view.SetSubscriptionFilters(c.filters)

	return nil
}

func (c *SubscriptionFiltersViewController) Unsubscribe(topic string) error {
	if err := c.ctrl.Unsubscribe(topic); err != nil {
		return err
	}
	c.filters.Remove(topic)
	c.view.SetSubscriptionFilters(c.filters)

	return nil
}

func (c *SubscriptionFiltersViewController) Cancel() {
	c.ctrl.Hide(subscriptionFiltersViewLabel)
}

func (c *SubscriptionFiltersViewController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}
