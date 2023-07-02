// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// AliasGenerator is an autogenerated mock type for the AliasGenerator type
type AliasGenerator struct {
	mock.Mock
}

// NewAlias provides a mock function with given fields: size
func (_m *AliasGenerator) NewAlias(size int) string {
	ret := _m.Called(size)

	var r0 string
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(size)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewAliasGenerator interface {
	mock.TestingT
	Cleanup(func())
}

// NewAliasGenerator creates a new instance of AliasGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAliasGenerator(t mockConstructorTestingTNewAliasGenerator) *AliasGenerator {
	mock := &AliasGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
