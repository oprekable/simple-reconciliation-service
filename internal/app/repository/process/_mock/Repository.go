// Code generated by mockery v2.50.0. DO NOT EDIT.

package _mock

import (
	context "context"
	banks "simple-reconciliation-service/internal/pkg/reconcile/parser/banks"

	mock "github.com/stretchr/testify/mock"

	process "simple-reconciliation-service/internal/app/repository/process"

	systems "simple-reconciliation-service/internal/pkg/reconcile/parser/systems"

	time "time"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *Repository) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateReconciliationMap provides a mock function with given fields: ctx, minAmount, maxAmount
func (_m *Repository) GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) error {
	ret := _m.Called(ctx, minAmount, maxAmount)

	if len(ret) == 0 {
		panic("no return value specified for GenerateReconciliationMap")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, float64, float64) error); ok {
		r0 = rf(ctx, minAmount, maxAmount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetReconciliationSummary provides a mock function with given fields: ctx
func (_m *Repository) GetReconciliationSummary(ctx context.Context) (process.ReconciliationSummary, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetReconciliationSummary")
	}

	var r0 process.ReconciliationSummary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (process.ReconciliationSummary, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) process.ReconciliationSummary); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(process.ReconciliationSummary)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImportBankTrx provides a mock function with given fields: ctx, data
func (_m *Repository) ImportBankTrx(ctx context.Context, data []*banks.BankTrxData) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for ImportBankTrx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*banks.BankTrxData) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ImportSystemTrx provides a mock function with given fields: ctx, data
func (_m *Repository) ImportSystemTrx(ctx context.Context, data []*systems.SystemTrxData) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for ImportSystemTrx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*systems.SystemTrxData) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Post provides a mock function with given fields: ctx
func (_m *Repository) Post(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Post")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pre provides a mock function with given fields: ctx, listBank, startDate, toDate
func (_m *Repository) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) error {
	ret := _m.Called(ctx, listBank, startDate, toDate)

	if len(ret) == 0 {
		panic("no return value specified for Pre")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, time.Time, time.Time) error); ok {
		r0 = rf(ctx, listBank, startDate, toDate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}