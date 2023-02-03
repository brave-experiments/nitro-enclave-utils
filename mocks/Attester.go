// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Attester is an autogenerated mock type for the Attester type
type Attester struct {
	mock.Mock
}

// GetAttestDoc provides a mock function with given fields: nonce, publicKey, userData
func (_m *Attester) Attest(nonce []byte, publicKey []byte, userData []byte) ([]byte, error) {
	ret := _m.Called(nonce, publicKey, userData)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, []byte, []byte) []byte); ok {
		r0 = rf(nonce, publicKey, userData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte, []byte) error); ok {
		r1 = rf(nonce, publicKey, userData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAttester interface {
	mock.TestingT
	Cleanup(func())
}

// NewAttester creates a new instance of Attester. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAttester(t mockConstructorTestingTNewAttester) *Attester {
	mock := &Attester{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
