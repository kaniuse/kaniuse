package syncmap

import "sync"

type TypedSyncMap[K any, V any] interface {
	Load(key K) (V, bool)
	Store(key K, value V)
	Delete(key K)
	Range(f func(key K, value V) bool)
}

var _ TypedSyncMap[any, any] = (*SyncMapWrapper[any, any])(nil)

type SyncMapWrapper[K any, V any] struct {
	m *sync.Map
}

func NewSyncMapWrapper[K any, V any]() *SyncMapWrapper[K, V] {
	return &SyncMapWrapper[K, V]{
		m: &sync.Map{},
	}
}

func (t *SyncMapWrapper[K, V]) Store(key K, value V) {
	t.m.Store(key, value)
}

func (t *SyncMapWrapper[K, V]) Load(key K) (V, bool) {
	value, ok := t.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), ok
}

func (t *SyncMapWrapper[K, V]) Delete(key K) {
	t.m.Delete(key)
}

func (t *SyncMapWrapper[K, V]) Range(f func(key K, value V) bool) {
	t.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}
