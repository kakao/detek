package detek

import (
	"fmt"
	"reflect"
	"sync"
)

type TypedStored[T any] struct {
	Value T
	Stored
}

type Stored struct {
	Value      interface{}
	Type       reflect.Type
	ProducedBy *MetaInfo
}
type Store struct {
	kv map[string]Stored
	mu sync.RWMutex
}

func (s *Store) Get(key string) (interface{}, *Stored, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.kv[key]; ok {
		return v.Value, &v, nil
	}
	return nil, nil, NewError(nil, ErrKeyNotFound)
}

func (s *Store) Set(key string, val *Stored) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if val == nil {
		return fmt.Errorf("can not set nil value in store")
	}
	if val.Type == nil {
		val.Type = reflect.TypeOf(val.Value)
	}
	if val.Type != reflect.TypeOf(val.Value) {
		return fmt.Errorf("value and type of the value not matched: %q != %q", val.Type, reflect.TypeOf(val.Value))
	}
	if val.ProducedBy == nil {
		return fmt.Errorf("producer not specified")
	}
	if val.ProducedBy.ID == "" {
		return fmt.Errorf("producer name not specified")
	}
	s.kv[key] = *val
	return nil
}
