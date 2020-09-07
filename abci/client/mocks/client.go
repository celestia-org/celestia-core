// Code generated by mockery v1.1.1. DO NOT EDIT.

package mocks

import (
	abcicli "github.com/lazyledger/lazyledger-core/abci/client"
	log "github.com/lazyledger/lazyledger-core/libs/log"

	mock "github.com/stretchr/testify/mock"

	types "github.com/lazyledger/lazyledger-core/abci/types"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// ApplySnapshotChunkAsync provides a mock function with given fields: _a0
func (_m *Client) ApplySnapshotChunkAsync(_a0 types.RequestApplySnapshotChunk) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestApplySnapshotChunk) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// ApplySnapshotChunkSync provides a mock function with given fields: _a0
func (_m *Client) ApplySnapshotChunkSync(_a0 types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseApplySnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestApplySnapshotChunk) *types.ResponseApplySnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseApplySnapshotChunk)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestApplySnapshotChunk) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeginBlockAsync provides a mock function with given fields: _a0
func (_m *Client) BeginBlockAsync(_a0 types.RequestBeginBlock) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestBeginBlock) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// BeginBlockSync provides a mock function with given fields: _a0
func (_m *Client) BeginBlockSync(_a0 types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseBeginBlock
	if rf, ok := ret.Get(0).(func(types.RequestBeginBlock) *types.ResponseBeginBlock); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseBeginBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestBeginBlock) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckTxAsync provides a mock function with given fields: _a0
func (_m *Client) CheckTxAsync(_a0 types.RequestCheckTx) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestCheckTx) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// CheckTxSync provides a mock function with given fields: _a0
func (_m *Client) CheckTxSync(_a0 types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseCheckTx
	if rf, ok := ret.Get(0).(func(types.RequestCheckTx) *types.ResponseCheckTx); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseCheckTx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestCheckTx) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommitAsync provides a mock function with given fields:
func (_m *Client) CommitAsync() *abcicli.ReqRes {
	ret := _m.Called()

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func() *abcicli.ReqRes); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// CommitSync provides a mock function with given fields:
func (_m *Client) CommitSync() (*types.ResponseCommit, error) {
	ret := _m.Called()

	var r0 *types.ResponseCommit
	if rf, ok := ret.Get(0).(func() *types.ResponseCommit); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseCommit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeliverTxAsync provides a mock function with given fields: _a0
func (_m *Client) DeliverTxAsync(_a0 types.RequestDeliverTx) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestDeliverTx) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// DeliverTxSync provides a mock function with given fields: _a0
func (_m *Client) DeliverTxSync(_a0 types.RequestDeliverTx) (*types.ResponseDeliverTx, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseDeliverTx
	if rf, ok := ret.Get(0).(func(types.RequestDeliverTx) *types.ResponseDeliverTx); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseDeliverTx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestDeliverTx) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EchoAsync provides a mock function with given fields: msg
func (_m *Client) EchoAsync(msg string) *abcicli.ReqRes {
	ret := _m.Called(msg)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(string) *abcicli.ReqRes); ok {
		r0 = rf(msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// EchoSync provides a mock function with given fields: msg
func (_m *Client) EchoSync(msg string) (*types.ResponseEcho, error) {
	ret := _m.Called(msg)

	var r0 *types.ResponseEcho
	if rf, ok := ret.Get(0).(func(string) *types.ResponseEcho); ok {
		r0 = rf(msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseEcho)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EndBlockAsync provides a mock function with given fields: _a0
func (_m *Client) EndBlockAsync(_a0 types.RequestEndBlock) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestEndBlock) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// EndBlockSync provides a mock function with given fields: _a0
func (_m *Client) EndBlockSync(_a0 types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseEndBlock
	if rf, ok := ret.Get(0).(func(types.RequestEndBlock) *types.ResponseEndBlock); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseEndBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestEndBlock) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Error provides a mock function with given fields:
func (_m *Client) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FlushAsync provides a mock function with given fields:
func (_m *Client) FlushAsync() *abcicli.ReqRes {
	ret := _m.Called()

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func() *abcicli.ReqRes); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// FlushSync provides a mock function with given fields:
func (_m *Client) FlushSync() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InfoAsync provides a mock function with given fields: _a0
func (_m *Client) InfoAsync(_a0 types.RequestInfo) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestInfo) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// InfoSync provides a mock function with given fields: _a0
func (_m *Client) InfoSync(_a0 types.RequestInfo) (*types.ResponseInfo, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseInfo
	if rf, ok := ret.Get(0).(func(types.RequestInfo) *types.ResponseInfo); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestInfo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitChainAsync provides a mock function with given fields: _a0
func (_m *Client) InitChainAsync(_a0 types.RequestInitChain) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestInitChain) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// InitChainSync provides a mock function with given fields: _a0
func (_m *Client) InitChainSync(_a0 types.RequestInitChain) (*types.ResponseInitChain, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseInitChain
	if rf, ok := ret.Get(0).(func(types.RequestInitChain) *types.ResponseInitChain); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseInitChain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestInitChain) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsRunning provides a mock function with given fields:
func (_m *Client) IsRunning() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ListSnapshotsAsync provides a mock function with given fields: _a0
func (_m *Client) ListSnapshotsAsync(_a0 types.RequestListSnapshots) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestListSnapshots) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// ListSnapshotsSync provides a mock function with given fields: _a0
func (_m *Client) ListSnapshotsSync(_a0 types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseListSnapshots
	if rf, ok := ret.Get(0).(func(types.RequestListSnapshots) *types.ResponseListSnapshots); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseListSnapshots)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestListSnapshots) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadSnapshotChunkAsync provides a mock function with given fields: _a0
func (_m *Client) LoadSnapshotChunkAsync(_a0 types.RequestLoadSnapshotChunk) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestLoadSnapshotChunk) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// LoadSnapshotChunkSync provides a mock function with given fields: _a0
func (_m *Client) LoadSnapshotChunkSync(_a0 types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseLoadSnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestLoadSnapshotChunk) *types.ResponseLoadSnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseLoadSnapshotChunk)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestLoadSnapshotChunk) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OfferSnapshotAsync provides a mock function with given fields: _a0
func (_m *Client) OfferSnapshotAsync(_a0 types.RequestOfferSnapshot) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestOfferSnapshot) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// OfferSnapshotSync provides a mock function with given fields: _a0
func (_m *Client) OfferSnapshotSync(_a0 types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseOfferSnapshot
	if rf, ok := ret.Get(0).(func(types.RequestOfferSnapshot) *types.ResponseOfferSnapshot); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseOfferSnapshot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestOfferSnapshot) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OnReset provides a mock function with given fields:
func (_m *Client) OnReset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnStart provides a mock function with given fields:
func (_m *Client) OnStart() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnStop provides a mock function with given fields:
func (_m *Client) OnStop() {
	_m.Called()
}

// QueryAsync provides a mock function with given fields: _a0
func (_m *Client) QueryAsync(_a0 types.RequestQuery) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestQuery) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// QuerySync provides a mock function with given fields: _a0
func (_m *Client) QuerySync(_a0 types.RequestQuery) (*types.ResponseQuery, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseQuery
	if rf, ok := ret.Get(0).(func(types.RequestQuery) *types.ResponseQuery); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseQuery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestQuery) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Quit provides a mock function with given fields:
func (_m *Client) Quit() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *Client) Reset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLogger provides a mock function with given fields: _a0
func (_m *Client) SetLogger(_a0 log.Logger) {
	_m.Called(_a0)
}

// SetOptionAsync provides a mock function with given fields: _a0
func (_m *Client) SetOptionAsync(_a0 types.RequestSetOption) *abcicli.ReqRes {
	ret := _m.Called(_a0)

	var r0 *abcicli.ReqRes
	if rf, ok := ret.Get(0).(func(types.RequestSetOption) *abcicli.ReqRes); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	return r0
}

// SetOptionSync provides a mock function with given fields: _a0
func (_m *Client) SetOptionSync(_a0 types.RequestSetOption) (*types.ResponseSetOption, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseSetOption
	if rf, ok := ret.Get(0).(func(types.RequestSetOption) *types.ResponseSetOption); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseSetOption)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestSetOption) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetResponseCallback provides a mock function with given fields: _a0
func (_m *Client) SetResponseCallback(_a0 abcicli.Callback) {
	_m.Called(_a0)
}

// Start provides a mock function with given fields:
func (_m *Client) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *Client) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *Client) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
