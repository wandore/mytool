package set

import "sync"

var empty = struct{}{}

type Set struct {
	elems map[interface{}]struct{}
	mu    sync.RWMutex
}

func New() *Set {
	return &Set{
		elems: make(map[interface{}]struct{}, 0),
		mu:    sync.RWMutex{},
	}
}

func (s *Set) Add(elem interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.elems[elem] = empty
}

func (s *Set) Delete(elem interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.elems, elem)
}

func (s *Set) Exist(elem interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.elems[elem]

	return ok
}

func (s *Set) GetElems() []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	elems := make([]interface{}, 0)

	for elem := range s.elems {
		elems = append(elems, elem)
	}

	return elems
}

func (s *Set) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.elems)
}

func (s *Set) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.elems = make(map[interface{}]struct{}, 0)
}