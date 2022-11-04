// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	discovery "github.com/barcostreams/barco/internal/discovery"
	mock "github.com/stretchr/testify/mock"

	types "github.com/barcostreams/barco/internal/types"

	uuid "github.com/google/uuid"
)

// Discoverer is an autogenerated mock type for the Discoverer type
type Discoverer struct {
	mock.Mock
}

// Brokers provides a mock function with given fields:
func (_m *Discoverer) Brokers() []types.BrokerInfo {
	ret := _m.Called()

	var r0 []types.BrokerInfo
	if rf, ok := ret.Get(0).(func() []types.BrokerInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.BrokerInfo)
		}
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Discoverer) Close() {
	_m.Called()
}

// CurrentOrPastBroker provides a mock function with given fields: ordinal
func (_m *Discoverer) CurrentOrPastBroker(ordinal int) *types.BrokerInfo {
	ret := _m.Called(ordinal)

	var r0 *types.BrokerInfo
	if rf, ok := ret.Get(0).(func(int) *types.BrokerInfo); ok {
		r0 = rf(ordinal)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BrokerInfo)
		}
	}

	return r0
}

// Generation provides a mock function with given fields: token
func (_m *Discoverer) Generation(token types.Token) *types.Generation {
	ret := _m.Called(token)

	var r0 *types.Generation
	if rf, ok := ret.Get(0).(func(types.Token) *types.Generation); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Generation)
		}
	}

	return r0
}

// GenerationInfo provides a mock function with given fields: id
func (_m *Discoverer) GenerationInfo(id types.GenId) *types.Generation {
	ret := _m.Called(id)

	var r0 *types.Generation
	if rf, ok := ret.Get(0).(func(types.GenId) *types.Generation); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Generation)
		}
	}

	return r0
}

// GenerationProposed provides a mock function with given fields: token
func (_m *Discoverer) GenerationProposed(token types.Token) (*types.Generation, *types.Generation) {
	ret := _m.Called(token)

	var r0 *types.Generation
	if rf, ok := ret.Get(0).(func(types.Token) *types.Generation); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Generation)
		}
	}

	var r1 *types.Generation
	if rf, ok := ret.Get(1).(func(types.Token) *types.Generation); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*types.Generation)
		}
	}

	return r0, r1
}

// GetTokenHistory provides a mock function with given fields: token, clusterSize
func (_m *Discoverer) GetTokenHistory(token types.Token, clusterSize int) (*types.Generation, error) {
	ret := _m.Called(token, clusterSize)

	var r0 *types.Generation
	if rf, ok := ret.Get(0).(func(types.Token, int) *types.Generation); ok {
		r0 = rf(token, clusterSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Generation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Token, int) error); ok {
		r1 = rf(token, clusterSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasTokenHistory provides a mock function with given fields: token, clusterSize
func (_m *Discoverer) HasTokenHistory(token types.Token, clusterSize int) (bool, error) {
	ret := _m.Called(token, clusterSize)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Token, int) bool); ok {
		r0 = rf(token, clusterSize)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Token, int) error); ok {
		r1 = rf(token, clusterSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields:
func (_m *Discoverer) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsTokenInRange provides a mock function with given fields: token
func (_m *Discoverer) IsTokenInRange(token types.Token) bool {
	ret := _m.Called(token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Token) bool); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Leader provides a mock function with given fields: partitionKey
func (_m *Discoverer) Leader(partitionKey string) types.ReplicationInfo {
	ret := _m.Called(partitionKey)

	var r0 types.ReplicationInfo
	if rf, ok := ret.Get(0).(func(string) types.ReplicationInfo); ok {
		r0 = rf(partitionKey)
	} else {
		r0 = ret.Get(0).(types.ReplicationInfo)
	}

	return r0
}

// LocalInfo provides a mock function with given fields:
func (_m *Discoverer) LocalInfo() *types.BrokerInfo {
	ret := _m.Called()

	var r0 *types.BrokerInfo
	if rf, ok := ret.Get(0).(func() *types.BrokerInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BrokerInfo)
		}
	}

	return r0
}

// NextGeneration provides a mock function with given fields: id
func (_m *Discoverer) NextGeneration(id types.GenId) []types.Generation {
	ret := _m.Called(id)

	var r0 []types.Generation
	if rf, ok := ret.Get(0).(func(types.GenId) []types.Generation); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Generation)
		}
	}

	return r0
}

// ParentRanges provides a mock function with given fields: gen, indices
func (_m *Discoverer) ParentRanges(gen *types.Generation, indices []types.RangeIndex) []types.TokenRanges {
	ret := _m.Called(gen, indices)

	var r0 []types.TokenRanges
	if rf, ok := ret.Get(0).(func(*types.Generation, []types.RangeIndex) []types.TokenRanges); ok {
		r0 = rf(gen, indices)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.TokenRanges)
		}
	}

	return r0
}

// RegisterListener provides a mock function with given fields: l
func (_m *Discoverer) RegisterListener(l discovery.TopologyChangeListener) {
	_m.Called(l)
}

// RepairCommitted provides a mock function with given fields: gen
func (_m *Discoverer) RepairCommitted(gen *types.Generation) error {
	ret := _m.Called(gen)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Generation) error); ok {
		r0 = rf(gen)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetAsCommitted provides a mock function with given fields: token1, token2, tx, origin
func (_m *Discoverer) SetAsCommitted(token1 types.Token, token2 *types.Token, tx uuid.UUID, origin int) error {
	ret := _m.Called(token1, token2, tx, origin)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Token, *types.Token, uuid.UUID, int) error); ok {
		r0 = rf(token1, token2, tx, origin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetGenerationProposed provides a mock function with given fields: gen, gen2, expectedTx
func (_m *Discoverer) SetGenerationProposed(gen *types.Generation, gen2 *types.Generation, expectedTx *uuid.UUID) error {
	ret := _m.Called(gen, gen2, expectedTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Generation, *types.Generation, *uuid.UUID) error); ok {
		r0 = rf(gen, gen2, expectedTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Topology provides a mock function with given fields:
func (_m *Discoverer) Topology() *types.TopologyInfo {
	ret := _m.Called()

	var r0 *types.TopologyInfo
	if rf, ok := ret.Get(0).(func() *types.TopologyInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.TopologyInfo)
		}
	}

	return r0
}

type mockConstructorTestingTNewDiscoverer interface {
	mock.TestingT
	Cleanup(func())
}

// NewDiscoverer creates a new instance of Discoverer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDiscoverer(t mockConstructorTestingTNewDiscoverer) *Discoverer {
	mock := &Discoverer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
