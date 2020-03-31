package deathtax

import "sync"

// PooledFactory is a factory that holds between min and max
// hot instances.
type PooledFactory struct {
	instances    []*Session
	minInstances int
	maxInstances int
	mux          sync.RWMutex
}

// NewPooledFactory creates a PooledFactory instance with the
// specified min and max hot spares
func NewPooledFactory(min, max int) *PooledFactory {
	pf := &PooledFactory{
		instances:    []*Session{},
		minInstances: min,
		maxInstances: max,
		mux:          sync.RWMutex{},
	}

	pf.addInstance()

	return pf
}

func (pf *PooledFactory) addInstance() {
	pf.mux.RLock()
	numInstances := len(pf.instances)
	pf.mux.RUnlock()

	if numInstances < pf.minInstances {
		// Fill to atleast min
		for i := 0; i < pf.minInstances-numInstances; i++ {
			pf.mux.Lock()
			pf.instances = append(pf.instances, NewSession())
			pf.mux.Unlock()
		}
	} else if numInstances < pf.maxInstances {
		// Dont go over max
		pf.mux.Lock()
		pf.instances = append(pf.instances, NewSession())
		pf.mux.Unlock()
	}
}

// GetInstance returns an instance of Session
func (pf *PooledFactory) GetInstance() *Session {
	var session *Session
	pf.mux.Lock()
	session, pf.instances = pf.instances[0], pf.instances[1:]
	pf.mux.Unlock()

	pf.addInstance()

	return session
}
