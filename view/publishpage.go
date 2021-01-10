package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	PublishPageControl interface {
		OnChangeFocus(p tview.Primitive)
		OnLaunchEditor()
		OnOpenFile()
		Publish(topic string, qos network.Qos, retained bool, message []byte) error
		Cancel()
	}

	PublishPage struct {
		*tview.Flex
		ctrl     PublishPageControl
		dataView *tview.TextView
		data     []byte
		renderer model.Renderer
	}
)

func NewPublishPage(ctrl PublishPageControl) *PublishPage {
	const (
		topicLabel    = "Topic"
		qosLabel      = "Qos"
		retainedLabel = "Retained"
		topicWidth    = 32
	)
	p := &PublishPage{
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
	p.dataView = tview.NewTextView()

	fc := NewFocusChain(topicField, qosDropDown, retainedCheckbox, launchEditorButton, openFileButton)
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

			if err := ctrl.Publish(topic, qos, retained, p.data); err != nil {
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

func (p *PublishPage) SetData(data []byte) {
	p.dataView.Clear()
	printable, _ := p.renderer.Render(data)
	p.dataView.Write(printable)
	p.data = data
}
