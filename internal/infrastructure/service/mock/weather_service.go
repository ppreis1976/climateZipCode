// Code generated by mockery. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "climateZpCode/internal/business/model"
)

// MockWheaterService is an autogenerated mock type for the WheaterService type
type MockWheaterService struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *MockWheaterService) Get(_a0 context.Context, _a1 model.WeatherID) (*model.Weather, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Weather
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.WeatherID) (*model.Weather, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.WeatherID) *model.Weather); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Weather)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.WeatherID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockWheaterService creates a new instance of MockWheaterService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWheaterService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWheaterService {
	mock := &MockWheaterService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}