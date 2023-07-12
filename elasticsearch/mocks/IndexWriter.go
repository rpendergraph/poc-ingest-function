// Code generated by mockery v2.31.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IndexWriter is an autogenerated mock type for the IndexWriter type
type IndexWriter struct {
	mock.Mock
}

// WriteToIndex provides a mock function with given fields: ctx, items
func (_m *IndexWriter) WriteToIndex(ctx context.Context, items []interface{}) error {
	ret := _m.Called(ctx, items)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}) error); ok {
		r0 = rf(ctx, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIndexWriter creates a new instance of IndexWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIndexWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *IndexWriter {
	mock := &IndexWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
