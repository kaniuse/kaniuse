package main

import "sync"

type TypedSyncMap[K any, V any] interface {
	Load(key K) (V, bool)
	Store(key K, value V)
	Delete(key K)
	Range(f func(key K, value V) bool)
}

var _ TypedSyncMap[any, any] = (*syncMapWrapper[any, any])(nil)

type syncMapWrapper[K any, V any] struct {
	m *sync.Map
}

func newSyncMapWrapper[K any, V any]() *syncMapWrapper[K, V] {
	return &syncMapWrapper[K, V]{
		m: &sync.Map{},
	}
}

func (t *syncMapWrapper[K, V]) Store(key K, value V) {
	t.m.Store(key, value)
}

func (t *syncMapWrapper[K, V]) Load(key K) (V, bool) {
	value, ok := t.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), ok
}

func (t *syncMapWrapper[K, V]) Delete(key K) {
	t.m.Delete(key)
}

func (t *syncMapWrapper[K, V]) Range(f func(key K, value V) bool) {
	t.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}
