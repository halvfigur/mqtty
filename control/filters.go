package control

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
	"github.com/rivo/tview"
)

type Filters struct {
	ctrl    Control
	view    *view.Filters
	filters *model.SubscriptionFilters
}

func NewFilters(ctrl Control) *Filters {
	f := &Filters{
		ctrl:    ctrl,
		filters: model.NewSubscriptionFilters(),
	}

	f.view = view.NewFilters(f)

	ctrl.Register(filtersLabel, f.view, false)

	return f
}

func (f *Filters) AddFilter(filter string, qos network.Qos) {
	f.filters.Add(model.NewSubscriptionFilter(filter, qos))
	f.view.SetSubscriptionFilters(f.filters)
}

func (f *Filters) RemoveFilter(filter string) {
	f.filters.Remove(filter)
	f.view.SetSubscriptionFilters(f.filters)
}

func (f *Filters) OnSubscribe(topic string, qos network.Qos) {
	f.ctrl.OnSubscribe(topic, qos)
}

func (f *Filters) OnUnsubscribe(topic string) {
	f.ctrl.OnUnsubscribe(topic)
}

func (f *Filters) Cancel() {
	f.ctrl.Hide(filtersLabel)
}

func (f *Filters) OnChangeFocus(p tview.Primitive) {
	f.ctrl.Focus(p)
}
