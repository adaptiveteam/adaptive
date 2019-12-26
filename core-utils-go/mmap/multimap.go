package mmap

type MultiMap interface {
	Put(key, value interface{})
	Get(key interface{}) (values []interface{}, found bool)
	PutAll(key interface{}, values []interface{})
	Remove(key interface{}, value interface{})
	RemoveAll(key interface{})
	Contains(key interface{}, value interface{}) bool
	ContainsKey(key interface{}) (found bool)
	ContainsValue(value interface{}) bool
	Size() int
	KeySet() []interface{}
	Keys() []interface{}
	Values() []interface{}
	Empty() bool
	Clear()
}

// MultiMap holds the elements in go's native map.
type multiMap struct {
	m map[interface{}][]interface{}
}

func NewMultiMap() MultiMap {
	return &multiMap{m: make(map[interface{}][]interface{})}
}

// Put stores a key-value pair in this multimap.
func (m *multiMap) Put(key interface{}, value interface{}) {
	m.m[key] = append(m.m[key], value)
}

// Get searches the element in the multimap by key.
// It returns its value or nil if key is not found in multimap.
// Second return parameter is true if key was found, otherwise false.
func (m *multiMap) Get(key interface{}) (values []interface{}, found bool) {
	values, found = m.m[key]
	return
}

// PutAll stores a key-value pair in this multimap for each of the values, all using the same key key.
func (m *multiMap) PutAll(key interface{}, values []interface{}) {
	for _, value := range values {
		m.Put(key, value)
	}
}

// Contains returns true if this multimap contains at least one key-value pair with the key key and the value value.
func (m *multiMap) Contains(key interface{}, value interface{}) bool {
	values, found := m.m[key]
	for _, v := range values {
		if v == value {
			return true && found
		}
	}
	return false && found
}

// ContainsKey returns true if this multimap contains at least one key-value pair with the key key.
func (m *multiMap) ContainsKey(key interface{}) (found bool) {
	_, found = m.m[key]
	return
}

// ContainsValue returns true if this multimap contains at least one key-value pair with the value value.
func (m *multiMap) ContainsValue(value interface{}) bool {
	for _, values := range m.m {
		for _, v := range values {
			if v == value {
				return true
			}
		}
	}
	return false
}

// Remove removes a single key-value pair from this multimap, if such exists.
func (m *multiMap) Remove(key interface{}, value interface{}) {
	values, found := m.m[key]
	if found {
		for i, v := range values {
			if v == value {
				m.m[key] = append(values[:i], values[i+1:]...)
			}
		}
	}
	if len(m.m[key]) == 0 {
		delete(m.m, key)
	}
}

// RemoveAll removes all values associated with the key from the multimap.
func (m *multiMap) RemoveAll(key interface{}) {
	delete(m.m, key)
}

// Size returns number of key-value pairs in the multimap.
func (m *multiMap) Size() int {
	size := 0
	for _, value := range m.m {
		size += len(value)
	}
	return size
}

// Keys returns a view collection containing the key from each key-value pair in this multimap.
// This is done without collapsing duplicates.
func (m *multiMap) Keys() []interface{} {
	keys := make([]interface{}, m.Size())
	count := 0
	for key, value := range m.m {
		for range value {
			keys[count] = key
			count++
		}
	}
	return keys
}

// KeySet returns all distinct keys contained in this multimap.
func (m *multiMap) KeySet() []interface{} {
	keys := make([]interface{}, len(m.m))
	count := 0
	for key := range m.m {
		keys[count] = key
		count++
	}
	return keys
}

// Values returns all values from each key-value pair contained in this multimap.
// This is done without collapsing duplicates. (size of Values() = MultiMap.Size()).
func (m *multiMap) Values() []interface{} {
	values := make([]interface{}, m.Size())
	count := 0
	for _, vs := range m.m {
		for _, value := range vs {
			values[count] = value
			count++
		}
	}
	return values
}

func (m *multiMap) Empty() bool {
	return m.Size() == 0
}

// Clear removes all elements from the map.
func (m *multiMap) Clear() {
	m.m = make(map[interface{}][]interface{})
}
