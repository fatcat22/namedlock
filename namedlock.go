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

func Create(key string)      { globalLock.Create(key) }
func Delete(key string)      { globalLock.Delete(key) }
func Lock(key string) error  { return globalLock.Lock(key) }
func Unlock(key string)      { globalLock.Unlock(key) }
func RLock(key string) error { return globalLock.RLock(key) }
func RUnlock(key string)     { globalLock.RUnlock(key) }

func (nlk *NamedLock) Create(key string) {
	lk := new(sync.RWMutex)

	nlk.locksLock.Lock()
	defer nlk.locksLock.Unlock()
	if _, ok := nlk.locks[key]; !ok {
		nlk.locks[key] = lk
	}
}

func (nlk *NamedLock) Delete(key string) {
	nlk.locksLock.Lock()
	defer nlk.locksLock.Unlock()
	delete(nlk.locks, key)
}

func (nlk *NamedLock) Lock(key string) error {
	l, ok := nlk.getLock(key)
	if !ok {
		return ErrLockNotExist
	}

	l.Lock()
	return nil
}

func (nlk *NamedLock) Unlock(key string) {
	l, ok := nlk.getLock(key)
	if !ok {
		return
	}
	l.Unlock()
}

func (nlk *NamedLock) RLock(key string) error {
	l, ok := nlk.getLock(key)
	if !ok {
		return ErrLockNotExist
	}
	l.RLock()
	return nil
}

func (nlk *NamedLock) RUnlock(key string) {
	l, ok := nlk.getLock(key)
	if !ok {
		return
	}
	l.RUnlock()
}

func (nlk *NamedLock) getLock(key string) (*sync.RWMutex, bool) {
	nlk.locksLock.RLock()
	defer nlk.locksLock.RUnlock()

	l, ok := nlk.locks[key]
	return l, ok
}
