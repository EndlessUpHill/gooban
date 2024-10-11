package gooban

import "sync"

const DefaultRegistrySize = 100

type functionRegistry struct {
	sync.RWMutex
	functions map[string]func(args JobArgs) JobResult
}

func NewFunctionRegistry(size ...int) *functionRegistry {
	registrySize := DefaultRegistrySize
	if len(size) > 0 {
		registrySize = size[0]
	}

	return &functionRegistry{
		functions: make(map[string]func(args JobArgs) JobResult, registrySize),
	}
}

func (fr *functionRegistry) RegisterFunction(name string, actor func(args JobArgs) JobResult) {
	fr.Lock()
	defer fr.Unlock()
	fr.functions[name] = actor
}

func (ar *functionRegistry) GetActor(id string) (func(args JobArgs) JobResult, bool) {
	ar.RLock()
	defer ar.RUnlock()
	function, exists := ar.functions[id]
	return function, exists
}