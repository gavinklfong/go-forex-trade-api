// Code generated by mockery v2.53.3. DO NOT EDIT.

package apiclient

import (
	model "github.com/gavinklfong/go-forex-trade-api/apiclient/model"
	mock "github.com/stretchr/testify/mock"
)

// MockForexApiClient is an autogenerated mock type for the ForexApiClient type
type MockForexApiClient struct {
	mock.Mock
}

type MockForexApiClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockForexApiClient) EXPECT() *MockForexApiClient_Expecter {
	return &MockForexApiClient_Expecter{mock: &_m.Mock}
}

// GetRateByBaseCurrency provides a mock function with given fields: base
func (_m *MockForexApiClient) GetRateByBaseCurrency(base string) (*model.ForexRateResponse, error) {
	ret := _m.Called(base)

	if len(ret) == 0 {
		panic("no return value specified for GetRateByBaseCurrency")
	}

	var r0 *model.ForexRateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.ForexRateResponse, error)); ok {
		return rf(base)
	}
	if rf, ok := ret.Get(0).(func(string) *model.ForexRateResponse); ok {
		r0 = rf(base)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ForexRateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(base)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockForexApiClient_GetRateByBaseCurrency_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRateByBaseCurrency'
type MockForexApiClient_GetRateByBaseCurrency_Call struct {
	*mock.Call
}

// GetRateByBaseCurrency is a helper method to define mock.On call
//   - base string
func (_e *MockForexApiClient_Expecter) GetRateByBaseCurrency(base interface{}) *MockForexApiClient_GetRateByBaseCurrency_Call {
	return &MockForexApiClient_GetRateByBaseCurrency_Call{Call: _e.mock.On("GetRateByBaseCurrency", base)}
}

func (_c *MockForexApiClient_GetRateByBaseCurrency_Call) Run(run func(base string)) *MockForexApiClient_GetRateByBaseCurrency_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockForexApiClient_GetRateByBaseCurrency_Call) Return(_a0 *model.ForexRateResponse, _a1 error) *MockForexApiClient_GetRateByBaseCurrency_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockForexApiClient_GetRateByBaseCurrency_Call) RunAndReturn(run func(string) (*model.ForexRateResponse, error)) *MockForexApiClient_GetRateByBaseCurrency_Call {
	_c.Call.Return(run)
	return _c
}

// GetRateByCurrencyPair provides a mock function with given fields: base, counter
func (_m *MockForexApiClient) GetRateByCurrencyPair(base string, counter string) (*model.ForexRateResponse, error) {
	ret := _m.Called(base, counter)

	if len(ret) == 0 {
		panic("no return value specified for GetRateByCurrencyPair")
	}

	var r0 *model.ForexRateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*model.ForexRateResponse, error)); ok {
		return rf(base, counter)
	}
	if rf, ok := ret.Get(0).(func(string, string) *model.ForexRateResponse); ok {
		r0 = rf(base, counter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ForexRateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(base, counter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockForexApiClient_GetRateByCurrencyPair_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRateByCurrencyPair'
type MockForexApiClient_GetRateByCurrencyPair_Call struct {
	*mock.Call
}

// GetRateByCurrencyPair is a helper method to define mock.On call
//   - base string
//   - counter string
func (_e *MockForexApiClient_Expecter) GetRateByCurrencyPair(base interface{}, counter interface{}) *MockForexApiClient_GetRateByCurrencyPair_Call {
	return &MockForexApiClient_GetRateByCurrencyPair_Call{Call: _e.mock.On("GetRateByCurrencyPair", base, counter)}
}

func (_c *MockForexApiClient_GetRateByCurrencyPair_Call) Run(run func(base string, counter string)) *MockForexApiClient_GetRateByCurrencyPair_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockForexApiClient_GetRateByCurrencyPair_Call) Return(_a0 *model.ForexRateResponse, _a1 error) *MockForexApiClient_GetRateByCurrencyPair_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockForexApiClient_GetRateByCurrencyPair_Call) RunAndReturn(run func(string, string) (*model.ForexRateResponse, error)) *MockForexApiClient_GetRateByCurrencyPair_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockForexApiClient creates a new instance of MockForexApiClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockForexApiClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockForexApiClient {
	mock := &MockForexApiClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
