package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

const filterMaxWidth = 32

type (
	FiltersController interface {
		OnSubscribe(topic string, qos network.Qos)
		OnUnsubscribe(topic string)
		OnChangeFocus(p tview.Primitive)
		Cancel()
	}

	Filters struct {
		*tview.Flex
		filters *tview.List
		model   *model.SubscriptionFilters
	}
)

func NewFilters(ctrl FiltersController) *Filters {
	filterInput := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldWidth(filterMaxWidth)
	filterInput.SetBorderPadding(1, 1, 1, 1)

	qosOpts := []string{"At most once", "At least once", "Exatly once"}
	qosDropDown := tview.NewDropDown().
		SetOptions(qosOpts, nil).
		SetLabel("Qos: ").
		SetFieldWidth(0).
		SetCurrentOption(0)
	qosDropDown.SetBorderPadding(1, 1, 1, 1)

	addButton := tview.NewButton("Add")
	clearButton := tview.NewButton("Clear")
	filterList := tview.NewList()

	fc := NewFocusChain(filterInput, qosDropDown, addButton, clearButton, filterList)

	subscribe := func() {
		defer ctrl.OnChangeFocus(fc.Reset())

		// Get filter name
		if filterInput.GetText() == "" {
			return
		}

		// Get qos option
		o, _ := qosDropDown.GetCurrentOption()

		var qos network.Qos
		switch o {
		case 0:
			qos = network.QosAtMostOnce
		case 1:
			qos = network.QosAtLeastOnce
		case 2:
			qos = network.QosExatlyOnce
		}

		ctrl.OnSubscribe(filterInput.GetText(), qos)
	}

	filterFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(filterInput, 0, 5, true).
		AddItem(qosDropDown, 0, 5, false)

	filterList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
			i := filterList.GetCurrentItem()
			name, _ := filterList.GetItemText(i)
			ctrl.OnUnsubscribe(name)
		}

		return event
	})

	addButton.SetSelectedFunc(func() {
		subscribe()
	})

	clearButton.SetSelectedFunc(func() {
		filterInput.SetText("")
		qosDropDown.SetCurrentOption(0)
		ctrl.OnChangeFocus(fc.Reset())
	})

	buttonFlex := Space(tview.FlexColumn, addButton, clearButton)

	viewFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(filterFlex, 3, 0, true).
		AddItem(widget.NewDivider().SetLabel("Applied"), 1, 0, false).
		AddItem(filterList, 0, 1, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(buttonFlex, 1, 0, false)

	viewFlex.SetTitle("Filters").SetBorder(true)

	viewFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			ctrl.OnChangeFocus(fc.Prev())
		case tcell.KeyEscape:
			ctrl.Cancel()
		case tcell.KeyEnter:
			if filterInput.HasFocus() {
				subscribe()
			}
		}

		return event
	})

	return &Filters{
		Flex:    Center(viewFlex, 100, 100),
		filters: filterList,
	}
}

func (f *Filters) SetSubscriptionFilters(filters *model.SubscriptionFilters) {
	f.model = filters

	maxLen := 0
	for _, f := range f.model.Filters() {
		if l := len(f.Name()); l > maxLen {
			maxLen = l
		}
	}

	f.filters.Clear()
	f.filters.ShowSecondaryText(false)
	for _, filter := range f.model.Filters() {
		f.filters.AddItem(filter.Name(), "", 0, nil)
	}
}
