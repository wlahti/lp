// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/wlahti/lp/server"
)

type NameGetter struct {
	GetNameStub        func(int) (string, error)
	getNameMutex       sync.RWMutex
	getNameArgsForCall []struct {
		arg1 int
	}
	getNameReturns struct {
		result1 string
		result2 error
	}
	getNameReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *NameGetter) GetName(arg1 int) (string, error) {
	fake.getNameMutex.Lock()
	ret, specificReturn := fake.getNameReturnsOnCall[len(fake.getNameArgsForCall)]
	fake.getNameArgsForCall = append(fake.getNameArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("GetName", []interface{}{arg1})
	fake.getNameMutex.Unlock()
	if fake.GetNameStub != nil {
		return fake.GetNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getNameReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NameGetter) GetNameCallCount() int {
	fake.getNameMutex.RLock()
	defer fake.getNameMutex.RUnlock()
	return len(fake.getNameArgsForCall)
}

func (fake *NameGetter) GetNameCalls(stub func(int) (string, error)) {
	fake.getNameMutex.Lock()
	defer fake.getNameMutex.Unlock()
	fake.GetNameStub = stub
}

func (fake *NameGetter) GetNameArgsForCall(i int) int {
	fake.getNameMutex.RLock()
	defer fake.getNameMutex.RUnlock()
	argsForCall := fake.getNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NameGetter) GetNameReturns(result1 string, result2 error) {
	fake.getNameMutex.Lock()
	defer fake.getNameMutex.Unlock()
	fake.GetNameStub = nil
	fake.getNameReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *NameGetter) GetNameReturnsOnCall(i int, result1 string, result2 error) {
	fake.getNameMutex.Lock()
	defer fake.getNameMutex.Unlock()
	fake.GetNameStub = nil
	if fake.getNameReturnsOnCall == nil {
		fake.getNameReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getNameReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *NameGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getNameMutex.RLock()
	defer fake.getNameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *NameGetter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ server.NameGetter = new(NameGetter)
