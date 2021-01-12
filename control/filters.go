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

func (f *Filters) OnSubscribe(topic string, qos network.Qos) {
	f.ctrl.OnSubscribe(topic, qos, func(err error) {
		f.ctrl.QueueUpdate(func() {
			if err != nil {
				//TODO handle error
				return
			}

			f.filters.Add(model.NewSubscriptionFilter(topic, qos))
			f.view.SetSubscriptionFilters(f.filters)
		})
	})
}

func (f *Filters) OnUnsubscribe(topic string) {
	f.ctrl.OnUnsubscribe(topic, func(err error) {
		f.ctrl.QueueUpdate(func() {
			if err != nil {
				//TODO handle error
				return
			}

			f.filters.Remove(topic)
			f.view.SetSubscriptionFilters(f.filters)
		})
	})
}

func (f *Filters) Cancel() {
	f.ctrl.Hide(filtersLabel)
}

func (f *Filters) OnChangeFocus(p tview.Primitive) {
	f.ctrl.Focus(p)
}
