package core

import (
	rpc "github.com/tendermint/tendermint/rpc/lib/server"

	data "github.com/tendermint/go-wire/data"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpc "github.com/tendermint/tendermint/rpc/lib/server"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/tendermint/tendermint/types"
)

// TODO: better system than "unsafe" prefix
var Routes = map[string]*rpc.RPCFunc{
	// subscribe/unsubscribe are reserved for websocket events.
	"subscribe":   rpc.NewWSRPCFunc(Subscribe, "event"),
	"unsubscribe": rpc.NewWSRPCFunc(Unsubscribe, "event"),

	// info API
	"status":               rpc.NewRPCFunc(Status, ""),
	"net_info":             rpc.NewRPCFunc(NetInfo, ""),
	"blockchain":           rpc.NewRPCFunc(BlockchainInfo, "minHeight,maxHeight"),
	"genesis":              rpc.NewRPCFunc(Genesis, ""),
	"block":                rpc.NewRPCFunc(Block, "height"),
	"commit":               rpc.NewRPCFunc(Commit, "height"),
	"tx":                   rpc.NewRPCFunc(Tx, "hash,prove"),
	"validators":           rpc.NewRPCFunc(Validators, ""),
	"dump_consensus_state": rpc.NewRPCFunc(DumpConsensusState, ""),
	"unconfirmed_txs":      rpc.NewRPCFunc(UnconfirmedTxs, ""),
	"num_unconfirmed_txs":  rpc.NewRPCFunc(NumUnconfirmedTxs, ""),

	// broadcast API
	"broadcast_tx_commit": rpc.NewRPCFunc(BroadcastTxCommit, "tx"),
	"broadcast_tx_sync":   rpc.NewRPCFunc(BroadcastTxSync, "tx"),
	"broadcast_tx_async":  rpc.NewRPCFunc(BroadcastTxAsync, "tx"),

	// abci API
	"abci_query": rpc.NewRPCFunc(ABCIQuery, "path,data,prove"),
	"abci_info":  rpc.NewRPCFunc(ABCIInfo, ""),

	// control API
	"dial_seeds":           rpc.NewRPCFunc(UnsafeDialSeeds, "seeds"),
	"unsafe_flush_mempool": rpc.NewRPCFunc(UnsafeFlushMempool, ""),

	// profiler API
	"unsafe_start_cpu_profiler": rpc.NewRPCFunc(UnsafeStartCPUProfiler, "filename"),
	"unsafe_stop_cpu_profiler":  rpc.NewRPCFunc(UnsafeStopCPUProfiler, ""),
	"unsafe_write_heap_profile": rpc.NewRPCFunc(UnsafeWriteHeapProfile, "filename"),
}
