package deathtax

import (
	"log"
	"sync"
)

// PooledFactory is a factory that holds between min and max
// hot instances.
type PooledFactory struct {
	mux         sync.RWMutex
	instances   []*Session
	hotSpares   int
	initialized bool
}

// NewPooledFactory creates a PooledFactory instance with the
// specified min and max hot spares
func NewPooledFactory(size int) *PooledFactory {
	pf := &PooledFactory{
		mux:         sync.RWMutex{},
		instances:   []*Session{},
		hotSpares:   size,
		initialized: false,
	}

	pf.addInstance()

	return pf
}

// Initialized returns true when all instances in the pool are ready
// after initial create. This always returns true once the initial
// processes are booted
func (pf *PooledFactory) Initialized() bool {
	log.Println("Checking Inititalized")
	for i, inst := range pf.instances {
		log.Printf("%d: %t\n", i, inst.IsReady())
	}

	pf.mux.RLock()
	if pf.initialized {
		return true
	}

	for _, s := range pf.instances {
		if !s.IsReady() {
			return false
		}
	}
	pf.mux.RUnlock()

	pf.mux.Lock()
	pf.initialized = true
	pf.mux.Unlock()

	return true
}

// GetInstance returns an instance of Session
func (pf *PooledFactory) GetInstance() *Session {
	var session *Session

	pf.mux.Lock()
	if len(pf.instances) > 0 {
		session, pf.instances = pf.instances[0], pf.instances[1:]
	} else {
		// Start new session if pool is empty
		session = NewSession()
	}
	pf.mux.Unlock()

	pf.addInstance()

	return session
}

func (pf *PooledFactory) addInstance() {
	pf.mux.RLock()
	numInstances := len(pf.instances)
	pf.mux.RUnlock()

	if numInstances < pf.hotSpares {
		// Fill to atleast min
		for i := 0; i < pf.hotSpares-numInstances; i++ {
			pf.mux.Lock()
			pf.instances = append(pf.instances, NewSession())
			pf.mux.Unlock()
		}
	}
}
