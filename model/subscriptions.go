package model

import (
	"github.com/halvfigur/mqtty/network"
)

const initialSize = 16

type SubscriptionFilter struct {
	name string
	qos  network.Qos
}

func NewSubscriptionFilter(filter string, qos network.Qos) *SubscriptionFilter {
	return &SubscriptionFilter{
		name: filter,
		qos:  qos,
	}
}

func (f *SubscriptionFilter) Name() string {
	return f.name
}

func (f *SubscriptionFilter) Qos() network.Qos {
	return f.qos
}

type SubscriptionFilters []*SubscriptionFilter

func NewSubscriptionFilters() *SubscriptionFilters {
	f := SubscriptionFilters(make([]*SubscriptionFilter, 0, initialSize))
	return &f
}

func (f *SubscriptionFilters) indexOf(name string) int {
	for i, x := range *f {
		if x.name == name {
			return i
		}
	}

	return -1
}

func (f *SubscriptionFilters) Contains(name string) bool {
	return f.indexOf(name) != -1
}

func (f *SubscriptionFilters) Add(filter *SubscriptionFilter) {
	if f.Contains(filter.Name()) {
		return
	}

	*f = append(*f, filter)
}

func (f *SubscriptionFilters) Remove(name string) {
	var idx int
	if idx = f.indexOf(name); idx == -1 {
		return
	}

	*f = append((*f)[:idx], (*f)[idx+1:]...)
}

func (f *SubscriptionFilters) Filters() []*SubscriptionFilter {
	c := make([]*SubscriptionFilter, len(*f))
	copy(c, *f)

	return c
}
