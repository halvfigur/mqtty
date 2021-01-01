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
	ts := time.Now()

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &Document{
		ts:       ts,
		contents: contents,
	}, nil
}

func (d *Document) Contents() []byte {
	return d.contents
}
