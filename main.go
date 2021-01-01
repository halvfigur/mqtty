package main

import (
	"bytes"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

func incomingPage() tview.Primitive {
	buf := bytes.NewBufferString(`
	{
		"alpha": {
			"beta": [0, 1, 2, 3},
			"gamma": "hello world",
			"epsilon": null
		}
	}
	`)

	// Document data
	d, err := data.NewDocument(buf)
	if err != nil {
		log.Fatal("Failed to read document")
	}

	// Document model
	m := model.NewDocument()
	m.SetRenderer(new(model.HexRenderer))
	m.SetDocument(d)

	// Document view
	v := view.NewDocumentView()
	v.SetDocument(m)

	// Root window layout
	flex := tview.NewFlex().AddItem(v.SetBorder(true), 0, 1, false)

	return flex
}

func main() {
	app := tview.NewApplication()

	pages := tview.NewPages()
	pages.AddPage("incoming-page", incomingPage(), true)
	pages.AddPage("quit-modal", tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}
			if buttonLabel == "Cancel" {
				pages.SwitchToPage("incoming-page")
			}
		}), false)

	app.SetRoot(pages, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() != 'c' {
			return event
		}

		modal := tview.NewModal().SetText("Do you want to quit the application?").
			AddButtons([]string{"Quit", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					app.Stop()
				}
				if buttonLabel == "Cancel" {
				}
			})

		app.SetRoot("quit-page", true)

		// Returning nil prevents the event to stop further event processing
		return nil
	})
	if err := app..Run(); err != nil {
		log.Fatal(err)
	}
}
