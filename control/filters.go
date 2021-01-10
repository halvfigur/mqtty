package control

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
	"github.com/rivo/tview"
)

const filtersLabel = "filters"

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

func (f *Filters) GetView() *view.Filters {
	return f.view
}

func (f *Filters) Subscribe(topic string, qos network.Qos) error {
	if err := f.ctrl.Subscribe(topic, qos); err != nil {
		return err
	}
	f.filters.Add(model.NewSubscriptionFilter(topic, qos))
	f.view.SetSubscriptionFilters(f.filters)

	return nil
}

func (f *Filters) Unsubscribe(topic string) error {
	if err := f.ctrl.Unsubscribe(topic); err != nil {
		return err
	}
	f.filters.Remove(topic)
	f.view.SetSubscriptionFilters(f.filters)

	return nil
}

func (f *Filters) Cancel() {
	f.ctrl.Hide(filtersLabel)
}

func (f *Filters) OnChangeFocus(p tview.Primitive) {
	f.ctrl.Focus(p)
}
