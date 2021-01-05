package model

import (
	"github.com/halvfigur/mqtty/data"
)

type (
	Renderer interface {
		Name() string
		Render(data []byte) ([]byte, bool)
	}

	RawRenderer struct{}
)

func NewRawRenderer() *RawRenderer {
	return new(RawRenderer)
}

func (r *RawRenderer) Name() string {
	return "Raw"
}

func (r *RawRenderer) Render(data []byte) ([]byte, bool) {
	return data, false
}

type (
	Document struct {
		doc      *data.Document
		renderer Renderer
	}
)

func NewDocument() *Document {
	return &Document{
		renderer: NewRawRenderer(),
	}
}

func (d *Document) SetDocument(doc *data.Document) {
	d.doc = doc
}

func (d *Document) SetRenderer(r Renderer) {
	d.renderer = r
}

func (d *Document) Contents() ([]byte, bool) {
	return d.renderer.Render(d.doc.Contents())
}
