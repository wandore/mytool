package set

import "sync"

type Map struct {
	data map[interface{}]interface{}
	mu   sync.RWMutex
}

func New() *Map {
	return &Map{
		data: make(map[interface{}]interface{}, 0),
		mu:   sync.RWMutex{},
	}
}

func (m *Map) Set(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

func (m *Map) Delete(key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

func (m *Map) Get(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]

	return value, ok
}

func (m *Map) Exist(key interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[key]

	return ok
}

func (m *Map) GetKeys(key interface{}) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]interface{}, 0)

	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}

func (m *Map) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}

func (m *Map) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[interface{}]interface{}, 0)
}
