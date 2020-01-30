package gintunnel

import "sync"

type ForwardMap struct {
	m  map[string]string
	rw sync.RWMutex
}

func NewForwardMap() ForwardMap {
	return ForwardMap{m: make(map[string]string)}
}
func (f *ForwardMap) get(key string) string {
	f.rw.RLock()
	defer f.rw.RUnlock()
	return f.m[key]
}
func (f *ForwardMap) set(key string, value string) bool {
	f.rw.Lock()
	defer f.rw.Unlock()
	if f.m[key] != "" {
		return false
	}
	f.m[key] = value
	return true
}
func (f *ForwardMap) remove(key string) {
	f.rw.Lock()
	defer f.rw.Unlock()
	delete(f.m, key)
}
