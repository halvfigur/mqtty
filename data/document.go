package data

import (
	"io"
	"io/ioutil"
	"time"
)

type Document struct {
	contents []byte
	ts       time.Time
}

func NewDocument(r io.Reader) (*Document, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewDocumentBytes(contents), nil
}

func NewDocumentBytes(contents []byte) *Document {
	return &Document{
		ts:       time.Now(),
		contents: contents,
	}
}

func NewDocumentEmpty() *Document {
	return new(Document)
}

func (d *Document) Contents() []byte {
	return d.contents
}
