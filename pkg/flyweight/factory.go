package flyweight

import (
	"runtime"
	"sync"
	"weak"
)

type Factory[K comparable, V any] interface {
	Get(K, func() (*V, error)) (*V, error)
}

func NewFactory[K comparable, V any]() Factory[K, V] {
	return &factory[K, V]{}
}

type factory[K comparable, V any] struct {
	cache sync.Map
}

func (f *factory[K, V]) Get(key K, builder func() (*V, error)) (*V, error) {
	var newValue *V
	for {
		cached, ok := f.cache.Load(key)
		if !ok {
			// No cached object found. Create the new object.
			if newValue == nil {
				var err error
				newValue, err = builder()
				if err != nil {
					return nil, err
				}
			}
			// Try to install the new object into cache
			wp := weak.Make(newValue)
			var loaded bool
			cached, loaded = f.cache.LoadOrStore(key, wp)
			if !loaded {
				// New object was installed into cache.
				runtime.AddCleanup(newValue, func(key K) {
					// Only delete if the weak pointer is equal.
					// If it's not, someone else already deleted the entry and, probably, installed a new one.
					f.cache.CompareAndDelete(key, wp)
				}, key)
				return newValue, nil
			}
			// Someone has installed the cache entry before us.
		}
		// Check if our cache entry is valid.
		if value := cached.(weak.Pointer[V]).Value(); value != nil {
			return value, nil
		}
		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		f.cache.CompareAndDelete(key, cached)
	}
}
