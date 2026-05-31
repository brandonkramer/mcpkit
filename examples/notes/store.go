package main

import (
	"fmt"
	"sync"
)

type Note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Store struct {
	mu     sync.Mutex
	nextID int
	items  []Note
}

func NewStore(seed ...Note) *Store {
	s := &Store{nextID: 1}
	for _, n := range seed {
		s.items = append(s.items, n)
		s.nextID++
	}
	return s
}

func (s *Store) List() []Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Note, len(s.items))
	copy(out, s.items)
	return out
}

func (s *Store) Get(id string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range s.items {
		if n.ID == id {
			return n, nil
		}
	}
	return Note{}, fmt.Errorf("note %q not found", id)
}

func (s *Store) Create(title, body string) Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := Note{ID: fmt.Sprintf("%02d", s.nextID), Title: title, Body: body}
	s.nextID++
	s.items = append(s.items, n)
	return n
}
