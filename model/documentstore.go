package model

import "github.com/halvfigur/mqtty/data"

type DocumentIndex struct {
	current   int
	documents []*data.Document
}

func newDocumentIndxex() *DocumentIndex {
	return &DocumentIndex{
		current:   -1,
		documents: make([]*data.Document, 0, defaultDocumentIndexSize),
	}
}

func (i *DocumentIndex) Add(d *data.Document) {
	i.documents = append(i.documents, d)

	if i.current == -1 {
		i.current = 0
	}
}

func (i *DocumentIndex) Current() (int, *data.Document) {
	if i.documents == nil {
		return -1, data.NewDocumentEmpty()
	}
	return i.current, i.documents[i.current]
}

func (i *DocumentIndex) Next() (int, *data.Document) {
	i.current = (i.current + 1) % len(i.documents)
	return i.Current()
}

func (i *DocumentIndex) Prev() (int, *data.Document) {
	if i.current == -1 {
		panic("index is empty")
	}

	if i.current == 0 {
		i.current = len(i.documents)
	}

	i.current--
	return i.Current()
}

func (i *DocumentIndex) Len() int {
	return len(i.documents)
}

type DocumentStore struct {
	current string
	index   map[string]*DocumentIndex
}

const defaultDocumentIndexSize = 32

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		index: make(map[string]*DocumentIndex),
	}
}

func (s *DocumentStore) SetCurrent(name string) {
	if name == "" {
		panic("invalid index")
	}

	index := s.index[name]
	if index == nil {
		panic("name not in store")
	}

	s.current = name
}

func (s *DocumentStore) Store(t string, d *data.Document) {
	if s.index[t] == nil {
		s.index[t] = newDocumentIndxex()
	}

	s.index[t].Add(d)

	if s.current == "" {
		s.current = t
	}
}

func (s *DocumentStore) Current() (string, *DocumentIndex) {
	return s.current, s.index[s.current]
}

func (s *DocumentStore) Len() int {
	return len(s.index)
}
