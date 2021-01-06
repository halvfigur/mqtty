package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/rivo/tview"
)

const filterMaxWidth = 32

type (
	SubscriptionFiltersViewController interface {
		Subscribe(topic string, qos network.Qos)
		Unsubscribe(topic string)
		OnChangeFocus(p tview.Primitive)
		Cancel()
	}

	SubscriptionFiltersView struct {
		*tview.Flex
		filters *tview.List
		model   *model.SubscriptionFilters
	}
)

func NewSubscriptionFiltersView(ctrl SubscriptionFiltersViewController) *SubscriptionFiltersView {
	filterInput := tview.NewInputField().
		SetLabel("[blue]Filter:[-] ").
		SetFieldWidth(filterMaxWidth).
		SetText("hamweather/#")
	filterInput.SetBorderPadding(1, 1, 1, 1)

	qosOpts := []string{"At most once", "At least once", "Exatly once"}
	qosDropDown := tview.NewDropDown().
		SetOptions(qosOpts, nil).
		SetLabel("[blue]Qos:[-] ").
		SetFieldWidth(0).
		SetCurrentOption(0)
	qosDropDown.SetBorderPadding(1, 1, 1, 1)

	fc := NewFocusChain(filterInput, qosDropDown)

	addButton := tview.NewButton("Add").SetSelectedFunc(func() {
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

		ctrl.Subscribe(filterInput.GetText(), qos)
	})
	addButton.SetBorder(true)

	clearButton := tview.NewButton("Clear").SetSelectedFunc(func() {
		filterInput.SetText("")
		qosDropDown.SetCurrentOption(0)
		ctrl.OnChangeFocus(fc.Reset())
	})
	clearButton.SetBorder(true)

	filterFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(filterInput, 0, 5, true).
		AddItem(qosDropDown, 0, 5, false).
		AddItem(addButton, 0, 1, false).
		AddItem(clearButton, 0, 1, false)

	filterList := tview.NewList()
	filterList.SetBorder(true).SetTitle("[blue]Current[-]")
	filterList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
			i := filterList.GetCurrentItem()
			name, _ := filterList.GetItemText(i)
			ctrl.Unsubscribe(name)
		}

		return event
	})

	viewFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	viewFlex.SetTitle("Filters").SetBorder(true)
	viewFlex.AddItem(filterFlex, 3, 0, true)
	viewFlex.AddItem(filterList, 0, 1, false)

	fc.Add(addButton, clearButton, filterList)

	viewFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			ctrl.OnChangeFocus(fc.Prev())
		case tcell.KeyEscape:
			ctrl.Cancel()
		}

		return event
	})

	return &SubscriptionFiltersView{
		Flex:    center(viewFlex, 1, 1),
		filters: filterList,
	}
}

func (v *SubscriptionFiltersView) SetSubscriptionFilters(filters *model.SubscriptionFilters) {
	v.model = filters

	maxLen := 0
	for _, f := range v.model.Filters() {
		if l := len(f.Name()); l > maxLen {
			maxLen = l
		}
	}

	v.filters.Clear()
	v.filters.ShowSecondaryText(false)
	for _, f := range v.model.Filters() {
		v.filters.AddItem(f.Name(), "", 0, nil)
	}
}
