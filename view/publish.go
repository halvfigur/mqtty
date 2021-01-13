package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	PublishControl interface {
		OnChangeFocus(p tview.Primitive)
		OnLaunchEditor()
		OnOpenFile()
		OnOpenHistory()
		OnPublish(topic string, qos network.Qos, retained bool, message []byte)
		Cancel()
	}

	Publish struct {
		*tview.Flex
		ctrl     PublishControl
		dataView *tview.TextView
		data     []byte
		renderer model.Renderer
	}
)

func NewPublish(ctrl PublishControl) *Publish {
	const topicWidth = 32

	p := &Publish{
		ctrl:     ctrl,
		renderer: model.NewRawRenderer(),
	}

	topicField := tview.NewInputField().
		SetLabel("Topic:    ")
	qosDropDown := tview.NewDropDown().
		SetLabel("Qos:      ").
		SetOptions([]string{"At least once", "At most once", "Exactly once"}, nil)
	retainedCheckbox := tview.NewCheckbox().
		SetLabel("Retained: ")
	launchEditorButton := tview.NewButton("Launch editor").
		SetSelectedFunc(func() {
			ctrl.OnLaunchEditor()
		})
	openFileButton := tview.NewButton("Open file").
		SetSelectedFunc(func() {
			ctrl.OnOpenFile()
		})
	openHistoryButton := tview.NewButton("History").
		SetSelectedFunc(func() {
			ctrl.OnOpenHistory()
		})
	p.dataView = tview.NewTextView()
	publishButton := tview.NewButton("Publish")

	fc := NewFocusChain(topicField, qosDropDown, retainedCheckbox,
		launchEditorButton, openFileButton, openHistoryButton, publishButton)

	publish := func() {
		if topicField.GetText() == "" {
			ctrl.OnChangeFocus(fc.Reset())
		}

		topic := topicField.GetText()
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

		retained := retainedCheckbox.IsChecked()

		ctrl.OnPublish(topic, qos, retained, p.data)
		fc.Reset()
	}

	publishButton.SetSelectedFunc(publish)

	buttonFlex := Space(tview.FlexColumn, launchEditorButton, openFileButton, openHistoryButton)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Publish").SetBorder(true)
	flex.AddItem(topicField, 1, 0, true).
		AddItem(qosDropDown, 1, 0, false).
		AddItem(retainedCheckbox, 1, 0, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(buttonFlex, 1, 0, false).
		AddItem(widget.NewDivider().SetLabel("Data"), 1, 0, false).
		AddItem(p.dataView, 0, 1, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(publishButton, 1, 1, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[blue](TAB):[-] navigate  [blue](^E):[-] launch editor  [blue](^F):[-] open file  [blue](^H):[-] open history  [blue](^P):[-] publish"),
			1, 0, false)

	flex = Center(flex, 300, 200)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			ctrl.OnChangeFocus(fc.Prev())
		case tcell.KeyEscape:
			ctrl.Cancel()
		case tcell.KeyCtrlE:
			ctrl.OnLaunchEditor()
		case tcell.KeyCtrlF:
			ctrl.OnOpenFile()
		case tcell.KeyCtrlH:
			ctrl.OnOpenHistory()
		case tcell.KeyCtrlP:
			publish()
		}

		return event
	})

	p.Flex = flex
	return p
}

func (p *Publish) SetData(data []byte) {
	p.dataView.Clear()
	printable, _ := p.renderer.Render(data)
	p.dataView.Write(printable)
	p.data = data
}
