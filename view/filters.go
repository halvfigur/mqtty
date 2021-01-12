package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
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

	errorMsgView := tview.NewTextView().
		SetWrap(true).
		SetWordWrap(true).
		SetDynamicColors(true)

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

		errorMsgView.Clear()
		ctrl.OnSubscribe(filterInput.GetText(), qos)
		//errorMsgView.SetText(fmt.Sprint("[red]Failed to subscribe:[-] ", err.Error()))
	}

	addButton := tview.NewButton("Add").SetSelectedFunc(func() {
		subscribe()
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
	filterList.SetBorder(true).SetTitle("Current")
	filterList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
			i := filterList.GetCurrentItem()
			name, _ := filterList.GetItemText(i)
			ctrl.OnUnsubscribe(name)
			//errorMsgView.SetText(fmt.Sprint("[red]Unsubscribe failed:[-], ", err.Error()))
		}

		return event
	})

	viewFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	viewFlex.SetTitle("Filters").SetBorder(true)
	viewFlex.AddItem(filterFlex, 3, 0, true)
	viewFlex.AddItem(filterList, 0, 1, false)
	viewFlex.AddItem(errorMsgView, 1, 0, false)

	fc.Add(addButton, clearButton, filterList)

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
		Flex:    Center(viewFlex, 1, 1),
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
