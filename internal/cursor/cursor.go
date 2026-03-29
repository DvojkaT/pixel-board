package cursor

import (
	"board/internal/color"
	"board/internal/name"
	"errors"
	"sync"
)

type Store struct {
	cursors map[string]Cursor
	mtx     sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		cursors: make(map[string]Cursor),
	}
}

type Cursor struct {
	X     float64     `json:"x"`
	Y     float64     `json:"y"`
	Name  name.Name   `json:"name"`
	Color color.Color `json:"color"`
}

func NewCursor(x float64, y float64, name name.Name, color color.Color) *Cursor {
	return &Cursor{
		X:     x,
		Y:     y,
		Name:  name,
		Color: color,
	}
}

func (s *Store) Get(userId string) (*Cursor, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	cursor, ok := s.cursors[userId]
	if !ok {
		return nil, errors.New("cursor not found")
	}
	return &cursor, nil
}

func (s *Store) Update(userID string, x, y float64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if c, ok := s.cursors[userID]; ok {
		c.X = x
		c.Y = y
		s.cursors[userID] = c
	}
}

func (s *Store) Set(userID string, cursor *Cursor) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.cursors[userID] = *cursor
}

func (s *Store) Delete(userID string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.cursors, userID)
}
