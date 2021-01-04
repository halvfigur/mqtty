package main

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

const (
	defaultDocumentIndexSize = 32
)

type (
	documentIndex struct {
		current   int
		documents []*data.Document
	}

	documentStore struct {
		current string
		index   map[string]*documentIndex
	}

	MainPageController struct {
		app       *tview.Application
		view      *view.MainPage
		model     *model.Document
		documents *documentStore
	}
)

func newDocumentIndxex() *documentIndex {
	return &documentIndex{
		current:   -1,
		documents: make([]*data.Document, 0, defaultDocumentIndexSize),
	}
}

func (i *documentIndex) Add(d *data.Document) {
	i.documents = append(i.documents, d)

	if i.current == -1 {
		i.current = 0
	}
}

func (i *documentIndex) Current() (int, *data.Document) {
	return i.current, i.documents[i.current]
}

func (i *documentIndex) Next() (int, *data.Document) {
	i.current = (i.current + 1) % len(i.documents)
	return i.Current()
}

func (i *documentIndex) Prev() (int, *data.Document) {
	if i.current == -1 {
		panic("index is empty")
	}

	if i.current == 0 {
		i.current = len(i.documents)
	}

	i.current--
	return i.Current()
}

func (i *documentIndex) Len() int {
	return len(i.documents)
}

func newDocumentStore() *documentStore {
	return &documentStore{
		index: make(map[string]*documentIndex),
	}
}

func (s *documentStore) SetCurrent(name string) {
	if name == "" {
		panic("invalid index")
	}

	index := s.index[name]
	if index == nil {
		panic("name not in store")
	}

	s.current = name
}

func (s *documentStore) Store(t string, d *data.Document) {
	if s.index[t] == nil {
		s.index[t] = newDocumentIndxex()
	}

	s.index[t].Add(d)

	if s.current == "" {
		s.current = t
	}
}

func (s *documentStore) Current() (string, *documentIndex) {
	return s.current, s.index[s.current]
}

func (s *documentStore) Len() int {
	return len(s.index)
}

func NewMainPageController(a *tview.Application) *MainPageController {
	return &MainPageController{
		app:       a,
		model:     model.NewDocument(),
		documents: newDocumentStore(),
	}
}

func (c *MainPageController) SetView(v *view.MainPage) {
	c.view = v
}

func (c *MainPageController) SetDocument(d *data.Document) {
	c.model.SetDocument(d)
	if c.view != nil {
		c.view.SetDocument(c.model)
	}
}

func (c *MainPageController) AddDocument(t string, d *data.Document) {
	c.documents.Store(t, d)

	if c.view != nil {
		c.view.AddTopic(t)
		c.updateDocumentView()
	}
}

func (c *MainPageController) OnTopicSelected(t string) {
	c.documents.SetCurrent(t)
	c.updateDocumentView()
}

func (c *MainPageController) OnRendererSelected(r view.Renderer) {
	if c.view != nil {
		c.view.SetRenderer(r)
	}
}

func (c *MainPageController) OnChangeFocus(p tview.Primitive) {
	c.app.SetFocus(p)
}

func (c *MainPageController) OnNextDocument() {
	_, index := c.documents.Current()
	index.Next()
	c.updateDocumentView()
}

func (c *MainPageController) OnPrevDocument() {
	_, index := c.documents.Current()
	index.Prev()
	c.updateDocumentView()
}

func (c *MainPageController) updateDocumentView() {
	t, index := c.documents.Current()
	i, d := index.Current()
	c.model.SetDocument(d)
	c.view.SetDocument(c.model)
	c.view.SetTopicsTitle(fmt.Sprintf("Topics %d", c.documents.Len()))
	c.view.SetDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))
}
