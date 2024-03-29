// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// FnWithContextAndError is an autogenerated mock type for the FnWithContextAndError type
type FnWithContextAndError struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *FnWithContextAndError) Execute(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFnWithContextAndError creates a new instance of FnWithContextAndError. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFnWithContextAndError(t interface {
	mock.TestingT
	Cleanup(func())
}) *FnWithContextAndError {
	mock := &FnWithContextAndError{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
