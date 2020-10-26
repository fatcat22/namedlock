package namedlock

import (
	"errors"
	"sync"
)

var (
	ErrLockNotExist = errors.New("lock does not exist")
)

var globalLock = NamedLock{
	locks: make(map[string]*sync.RWMutex),
}

type NamedLock struct {
	locks     map[string]*sync.RWMutex
	locksLock sync.RWMutex
}

func NewNamedLock() *NamedLock {
	return &NamedLock{
		locks: make(map[string]*sync.RWMutex),
	}
}

func Lock(key string)    { globalLock.Lock(key) }
func Unlock(key string)  { globalLock.Unlock(key) }
func RLock(key string)   { globalLock.RLock(key) }
func RUnlock(key string) { globalLock.RUnlock(key) }
func Delete(key string)  { globalLock.Delete(key) }

func (nlk *NamedLock) Lock(key string) {
	l := nlk.getOrNewLock(key)

	l.Lock()
}

func (nlk *NamedLock) Unlock(key string) {
	l := nlk.getLock(key)
	if l == nil {
		return
	}
	l.Unlock()
}

func (nlk *NamedLock) RLock(key string) {
	l := nlk.getOrNewLock(key)
	l.RLock()
}

func (nlk *NamedLock) RUnlock(key string) {
	l := nlk.getLock(key)
	if l == nil {
		return
	}
	l.RUnlock()
}

func (nlk *NamedLock) Delete(key string) {
	nlk.locksLock.Lock()
	defer nlk.locksLock.Unlock()
	delete(nlk.locks, key)
}

func (nlk *NamedLock) getOrNewLock(key string) *sync.RWMutex {
	if lk := nlk.getLock(key); lk != nil {
		return lk
	}
	return nlk.newLock(key)
}

func (nlk *NamedLock) getLock(key string) *sync.RWMutex {
	nlk.locksLock.RLock()
	defer nlk.locksLock.RUnlock()

	if l, ok := nlk.locks[key]; ok {
		return l
	}
	return nil
}

func (nlk *NamedLock) newLock(key string) *sync.RWMutex {
	nlk.locksLock.Lock()
	defer nlk.locksLock.Unlock()

	if lk, ok := nlk.locks[key]; ok {
		return lk
	}

	newL := new(sync.RWMutex)
	nlk.locks[key] = newL
	return newL
}
