// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/KevinJCross/cf-test-helpers/v2/internal"
)

type FakeRedactor struct {
	RedactStub        func(toRedact string) string
	redactMutex       sync.RWMutex
	redactArgsForCall []struct {
		toRedact string
	}
	redactReturns struct {
		result1 string
	}
}

func (fake *FakeRedactor) Redact(toRedact string) string {
	fake.redactMutex.Lock()
	fake.redactArgsForCall = append(fake.redactArgsForCall, struct {
		toRedact string
	}{toRedact})
	fake.redactMutex.Unlock()
	if fake.RedactStub != nil {
		return fake.RedactStub(toRedact)
	} else {
		return fake.redactReturns.result1
	}
}

func (fake *FakeRedactor) RedactCallCount() int {
	fake.redactMutex.RLock()
	defer fake.redactMutex.RUnlock()
	return len(fake.redactArgsForCall)
}

func (fake *FakeRedactor) RedactArgsForCall(i int) string {
	fake.redactMutex.RLock()
	defer fake.redactMutex.RUnlock()
	return fake.redactArgsForCall[i].toRedact
}

func (fake *FakeRedactor) RedactReturns(result1 string) {
	fake.RedactStub = nil
	fake.redactReturns = struct {
		result1 string
	}{result1}
}

var _ internal.Redactor = new(FakeRedactor)
