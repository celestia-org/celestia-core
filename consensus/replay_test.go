package consensus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/tendermint/abci/example/dummy"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/log"
)

func init() {
	config = ResetConfig("consensus_replay_test")
}

// These tests ensure we can always recover from failure at any part of the consensus process.
// There are two general failure scenarios: failure during consensus, and failure while applying the block.
// Only the latter interacts with the app and store,
// but the former has to deal with restrictions on re-use of priv_validator keys.
// The `WAL Tests` are for failures during the consensus;
// the `Handshake Tests` are for failures in applying the block.
// With the help of the WAL, we can recover from it all!

// NOTE: Files in this dir are generated by running the `build.sh` therein.
// It's a simple way to generate wals for a single block, or multiple blocks, with random transactions,
// and different part sizes. The output is not deterministic, and the stepChanges may need to be adjusted
// after running it (eg. sometimes small_block2 will have 5 block parts, sometimes 6).
// It should only have to be re-run if there is some breaking change to the consensus data structures (eg. blocks, votes)
// or to the behaviour of the app (eg. computes app hash differently)
var data_dir = path.Join(cmn.GoPath, "src/github.com/tendermint/tendermint/consensus", "test_data")

//------------------------------------------------------------------------------------------
// WAL Tests

// TODO: It would be better to verify explicitly which states we can recover from without the wal
// and which ones we need the wal for - then we'd also be able to only flush the
// wal writer when we need to, instead of with every message.

// the priv validator changes step at these lines for a block with 1 val and 1 part
var baseStepChanges = []int{3, 6, 8}

// test recovery from each line in each testCase
var testCases = []*testCase{
	newTestCase("empty_block", baseStepChanges),   // empty block (has 1 block part)
	newTestCase("small_block1", baseStepChanges),  // small block with txs in 1 block part
	newTestCase("small_block2", []int{3, 11, 13}), // small block with txs across 6 smaller block parts
}

type testCase struct {
	name    string
	log     string       //full cs wal
	stepMap map[int]int8 // map lines of log to privval step

	proposeLine   int
	prevoteLine   int
	precommitLine int
}

func newTestCase(name string, stepChanges []int) *testCase {
	if len(stepChanges) != 3 {
		panic(cmn.Fmt("a full wal has 3 step changes! Got array %v", stepChanges))
	}
	return &testCase{
		name:    name,
		log:     readWAL(path.Join(data_dir, name+".cswal")),
		stepMap: newMapFromChanges(stepChanges),

		proposeLine:   stepChanges[0],
		prevoteLine:   stepChanges[1],
		precommitLine: stepChanges[2],
	}
}

func newMapFromChanges(changes []int) map[int]int8 {
	changes = append(changes, changes[2]+1) // so we add the last step change to the map
	m := make(map[int]int8)
	var count int
	for changeNum, nextChange := range changes {
		for ; count < nextChange; count++ {
			m[count] = int8(changeNum)
		}
	}
	return m
}

func readWAL(p string) string {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func writeWAL(walMsgs string) string {
	tempDir := os.TempDir()
	walDir := path.Join(tempDir, "/wal"+cmn.RandStr(12))
	walFile := path.Join(walDir, "wal")
	// Create WAL directory
	err := cmn.EnsureDir(walDir, 0700)
	if err != nil {
		panic(err)
	}
	// Write the needed WAL to file
	err = cmn.WriteFile(walFile, []byte(walMsgs), 0600)
	if err != nil {
		panic(err)
	}
	return walFile
}

func waitForBlock(newBlockCh chan interface{}, thisCase *testCase, i int) {
	after := time.After(time.Second * 10)
	select {
	case <-newBlockCh:
	case <-after:
		panic(cmn.Fmt("Timed out waiting for new block for case '%s' line %d", thisCase.name, i))
	}
}

func runReplayTest(t *testing.T, cs *ConsensusState, walFile string, newBlockCh chan interface{},
	thisCase *testCase, i int) {

	cs.config.SetWalFile(walFile)
	started, err := cs.Start()
	if err != nil {
		t.Fatalf("Cannot start consensus: %v", err)
	}
	if !started {
		t.Error("Consensus did not start")
	}
	// Wait to make a new block.
	// This is just a signal that we haven't halted; its not something contained in the WAL itself.
	// Assuming the consensus state is running, replay of any WAL, including the empty one,
	// should eventually be followed by a new block, or else something is wrong
	waitForBlock(newBlockCh, thisCase, i)
	cs.evsw.Stop()
	cs.Stop()
LOOP:
	for {
		select {
		case <-newBlockCh:
		default:
			break LOOP
		}
	}
	cs.Wait()
}

func toPV(pv PrivValidator) *types.PrivValidator {
	return pv.(*types.PrivValidator)
}

func setupReplayTest(t *testing.T, thisCase *testCase, nLines int, crashAfter bool) (*ConsensusState, chan interface{}, string, string) {
	t.Log("-------------------------------------")
	t.Logf("Starting replay test %v (of %d lines of WAL). Crash after = %v", thisCase.name, nLines, crashAfter)

	lineStep := nLines
	if crashAfter {
		lineStep -= 1
	}

	split := strings.Split(thisCase.log, "\n")
	lastMsg := split[nLines]

	// we write those lines up to (not including) one with the signature
	walFile := writeWAL(strings.Join(split[:nLines], "\n") + "\n")

	cs := fixedConsensusStateDummy()

	// set the last step according to when we crashed vs the wal
	toPV(cs.privValidator).LastHeight = 1 // first block
	toPV(cs.privValidator).LastStep = thisCase.stepMap[lineStep]

	t.Logf("[WARN] setupReplayTest LastStep=%v", toPV(cs.privValidator).LastStep)

	newBlockCh := subscribeToEvent(cs.evsw, "tester", types.EventStringNewBlock(), 1)

	return cs, newBlockCh, lastMsg, walFile
}

func readTimedWALMessage(t *testing.T, walMsg string) TimedWALMessage {
	var err error
	var msg TimedWALMessage
	wire.ReadJSON(&msg, []byte(walMsg), &err)
	if err != nil {
		t.Fatalf("Error reading json data: %v", err)
	}
	return msg
}

//-----------------------------------------------
// Test the log at every iteration, and set the privVal last step
// as if the log was written after signing, before the crash

func TestWALCrashAfterWrite(t *testing.T) {
	for _, thisCase := range testCases {
		split := strings.Split(thisCase.log, "\n")
		for i := 0; i < len(split)-1; i++ {
			cs, newBlockCh, _, walFile := setupReplayTest(t, thisCase, i+1, true)
			runReplayTest(t, cs, walFile, newBlockCh, thisCase, i+1)
		}
	}
}

//-----------------------------------------------
// Test the log as if we crashed after signing but before writing.
// This relies on privValidator.LastSignature being set

func TestWALCrashBeforeWritePropose(t *testing.T) {
	for _, thisCase := range testCases {
		lineNum := thisCase.proposeLine
		// setup replay test where last message is a proposal
		cs, newBlockCh, proposalMsg, walFile := setupReplayTest(t, thisCase, lineNum, false)
		msg := readTimedWALMessage(t, proposalMsg)
		proposal := msg.Msg.(msgInfo).Msg.(*ProposalMessage)
		// Set LastSig
		toPV(cs.privValidator).LastSignBytes = types.SignBytes(cs.state.ChainID, proposal.Proposal)
		toPV(cs.privValidator).LastSignature = proposal.Proposal.Signature
		runReplayTest(t, cs, walFile, newBlockCh, thisCase, lineNum)
	}
}

func TestWALCrashBeforeWritePrevote(t *testing.T) {
	for _, thisCase := range testCases {
		testReplayCrashBeforeWriteVote(t, thisCase, thisCase.prevoteLine, types.EventStringCompleteProposal())
	}
}

func TestWALCrashBeforeWritePrecommit(t *testing.T) {
	for _, thisCase := range testCases {
		testReplayCrashBeforeWriteVote(t, thisCase, thisCase.precommitLine, types.EventStringPolka())
	}
}

func testReplayCrashBeforeWriteVote(t *testing.T, thisCase *testCase, lineNum int, eventString string) {
	// setup replay test where last message is a vote
	cs, newBlockCh, voteMsg, walFile := setupReplayTest(t, thisCase, lineNum, false)
	types.AddListenerForEvent(cs.evsw, "tester", eventString, func(data types.TMEventData) {
		msg := readTimedWALMessage(t, voteMsg)
		vote := msg.Msg.(msgInfo).Msg.(*VoteMessage)
		// Set LastSig
		toPV(cs.privValidator).LastSignBytes = types.SignBytes(cs.state.ChainID, vote.Vote)
		toPV(cs.privValidator).LastSignature = vote.Vote.Signature
	})
	runReplayTest(t, cs, walFile, newBlockCh, thisCase, lineNum)
}

//------------------------------------------------------------------------------------------
// Handshake Tests

var (
	NUM_BLOCKS = 6 // number of blocks in the test_data/many_blocks.cswal
	mempool    = types.MockMempool{}

	testPartSize int
)

//---------------------------------------
// Test handshake/replay

// 0 - all synced up
// 1 - saved block but app and state are behind
// 2 - save block and committed but state is behind
var modes = []uint{0, 1, 2}

// Sync from scratch
func TestHandshakeReplayAll(t *testing.T) {
	for _, m := range modes {
		testHandshakeReplay(t, 0, m)
	}
}

// Sync many, not from scratch
func TestHandshakeReplaySome(t *testing.T) {
	for _, m := range modes {
		testHandshakeReplay(t, 1, m)
	}
}

// Sync from lagging by one
func TestHandshakeReplayOne(t *testing.T) {
	for _, m := range modes {
		testHandshakeReplay(t, NUM_BLOCKS-1, m)
	}
}

// Sync from caught up
func TestHandshakeReplayNone(t *testing.T) {
	for _, m := range modes {
		testHandshakeReplay(t, NUM_BLOCKS, m)
	}
}

// Make some blocks. Start a fresh app and apply nBlocks blocks. Then restart the app and sync it up with the remaining blocks
func testHandshakeReplay(t *testing.T, nBlocks int, mode uint) {
	config := ResetConfig("proxy_test_")

	// copy the many_blocks file
	walBody, err := cmn.ReadFile(path.Join(data_dir, "many_blocks.cswal"))
	if err != nil {
		t.Fatal(err)
	}
	walFile := writeWAL(string(walBody))
	config.Consensus.SetWalFile(walFile)

	privVal := types.LoadPrivValidator(config.PrivValidatorFile())
	testPartSize = config.Consensus.BlockPartSize

	wal, err := NewWAL(walFile, false)
	if err != nil {
		t.Fatal(err)
	}
	wal.SetLogger(log.TestingLogger())
	if _, err := wal.Start(); err != nil {
		t.Fatal(err)
	}
	chain, commits, err := makeBlockchainFromWAL(wal)
	if err != nil {
		t.Fatalf(err.Error())
	}

	state, store := stateAndStore(config, privVal.PubKey)
	store.chain = chain
	store.commits = commits

	// run the chain through state.ApplyBlock to build up the tendermint state
	latestAppHash := buildTMStateFromChain(config, state, chain, mode)

	// make a new client creator
	dummyApp := dummy.NewPersistentDummyApplication(path.Join(config.DBDir(), "2"))
	clientCreator2 := proxy.NewLocalClientCreator(dummyApp)
	if nBlocks > 0 {
		// run nBlocks against a new client to build up the app state.
		// use a throwaway tendermint state
		proxyApp := proxy.NewAppConns(clientCreator2, nil)
		state, _ := stateAndStore(config, privVal.PubKey)
		buildAppStateFromChain(proxyApp, state, chain, nBlocks, mode)
	}

	// now start the app using the handshake - it should sync
	handshaker := NewHandshaker(state, store)
	proxyApp := proxy.NewAppConns(clientCreator2, handshaker)
	if _, err := proxyApp.Start(); err != nil {
		t.Fatalf("Error starting proxy app connections: %v", err)
	}

	// get the latest app hash from the app
	res, err := proxyApp.Query().InfoSync()
	if err != nil {
		t.Fatal(err)
	}

	// the app hash should be synced up
	if !bytes.Equal(latestAppHash, res.LastBlockAppHash) {
		t.Fatalf("Expected app hashes to match after handshake/replay. got %X, expected %X", res.LastBlockAppHash, latestAppHash)
	}

	expectedBlocksToSync := NUM_BLOCKS - nBlocks
	if nBlocks == NUM_BLOCKS && mode > 0 {
		expectedBlocksToSync += 1
	} else if nBlocks > 0 && mode == 1 {
		expectedBlocksToSync += 1
	}

	if handshaker.NBlocks() != expectedBlocksToSync {
		t.Fatalf("Expected handshake to sync %d blocks, got %d", expectedBlocksToSync, handshaker.NBlocks())
	}
}

func applyBlock(st *sm.State, blk *types.Block, proxyApp proxy.AppConns) {
	err := st.ApplyBlock(nil, proxyApp.Consensus(), blk, blk.MakePartSet(testPartSize).Header(), mempool)
	if err != nil {
		panic(err)
	}
}

func buildAppStateFromChain(proxyApp proxy.AppConns,
	state *sm.State, chain []*types.Block, nBlocks int, mode uint) {
	// start a new app without handshake, play nBlocks blocks
	if _, err := proxyApp.Start(); err != nil {
		panic(err)
	}

	validators := types.TM2PB.Validators(state.Validators)
	proxyApp.Consensus().InitChainSync(validators)

	defer proxyApp.Stop()
	switch mode {
	case 0:
		for i := 0; i < nBlocks; i++ {
			block := chain[i]
			applyBlock(state, block, proxyApp)
		}
	case 1, 2:
		for i := 0; i < nBlocks-1; i++ {
			block := chain[i]
			applyBlock(state, block, proxyApp)
		}

		if mode == 2 {
			// update the dummy height and apphash
			// as if we ran commit but not
			applyBlock(state, chain[nBlocks-1], proxyApp)
		}
	}

}

func buildTMStateFromChain(config *cfg.Config, state *sm.State, chain []*types.Block, mode uint) []byte {
	// run the whole chain against this client to build up the tendermint state
	clientCreator := proxy.NewLocalClientCreator(dummy.NewPersistentDummyApplication(path.Join(config.DBDir(), "1")))
	proxyApp := proxy.NewAppConns(clientCreator, nil) // sm.NewHandshaker(config, state, store, ReplayLastBlock))
	if _, err := proxyApp.Start(); err != nil {
		panic(err)
	}
	defer proxyApp.Stop()

	validators := types.TM2PB.Validators(state.Validators)
	proxyApp.Consensus().InitChainSync(validators)

	var latestAppHash []byte

	switch mode {
	case 0:
		// sync right up
		for _, block := range chain {
			applyBlock(state, block, proxyApp)
		}

		latestAppHash = state.AppHash
	case 1, 2:
		// sync up to the penultimate as if we stored the block.
		// whether we commit or not depends on the appHash
		for _, block := range chain[:len(chain)-1] {
			applyBlock(state, block, proxyApp)
		}

		// apply the final block to a state copy so we can
		// get the right next appHash but keep the state back
		stateCopy := state.Copy()
		applyBlock(stateCopy, chain[len(chain)-1], proxyApp)
		latestAppHash = stateCopy.AppHash
	}

	return latestAppHash
}

//--------------------------
// utils for making blocks

func makeBlockchainFromWAL(wal *WAL) ([]*types.Block, []*types.Commit, error) {
	// Search for height marker
	gr, found, err := wal.group.Search("#ENDHEIGHT: ", makeHeightSearchFunc(0))
	if err != nil {
		return nil, nil, err
	}
	if !found {
		return nil, nil, errors.New(cmn.Fmt("WAL does not contain height %d.", 1))
	}
	defer gr.Close()

	// log.Notice("Build a blockchain by reading from the WAL")

	var blockParts *types.PartSet
	var blocks []*types.Block
	var commits []*types.Commit
	for {
		line, err := gr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, nil, err
			}
		}

		piece, err := readPieceFromWAL([]byte(line))
		if err != nil {
			return nil, nil, err
		}
		if piece == nil {
			continue
		}

		switch p := piece.(type) {
		case *types.PartSetHeader:
			// if its not the first one, we have a full block
			if blockParts != nil {
				var n int
				block := wire.ReadBinary(&types.Block{}, blockParts.GetReader(), types.MaxBlockSize, &n, &err).(*types.Block)
				blocks = append(blocks, block)
			}
			blockParts = types.NewPartSetFromHeader(*p)
		case *types.Part:
			_, err := blockParts.AddPart(p, false)
			if err != nil {
				return nil, nil, err
			}
		case *types.Vote:
			if p.Type == types.VoteTypePrecommit {
				commit := &types.Commit{
					BlockID:    p.BlockID,
					Precommits: []*types.Vote{p},
				}
				commits = append(commits, commit)
			}
		}
	}
	// grab the last block too
	var n int
	block := wire.ReadBinary(&types.Block{}, blockParts.GetReader(), types.MaxBlockSize, &n, &err).(*types.Block)
	blocks = append(blocks, block)
	return blocks, commits, nil
}

func readPieceFromWAL(msgBytes []byte) (interface{}, error) {
	// Skip over empty and meta lines
	if len(msgBytes) == 0 || msgBytes[0] == '#' {
		return nil, nil
	}
	var err error
	var msg TimedWALMessage
	wire.ReadJSON(&msg, msgBytes, &err)
	if err != nil {
		fmt.Println("MsgBytes:", msgBytes, string(msgBytes))
		return nil, fmt.Errorf("Error reading json data: %v", err)
	}

	// for logging
	switch m := msg.Msg.(type) {
	case msgInfo:
		switch msg := m.Msg.(type) {
		case *ProposalMessage:
			return &msg.Proposal.BlockPartsHeader, nil
		case *BlockPartMessage:
			return msg.Part, nil
		case *VoteMessage:
			return msg.Vote, nil
		}
	}
	return nil, nil
}

// fresh state and mock store
func stateAndStore(config *cfg.Config, pubKey crypto.PubKey) (*sm.State, *mockBlockStore) {
	stateDB := dbm.NewMemDB()
	state := sm.MakeGenesisStateFromFile(stateDB, config.GenesisFile())
	state.SetLogger(log.TestingLogger().With("module", "state"))

	store := NewMockBlockStore(config)
	return state, store
}

//----------------------------------
// mock block store

type mockBlockStore struct {
	config  *cfg.Config
	chain   []*types.Block
	commits []*types.Commit
}

// TODO: NewBlockStore(db.NewMemDB) ...
func NewMockBlockStore(config *cfg.Config) *mockBlockStore {
	return &mockBlockStore{config, nil, nil}
}

func (bs *mockBlockStore) Height() int                       { return len(bs.chain) }
func (bs *mockBlockStore) LoadBlock(height int) *types.Block { return bs.chain[height-1] }
func (bs *mockBlockStore) LoadBlockMeta(height int) *types.BlockMeta {
	block := bs.chain[height-1]
	return &types.BlockMeta{
		BlockID: types.BlockID{block.Hash(), block.MakePartSet(bs.config.Consensus.BlockPartSize).Header()},
		Header:  block.Header,
	}
}
func (bs *mockBlockStore) LoadBlockPart(height int, index int) *types.Part { return nil }
func (bs *mockBlockStore) SaveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit) {
}
func (bs *mockBlockStore) LoadBlockCommit(height int) *types.Commit {
	return bs.commits[height-1]
}
func (bs *mockBlockStore) LoadSeenCommit(height int) *types.Commit {
	return bs.commits[height-1]
}
