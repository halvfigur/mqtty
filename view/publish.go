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
		OnPublish(topic string, qos network.Qos, retained bool, message []byte) error
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
	const (
		topicLabel    = "Topic"
		qosLabel      = "Qos"
		retainedLabel = "Retained"
		topicWidth    = 32
	)
	p := &Publish{
		ctrl:     ctrl,
		renderer: model.NewRawRenderer(),
	}

	topicField := tview.NewInputField().
		SetLabel("Topic: ")
	qosDropDown := tview.NewDropDown().
		SetLabel("Qos: ").
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

	fc := NewFocusChain(topicField, qosDropDown, retainedCheckbox,
		launchEditorButton, openFileButton, openHistoryButton)

	publishButton := tview.NewButton("Publish").
		SetSelectedFunc(func() {
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

			if err := ctrl.OnPublish(topic, qos, retained, p.data); err != nil {
				// TODO handle error
			}
		})

	fc.Add(publishButton)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(launchEditorButton, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(openFileButton, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(openHistoryButton, 0, 1, false).
		AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Publish").SetBorder(true)
	flex.AddItem(topicField, 1, 0, true).
		AddItem(qosDropDown, 1, 0, false).
		AddItem(retainedCheckbox, 1, 0, false).
		AddItem(buttonFlex, 1, 0, false).
		AddItem(widget.NewDivider().SetLabel("Data"), 1, 0, false).
		AddItem(p.dataView, 0, 1, false).
		AddItem(publishButton, 1, 1, false)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
	p.Flex = Center(flex, 3, 2)
	return p
}

func (p *Publish) SetData(data []byte) {
	p.dataView.Clear()
	printable, _ := p.renderer.Render(data)
	p.dataView.Write(printable)
	p.data = data
}
