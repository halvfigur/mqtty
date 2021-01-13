package model

import (
	"sort"
	"strings"

	"github.com/halvfigur/mqtty/data"
)

type DocumentIndex struct {
	current   int
	documents []*data.Document
	follow    bool
}

func newDocumentIndxex() *DocumentIndex {
	return &DocumentIndex{
		current:   -1,
		documents: make([]*data.Document, 0, defaultDocumentIndexSize),
		follow:    false,
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

func (i *DocumentIndex) MoveToFirst() {
	if i.current == -1 {
		return
	}

	i.current = 0
}

func (i *DocumentIndex) MoveToLast() {
	if len(i.documents) == 0 {
		return
	}

	i.current = len(i.documents) - 1
}

func (i *DocumentIndex) Follow(enabled bool) {
	if i.follow {
		i.MoveToLast()
	}

	i.follow = enabled
}

func (i *DocumentIndex) Len() int {
	return len(i.documents)
}

type DocumentStore struct {
	current string
	index   map[string]*DocumentIndex
	follow  bool
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

func (s *DocumentStore) Next() (int, *data.Document) {
	return s.index[s.current].Next()
}

func (s *DocumentStore) Prev() (int, *data.Document) {
	return s.index[s.current].Prev()
}

func (s *DocumentStore) MoveToFirst() {
	for _, i := range s.index {
		i.MoveToFirst()
	}
}

func (s *DocumentStore) MoveToLast() {
	for _, i := range s.index {
		i.MoveToLast()
	}
}

func (s *DocumentStore) Follow(enabled bool) {
	for _, i := range s.index {
		i.Follow(enabled)
	}

	s.follow = enabled
	if enabled {
		s.MoveToLast()
	}
}

func (s *DocumentStore) Store(t string, d *data.Document) {
	if s.index[t] == nil {
		s.index[t] = newDocumentIndxex()
	}

	s.index[t].Add(d)

	if s.current == "" {
		s.current = t
	}

	if s.follow {
		s.MoveToLast()
	}
}

func (s *DocumentStore) Current() (string, *DocumentIndex) {
	return s.current, s.index[s.current]
}

func (s *DocumentStore) Topics() []string {
	topics := make([]string, 0, len(s.index))

	for t := range s.index {
		topics = append(topics, t)
	}

	sort.Slice(topics, func(i, j int) bool {
		return strings.Compare(topics[i], topics[j]) < 0
	})

	return topics
}

func (s *DocumentStore) Len() int {
	return len(s.index)
}

func (s *DocumentStore) DocumentCount() int {
	c := 0

	for _, index := range s.index {
		c += index.Len()
	}

	return c
}
