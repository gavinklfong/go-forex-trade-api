// Code generated by mockery v2.53.3. DO NOT EDIT.

package dao

import (
	model "github.com/gavinklfong/go-forex-trade-api/model"
	mock "github.com/stretchr/testify/mock"
)

// MockForexRateDao is an autogenerated mock type for the ForexRateDao type
type MockForexRateDao struct {
	mock.Mock
}

type MockForexRateDao_Expecter struct {
	mock *mock.Mock
}

func (_m *MockForexRateDao) EXPECT() *MockForexRateDao_Expecter {
	return &MockForexRateDao_Expecter{mock: &_m.Mock}
}

// FindByBookingRef provides a mock function with given fields: bookingRef
func (_m *MockForexRateDao) FindByBookingRef(bookingRef string) (*model.ForexRateBooking, error) {
	ret := _m.Called(bookingRef)

	if len(ret) == 0 {
		panic("no return value specified for FindByBookingRef")
	}

	var r0 *model.ForexRateBooking
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.ForexRateBooking, error)); ok {
		return rf(bookingRef)
	}
	if rf, ok := ret.Get(0).(func(string) *model.ForexRateBooking); ok {
		r0 = rf(bookingRef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ForexRateBooking)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bookingRef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockForexRateDao_FindByBookingRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByBookingRef'
type MockForexRateDao_FindByBookingRef_Call struct {
	*mock.Call
}

// FindByBookingRef is a helper method to define mock.On call
//   - bookingRef string
func (_e *MockForexRateDao_Expecter) FindByBookingRef(bookingRef interface{}) *MockForexRateDao_FindByBookingRef_Call {
	return &MockForexRateDao_FindByBookingRef_Call{Call: _e.mock.On("FindByBookingRef", bookingRef)}
}

func (_c *MockForexRateDao_FindByBookingRef_Call) Run(run func(bookingRef string)) *MockForexRateDao_FindByBookingRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockForexRateDao_FindByBookingRef_Call) Return(_a0 *model.ForexRateBooking, _a1 error) *MockForexRateDao_FindByBookingRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockForexRateDao_FindByBookingRef_Call) RunAndReturn(run func(string) (*model.ForexRateBooking, error)) *MockForexRateDao_FindByBookingRef_Call {
	_c.Call.Return(run)
	return _c
}

// FindByID provides a mock function with given fields: id
func (_m *MockForexRateDao) FindByID(id string) (*model.ForexRateBooking, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *model.ForexRateBooking
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.ForexRateBooking, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.ForexRateBooking); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ForexRateBooking)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockForexRateDao_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type MockForexRateDao_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - id string
func (_e *MockForexRateDao_Expecter) FindByID(id interface{}) *MockForexRateDao_FindByID_Call {
	return &MockForexRateDao_FindByID_Call{Call: _e.mock.On("FindByID", id)}
}

func (_c *MockForexRateDao_FindByID_Call) Run(run func(id string)) *MockForexRateDao_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockForexRateDao_FindByID_Call) Return(_a0 *model.ForexRateBooking, _a1 error) *MockForexRateDao_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockForexRateDao_FindByID_Call) RunAndReturn(run func(string) (*model.ForexRateBooking, error)) *MockForexRateDao_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: booking
func (_m *MockForexRateDao) Insert(booking *model.ForexRateBooking) (int64, error) {
	ret := _m.Called(booking)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.ForexRateBooking) (int64, error)); ok {
		return rf(booking)
	}
	if rf, ok := ret.Get(0).(func(*model.ForexRateBooking) int64); ok {
		r0 = rf(booking)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(*model.ForexRateBooking) error); ok {
		r1 = rf(booking)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockForexRateDao_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockForexRateDao_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - booking *model.ForexRateBooking
func (_e *MockForexRateDao_Expecter) Insert() *MockForexRateDao_Insert_Call {
	return &MockForexRateDao_Insert_Call{Call: _e.mock.On("Insert", mock.Anything)}
}

func (_c *MockForexRateDao_Insert_Call) Run(run func(booking *model.ForexRateBooking)) *MockForexRateDao_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.ForexRateBooking))
	})
	return _c
}

func (_c *MockForexRateDao_Insert_Call) Return(_a0 int64, _a1 error) *MockForexRateDao_Insert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockForexRateDao_Insert_Call) RunAndReturn(run func(*model.ForexRateBooking) (int64, error)) *MockForexRateDao_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockForexRateDao creates a new instance of MockForexRateDao. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockForexRateDao(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockForexRateDao {
	mock := &MockForexRateDao{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
