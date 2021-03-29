package types

import (
	// it is ok to use math/rand here: we do not need a cryptographically secure random
	// number generator here and we can run the tests a bit faster
	stdbytes "bytes"
	"context"
	"encoding/hex"
	"math"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	gogotypes "github.com/gogo/protobuf/types"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	coremock "github.com/ipfs/go-ipfs/core/mock"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/lazyledger/lazyledger-core/p2p/ipld/plugin/nodes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lazyledger/lazyledger-core/crypto"
	"github.com/lazyledger/lazyledger-core/crypto/merkle"
	"github.com/lazyledger/lazyledger-core/crypto/tmhash"
	"github.com/lazyledger/lazyledger-core/libs/bits"
	"github.com/lazyledger/lazyledger-core/libs/bytes"
	tmrand "github.com/lazyledger/lazyledger-core/libs/rand"
	tmproto "github.com/lazyledger/lazyledger-core/proto/tendermint/types"
	tmversion "github.com/lazyledger/lazyledger-core/proto/tendermint/version"
	tmtime "github.com/lazyledger/lazyledger-core/types/time"
	"github.com/lazyledger/lazyledger-core/version"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestBlockAddEvidence(t *testing.T) {
	txs := []Tx{Tx("foo"), Tx("bar")}
	lastID := makeBlockIDRandom()
	h := int64(3)

	voteSet, _, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	require.NoError(t, err)

	ev := NewMockDuplicateVoteEvidenceWithValidator(h, time.Now(), vals[0], "block-test-chain")
	evList := []Evidence{ev}

	block := MakeBlock(h, txs, evList, nil, Messages{}, commit)
	require.NotNil(t, block)
	require.Equal(t, 1, len(block.Evidence.Evidence))
	require.NotNil(t, block.EvidenceHash)
}

func TestBlockValidateBasic(t *testing.T) {
	require.Error(t, (*Block)(nil).ValidateBasic())

	txs := []Tx{Tx("foo"), Tx("bar")}
	lastID := makeBlockIDRandom()
	h := int64(3)

	voteSet, valSet, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	require.NoError(t, err)

	ev := NewMockDuplicateVoteEvidenceWithValidator(h, time.Now(), vals[0], "block-test-chain")
	evList := []Evidence{ev}

	testCases := []struct {
		testName      string
		malleateBlock func(*Block)
		expErr        bool
	}{
		{"Make Block", func(blk *Block) {}, false},
		{"Make Block w/ proposer Addr", func(blk *Block) { blk.ProposerAddress = valSet.GetProposer().Address }, false},
		{"Negative Height", func(blk *Block) { blk.Height = -1 }, true},
		{"Remove 1/2 the commits", func(blk *Block) {
			blk.LastCommit.Signatures = commit.Signatures[:commit.Size()/2]
			blk.LastCommit.hash = nil // clear hash or change wont be noticed
		}, true},
		{"Remove LastCommitHash", func(blk *Block) { blk.LastCommitHash = []byte("something else") }, true},
		{"Tampered Data", func(blk *Block) {
			blk.Data.Txs[0] = Tx("something else")
			blk.DataHash = nil // clear hash or change wont be noticed
		}, true},
		{"Tampered DataHash", func(blk *Block) {
			blk.DataHash = tmrand.Bytes(len(blk.DataHash))
		}, true},
		{"Tampered EvidenceHash", func(blk *Block) {
			blk.EvidenceHash = tmrand.Bytes(len(blk.EvidenceHash))
		}, true},
		{"Incorrect block protocol version", func(blk *Block) {
			blk.Version.Block = 1
		}, true},
		{"Missing LastCommit", func(blk *Block) {
			blk.LastCommit = nil
		}, true},
		{"Invalid LastCommit", func(blk *Block) {
			blk.LastCommit = NewCommit(-1, 0, *voteSet.maj23, nil)
		}, true},
		{"Invalid Evidence", func(blk *Block) {
			emptyEv := &DuplicateVoteEvidence{}
			blk.Evidence = EvidenceData{Evidence: []Evidence{emptyEv}}
		}, true},
	}
	for i, tc := range testCases {
		tc := tc
		i := i
		t.Run(tc.testName, func(t *testing.T) {
			block := MakeBlock(h, txs, evList, nil, Messages{}, commit)
			block.ProposerAddress = valSet.GetProposer().Address
			tc.malleateBlock(block)
			err = block.ValidateBasic()
			t.Log(err)
			assert.Equal(t, tc.expErr, err != nil, "#%d: %v", i, err)
		})
	}
}

func TestBlockHash(t *testing.T) {
	assert.Nil(t, (*Block)(nil).Hash())
	assert.Nil(t, MakeBlock(int64(3), []Tx{Tx("Hello World")}, nil, nil, Messages{}, nil).Hash())
}

func TestBlockMakePartSet(t *testing.T) {
	assert.Nil(t, (*Block)(nil).MakePartSet(2))

	partSet := MakeBlock(int64(3), []Tx{Tx("Hello World")}, nil, nil, Messages{}, nil).MakePartSet(1024)
	assert.NotNil(t, partSet)
	assert.EqualValues(t, 1, partSet.Total())
}

func TestBlockMakePartSetWithEvidence(t *testing.T) {
	assert.Nil(t, (*Block)(nil).MakePartSet(2))

	lastID := makeBlockIDRandom()
	h := int64(3)

	voteSet, _, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	require.NoError(t, err)

	ev := NewMockDuplicateVoteEvidenceWithValidator(h, time.Now(), vals[0], "block-test-chain")
	evList := []Evidence{ev}

	partSet := MakeBlock(h, []Tx{Tx("Hello World")}, evList, nil, Messages{}, commit).MakePartSet(512)
	assert.NotNil(t, partSet)
	assert.EqualValues(t, 5, partSet.Total())
}

func TestBlockHashesTo(t *testing.T) {
	assert.False(t, (*Block)(nil).HashesTo(nil))

	lastID := makeBlockIDRandom()
	h := int64(3)
	voteSet, valSet, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	require.NoError(t, err)

	ev := NewMockDuplicateVoteEvidenceWithValidator(h, time.Now(), vals[0], "block-test-chain")
	evList := []Evidence{ev}

	block := MakeBlock(h, []Tx{Tx("Hello World")}, evList, nil, Messages{}, commit)
	block.ValidatorsHash = valSet.Hash()
	assert.False(t, block.HashesTo([]byte{}))
	assert.False(t, block.HashesTo([]byte("something else")))
	assert.True(t, block.HashesTo(block.Hash()))
}

func TestBlockSize(t *testing.T) {
	size := MakeBlock(int64(3), []Tx{Tx("Hello World")}, nil, nil, Messages{}, nil).Size()
	if size <= 0 {
		t.Fatal("Size of the block is zero or negative")
	}
}

func TestBlockString(t *testing.T) {
	assert.Equal(t, "nil-Block", (*Block)(nil).String())
	assert.Equal(t, "nil-Block", (*Block)(nil).StringIndented(""))
	assert.Equal(t, "nil-Block", (*Block)(nil).StringShort())

	block := MakeBlock(int64(3), []Tx{Tx("Hello World")}, nil, nil, Messages{}, nil)
	assert.NotEqual(t, "nil-Block", block.String())
	assert.NotEqual(t, "nil-Block", block.StringIndented(""))
	assert.NotEqual(t, "nil-Block", block.StringShort())
}

func makeBlockIDRandom() BlockID {
	var (
		blockHash   = make([]byte, tmhash.Size)
		partSetHash = make([]byte, tmhash.Size)
	)
	rand.Read(blockHash)
	rand.Read(partSetHash)
	return BlockID{blockHash, PartSetHeader{123, partSetHash}}
}

func makeBlockID(hash []byte, partSetSize uint32, partSetHash []byte) BlockID {
	var (
		h   = make([]byte, tmhash.Size)
		psH = make([]byte, tmhash.Size)
	)
	copy(h, hash)
	copy(psH, partSetHash)
	return BlockID{
		Hash: h,
		PartSetHeader: PartSetHeader{
			Total: partSetSize,
			Hash:  psH,
		},
	}
}

var nilBytes []byte

// This follows RFC-6962, i.e. `echo -n '' | sha256sum`
var emptyBytes = []byte{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8,
	0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b,
	0x78, 0x52, 0xb8, 0x55}

func TestNilHeaderHashDoesntCrash(t *testing.T) {
	assert.Equal(t, nilBytes, []byte((*Header)(nil).Hash()))
	assert.Equal(t, nilBytes, []byte((new(Header)).Hash()))
}

func TestNilDataAvailabilityHeaderHashDoesntCrash(t *testing.T) {
	assert.Equal(t, emptyBytes, (*DataAvailabilityHeader)(nil).Hash())
	assert.Equal(t, emptyBytes, new(DataAvailabilityHeader).Hash())
}

func TestCommit(t *testing.T) {
	lastID := makeBlockIDRandom()
	h := int64(3)
	voteSet, _, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	require.NoError(t, err)

	assert.Equal(t, h-1, commit.Height)
	assert.EqualValues(t, 1, commit.Round)
	assert.Equal(t, tmproto.PrecommitType, tmproto.SignedMsgType(commit.Type()))
	if commit.Size() <= 0 {
		t.Fatalf("commit %v has a zero or negative size: %d", commit, commit.Size())
	}

	require.NotNil(t, commit.BitArray())
	assert.Equal(t, bits.NewBitArray(10).Size(), commit.BitArray().Size())

	assert.Equal(t, voteSet.GetByIndex(0), commit.GetByIndex(0))
	assert.True(t, commit.IsCommit())
}

func TestCommitValidateBasic(t *testing.T) {
	testCases := []struct {
		testName       string
		malleateCommit func(*Commit)
		expectErr      bool
	}{
		{"Random Commit", func(com *Commit) {}, false},
		{"Incorrect signature", func(com *Commit) { com.Signatures[0].Signature = []byte{0} }, false},
		{"Incorrect height", func(com *Commit) { com.Height = int64(-100) }, true},
		{"Incorrect round", func(com *Commit) { com.Round = -100 }, true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			com := randCommit(time.Now())
			tc.malleateCommit(com)
			assert.Equal(t, tc.expectErr, com.ValidateBasic() != nil, "Validate Basic had an unexpected result")
		})
	}
}

func TestMaxCommitBytes(t *testing.T) {
	// time is varint encoded so need to pick the max.
	// year int, month Month, day, hour, min, sec, nsec int, loc *Location
	timestamp := time.Date(math.MaxInt64, 0, 0, 0, 0, 0, math.MaxInt64, time.UTC)

	cs := CommitSig{
		BlockIDFlag:      BlockIDFlagNil,
		ValidatorAddress: crypto.AddressHash([]byte("validator_address")),
		Timestamp:        timestamp,
		Signature:        crypto.CRandBytes(MaxSignatureSize),
	}

	pbSig := cs.ToProto()
	// test that a single commit sig doesn't exceed max commit sig bytes
	assert.EqualValues(t, MaxCommitSigBytes, pbSig.Size())

	// check size with a single commit
	commit := &Commit{
		Height: math.MaxInt64,
		Round:  math.MaxInt32,
		BlockID: BlockID{
			Hash: tmhash.Sum([]byte("blockID_hash")),
			PartSetHeader: PartSetHeader{
				Total: math.MaxInt32,
				Hash:  tmhash.Sum([]byte("blockID_part_set_header_hash")),
			},
		},
		Signatures: []CommitSig{cs},
	}

	pb := commit.ToProto()

	assert.EqualValues(t, MaxCommitBytes(1), int64(pb.Size()))

	// check the upper bound of the commit size
	for i := 1; i < MaxVotesCount; i++ {
		commit.Signatures = append(commit.Signatures, cs)
	}

	pb = commit.ToProto()

	assert.EqualValues(t, MaxCommitBytes(MaxVotesCount), int64(pb.Size()))

}

func TestHeaderHash(t *testing.T) {
	testCases := []struct {
		desc       string
		header     *Header
		expectHash bytes.HexBytes
	}{
		{"Generates expected hash", &Header{
			Version:            tmversion.Consensus{Block: 1, App: 2},
			ChainID:            "chainId",
			Height:             3,
			Time:               time.Date(2019, 10, 13, 16, 14, 44, 0, time.UTC),
			LastBlockID:        makeBlockID(make([]byte, tmhash.Size), 6, make([]byte, tmhash.Size)),
			LastCommitHash:     tmhash.Sum([]byte("last_commit_hash")),
			DataHash:           tmhash.Sum([]byte("data_hash")),
			ValidatorsHash:     tmhash.Sum([]byte("validators_hash")),
			NextValidatorsHash: tmhash.Sum([]byte("next_validators_hash")),
			ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
			AppHash:            tmhash.Sum([]byte("app_hash")),
			LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
			EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
			ProposerAddress:    crypto.AddressHash([]byte("proposer_address")),
		}, hexBytesFromString("F740121F553B5418C3EFBD343C2DBFE9E007BB67B0D020A0741374BAB65242A4")},
		{"nil header yields nil", nil, nil},
		{"nil ValidatorsHash yields nil", &Header{
			Version:            tmversion.Consensus{Block: 1, App: 2},
			ChainID:            "chainId",
			Height:             3,
			Time:               time.Date(2019, 10, 13, 16, 14, 44, 0, time.UTC),
			LastBlockID:        makeBlockID(make([]byte, tmhash.Size), 6, make([]byte, tmhash.Size)),
			LastCommitHash:     tmhash.Sum([]byte("last_commit_hash")),
			DataHash:           tmhash.Sum([]byte("data_hash")),
			ValidatorsHash:     nil,
			NextValidatorsHash: tmhash.Sum([]byte("next_validators_hash")),
			ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
			AppHash:            tmhash.Sum([]byte("app_hash")),
			LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
			EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
			ProposerAddress:    crypto.AddressHash([]byte("proposer_address")),
		}, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expectHash, tc.header.Hash())

			// We also make sure that all fields are hashed in struct order, and that all
			// fields in the test struct are non-zero.
			if tc.header != nil && tc.expectHash != nil {
				byteSlices := [][]byte{}

				s := reflect.ValueOf(*tc.header)
				for i := 0; i < s.NumField(); i++ {
					f := s.Field(i)

					assert.False(t, f.IsZero(), "Found zero-valued field %v",
						s.Type().Field(i).Name)

					switch f := f.Interface().(type) {
					case int64, bytes.HexBytes, string:
						byteSlices = append(byteSlices, cdcEncode(f))
					case time.Time:
						bz, err := gogotypes.StdTimeMarshal(f)
						require.NoError(t, err)
						byteSlices = append(byteSlices, bz)
					case tmversion.Consensus:
						bz, err := f.Marshal()
						require.NoError(t, err)
						byteSlices = append(byteSlices, bz)
					case BlockID:
						pbbi := f.ToProto()
						bz, err := pbbi.Marshal()
						require.NoError(t, err)
						byteSlices = append(byteSlices, bz)
					default:
						t.Errorf("unknown type %T", f)
					}
				}
				assert.Equal(t,
					bytes.HexBytes(merkle.HashFromByteSlices(byteSlices)), tc.header.Hash())
			}
		})
	}
}

func TestMaxHeaderBytes(t *testing.T) {
	// Construct a UTF-8 string of MaxChainIDLen length using the supplementary
	// characters.
	// Each supplementary character takes 4 bytes.
	// http://www.i18nguy.com/unicode/supplementary-test.html
	maxChainID := ""
	for i := 0; i < MaxChainIDLen; i++ {
		maxChainID += "𠜎"
	}

	// time is varint encoded so need to pick the max.
	// year int, month Month, day, hour, min, sec, nsec int, loc *Location
	timestamp := time.Date(math.MaxInt64, 0, 0, 0, 0, 0, math.MaxInt64, time.UTC)

	h := Header{
		Version:            tmversion.Consensus{Block: math.MaxInt64, App: math.MaxInt64},
		ChainID:            maxChainID,
		Height:             math.MaxInt64,
		Time:               timestamp,
		LastBlockID:        makeBlockID(make([]byte, tmhash.Size), math.MaxInt32, make([]byte, tmhash.Size)),
		LastCommitHash:     tmhash.Sum([]byte("last_commit_hash")),
		DataHash:           tmhash.Sum([]byte("data_hash")),
		ValidatorsHash:     tmhash.Sum([]byte("validators_hash")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators_hash")),
		ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
		AppHash:            tmhash.Sum([]byte("app_hash")),
		LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
		EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
		ProposerAddress:    crypto.AddressHash([]byte("proposer_address")),
	}

	bz, err := h.ToProto().Marshal()
	require.NoError(t, err)

	assert.EqualValues(t, MaxHeaderBytes, int64(len(bz)))
}

func randCommit(now time.Time) *Commit {
	lastID := makeBlockIDRandom()
	h := int64(3)
	voteSet, _, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, now)
	if err != nil {
		panic(err)
	}
	return commit
}

func hexBytesFromString(s string) bytes.HexBytes {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes.HexBytes(b)
}

func TestBlockMaxDataBytes(t *testing.T) {
	testCases := []struct {
		maxBytes      int64
		valsCount     int
		evidenceBytes int64
		panics        bool
		result        int64
	}{
		0: {-10, 1, 0, true, 0},
		1: {10, 1, 0, true, 0},
		2: {841, 1, 0, true, 0},
		3: {842, 1, 0, false, 0},
		4: {843, 1, 0, false, 1},
		5: {954, 2, 0, false, 1},
		6: {1053, 2, 100, false, 0},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.panics {
			assert.Panics(t, func() {
				MaxDataBytes(tc.maxBytes, tc.evidenceBytes, tc.valsCount)
			}, "#%v", i)
		} else {
			assert.Equal(t,
				tc.result,
				MaxDataBytes(tc.maxBytes, tc.evidenceBytes, tc.valsCount),
				"#%v", i)
		}
	}
}

func TestBlockMaxDataBytesNoEvidence(t *testing.T) {
	testCases := []struct {
		maxBytes  int64
		valsCount int
		panics    bool
		result    int64
	}{
		0: {-10, 1, true, 0},
		1: {10, 1, true, 0},
		2: {841, 1, true, 0},
		3: {842, 1, false, 0},
		4: {843, 1, false, 1},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.panics {
			assert.Panics(t, func() {
				MaxDataBytesNoEvidence(tc.maxBytes, tc.valsCount)
			}, "#%v", i)
		} else {
			assert.Equal(t,
				tc.result,
				MaxDataBytesNoEvidence(tc.maxBytes, tc.valsCount),
				"#%v", i)
		}
	}
}

func TestCommitToVoteSet(t *testing.T) {
	lastID := makeBlockIDRandom()
	h := int64(3)

	voteSet, valSet, vals := randVoteSet(h-1, 1, tmproto.PrecommitType, 10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals, time.Now())
	assert.NoError(t, err)

	chainID := voteSet.ChainID()
	voteSet2 := CommitToVoteSet(chainID, commit, valSet)

	for i := int32(0); int(i) < len(vals); i++ {
		vote1 := voteSet.GetByIndex(i)
		vote2 := voteSet2.GetByIndex(i)
		vote3 := commit.GetVote(i)

		vote1bz, err := vote1.ToProto().Marshal()
		require.NoError(t, err)
		vote2bz, err := vote2.ToProto().Marshal()
		require.NoError(t, err)
		vote3bz, err := vote3.ToProto().Marshal()
		require.NoError(t, err)
		assert.Equal(t, vote1bz, vote2bz)
		assert.Equal(t, vote1bz, vote3bz)
	}
}

func TestCommitToVoteSetWithVotesForNilBlock(t *testing.T) {
	blockID := makeBlockID([]byte("blockhash"), 1000, []byte("partshash"))

	const (
		height = int64(3)
		round  = 0
	)

	type commitVoteTest struct {
		blockIDs      []BlockID
		numVotes      []int // must sum to numValidators
		numValidators int
		valid         bool
	}

	testCases := []commitVoteTest{
		{[]BlockID{blockID, {}}, []int{67, 33}, 100, true},
	}

	for _, tc := range testCases {
		voteSet, valSet, vals := randVoteSet(height-1, round, tmproto.PrecommitType, tc.numValidators, 1)

		vi := int32(0)
		for n := range tc.blockIDs {
			for i := 0; i < tc.numVotes[n]; i++ {
				pubKey, err := vals[vi].GetPubKey()
				require.NoError(t, err)
				vote := &Vote{
					ValidatorAddress: pubKey.Address(),
					ValidatorIndex:   vi,
					Height:           height - 1,
					Round:            round,
					Type:             tmproto.PrecommitType,
					BlockID:          tc.blockIDs[n],
					Timestamp:        tmtime.Now(),
				}

				added, err := signAddVote(vals[vi], vote, voteSet)
				assert.NoError(t, err)
				assert.True(t, added)

				vi++
			}
		}

		if tc.valid {
			commit := voteSet.MakeCommit() // panics without > 2/3 valid votes
			assert.NotNil(t, commit)
			err := valSet.VerifyCommit(voteSet.ChainID(), blockID, height-1, commit)
			assert.Nil(t, err)
		} else {
			assert.Panics(t, func() { voteSet.MakeCommit() })
		}
	}
}

func TestBlockIDValidateBasic(t *testing.T) {
	validBlockID := BlockID{
		Hash: bytes.HexBytes{},
		PartSetHeader: PartSetHeader{
			Total: 1,
			Hash:  bytes.HexBytes{},
		},
	}

	invalidBlockID := BlockID{
		Hash: []byte{0},
		PartSetHeader: PartSetHeader{
			Total: 1,
			Hash:  []byte{0},
		},
	}

	testCases := []struct {
		testName             string
		blockIDHash          bytes.HexBytes
		blockIDPartSetHeader PartSetHeader
		expectErr            bool
	}{
		{"Valid BlockID", validBlockID.Hash, validBlockID.PartSetHeader, false},
		{"Invalid BlockID", invalidBlockID.Hash, validBlockID.PartSetHeader, true},
		{"Invalid BlockID", validBlockID.Hash, invalidBlockID.PartSetHeader, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			blockID := BlockID{
				Hash:          tc.blockIDHash,
				PartSetHeader: tc.blockIDPartSetHeader,
			}
			assert.Equal(t, tc.expectErr, blockID.ValidateBasic() != nil, "Validate Basic had an unexpected result")
		})
	}
}

func TestBlockProtoBuf(t *testing.T) {
	h := tmrand.Int63()
	c1 := randCommit(time.Now())
	b1 := MakeBlock(h, []Tx{Tx([]byte{1})}, []Evidence{}, nil, Messages{}, &Commit{Signatures: []CommitSig{}})
	b1.ProposerAddress = tmrand.Bytes(crypto.AddressSize)

	b2 := MakeBlock(h, []Tx{Tx([]byte{1})}, []Evidence{}, nil, Messages{}, c1)
	b2.ProposerAddress = tmrand.Bytes(crypto.AddressSize)
	evidenceTime := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	evi := NewMockDuplicateVoteEvidence(h, evidenceTime, "block-test-chain")
	b2.Evidence = EvidenceData{Evidence: EvidenceList{evi}}
	b2.EvidenceHash = b2.Evidence.Hash()

	b3 := MakeBlock(h, []Tx{}, []Evidence{}, nil, Messages{}, c1)
	b3.ProposerAddress = tmrand.Bytes(crypto.AddressSize)
	testCases := []struct {
		msg      string
		b1       *Block
		expPass  bool
		expPass2 bool
	}{
		{"nil block", nil, false, false},
		{"b1", b1, true, true},
		{"b2", b2, true, true},
		{"b3", b3, true, true},
	}
	for _, tc := range testCases {
		pb, err := tc.b1.ToProto()
		if tc.expPass {
			require.NoError(t, err, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}

		block, err := BlockFromProto(pb)
		if tc.expPass2 {
			require.NoError(t, err, tc.msg)
			require.EqualValues(t, tc.b1.Header, block.Header, tc.msg)
			// require.EqualValues(t, tc.b1.Data, block.Data, tc.msg)
			require.EqualValues(t, tc.b1.Evidence.Evidence, block.Evidence.Evidence, tc.msg)
			require.EqualValues(t, *tc.b1.LastCommit, *block.LastCommit, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func TestDataProtoBuf(t *testing.T) {
	data := &Data{Txs: Txs{Tx([]byte{1}), Tx([]byte{2}), Tx([]byte{3})}}
	data2 := &Data{Txs: Txs{}}
	testCases := []struct {
		msg     string
		data1   *Data
		expPass bool
	}{
		{"success", data, true},
		{"success data2", data2, true},
	}
	for _, tc := range testCases {
		protoData := tc.data1.ToProto()
		d, err := DataFromProto(&protoData)
		if tc.expPass {
			require.NoError(t, err, tc.msg)
			require.EqualValues(t, tc.data1, &d, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

// TestEvidenceDataProtoBuf ensures parity in converting to and from proto.
func TestEvidenceDataProtoBuf(t *testing.T) {
	const chainID = "mychain"
	ev := NewMockDuplicateVoteEvidence(math.MaxInt64, time.Now(), chainID)
	data := &EvidenceData{Evidence: EvidenceList{ev}}
	_ = data.ByteSize()
	testCases := []struct {
		msg      string
		data1    *EvidenceData
		expPass1 bool
		expPass2 bool
	}{
		{"success", data, true, true},
		{"empty evidenceData", &EvidenceData{Evidence: EvidenceList{}}, true, true},
		{"fail nil Data", nil, false, false},
	}

	for _, tc := range testCases {
		protoData, err := tc.data1.ToProto()
		if tc.expPass1 {
			require.NoError(t, err, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}

		eviD := new(EvidenceData)
		err = eviD.FromProto(protoData)
		if tc.expPass2 {
			require.NoError(t, err, tc.msg)
			require.Equal(t, tc.data1, eviD, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func makeRandHeader() Header {
	chainID := "test"
	t := time.Now()
	height := tmrand.Int63()
	randBytes := tmrand.Bytes(tmhash.Size)
	randAddress := tmrand.Bytes(crypto.AddressSize)
	h := Header{
		Version:            tmversion.Consensus{Block: version.BlockProtocol, App: 1},
		ChainID:            chainID,
		Height:             height,
		Time:               t,
		LastBlockID:        BlockID{},
		LastCommitHash:     randBytes,
		DataHash:           randBytes,
		ValidatorsHash:     randBytes,
		NextValidatorsHash: randBytes,
		ConsensusHash:      randBytes,
		AppHash:            randBytes,

		LastResultsHash: randBytes,

		EvidenceHash:    randBytes,
		ProposerAddress: randAddress,
	}

	return h
}

func TestHeaderProto(t *testing.T) {
	h1 := makeRandHeader()
	tc := []struct {
		msg     string
		h1      *Header
		expPass bool
	}{
		{"success", &h1, true},
		{"failure empty Header", &Header{}, false},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			pb := tt.h1.ToProto()
			h, err := HeaderFromProto(pb)
			if tt.expPass {
				require.NoError(t, err, tt.msg)
				require.Equal(t, tt.h1, &h, tt.msg)
			} else {
				require.Error(t, err, tt.msg)
			}

		})
	}
}

func TestBlockIDProtoBuf(t *testing.T) {
	blockID := makeBlockID([]byte("hash"), 2, []byte("part_set_hash"))
	testCases := []struct {
		msg     string
		bid1    *BlockID
		expPass bool
	}{
		{"success", &blockID, true},
		{"success empty", &BlockID{}, true},
		{"failure BlockID nil", nil, false},
	}
	for _, tc := range testCases {
		protoBlockID := tc.bid1.ToProto()

		bi, err := BlockIDFromProto(&protoBlockID)
		if tc.expPass {
			require.NoError(t, err)
			require.Equal(t, tc.bid1, bi, tc.msg)
		} else {
			require.NotEqual(t, tc.bid1, bi, tc.msg)
		}
	}
}

func TestSignedHeaderProtoBuf(t *testing.T) {
	commit := randCommit(time.Now())
	h := makeRandHeader()

	sh := SignedHeader{Header: &h, Commit: commit}

	testCases := []struct {
		msg     string
		sh1     *SignedHeader
		expPass bool
	}{
		{"empty SignedHeader 2", &SignedHeader{}, true},
		{"success", &sh, true},
		{"failure nil", nil, false},
	}
	for _, tc := range testCases {
		protoSignedHeader := tc.sh1.ToProto()

		sh, err := SignedHeaderFromProto(protoSignedHeader)

		if tc.expPass {
			require.NoError(t, err, tc.msg)
			require.Equal(t, tc.sh1, sh, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func TestBlockIDEquals(t *testing.T) {
	var (
		blockID          = makeBlockID([]byte("hash"), 2, []byte("part_set_hash"))
		blockIDDuplicate = makeBlockID([]byte("hash"), 2, []byte("part_set_hash"))
		blockIDDifferent = makeBlockID([]byte("different_hash"), 2, []byte("part_set_hash"))
		blockIDEmpty     = BlockID{}
	)

	assert.True(t, blockID.Equals(blockIDDuplicate))
	assert.False(t, blockID.Equals(blockIDDifferent))
	assert.False(t, blockID.Equals(blockIDEmpty))
	assert.True(t, blockIDEmpty.Equals(blockIDEmpty))
	assert.False(t, blockIDEmpty.Equals(blockIDDifferent))
}

func TestCommitSig_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		cs        CommitSig
		expectErr bool
		errString string
	}{
		{
			"invalid ID flag",
			CommitSig{BlockIDFlag: BlockIDFlag(0xFF)},
			true, "unknown BlockIDFlag",
		},
		{
			"BlockIDFlagAbsent validator address present",
			CommitSig{BlockIDFlag: BlockIDFlagAbsent, ValidatorAddress: crypto.Address("testaddr")},
			true, "validator address is present",
		},
		{
			"BlockIDFlagAbsent timestamp present",
			CommitSig{BlockIDFlag: BlockIDFlagAbsent, Timestamp: time.Now().UTC()},
			true, "time is present",
		},
		{
			"BlockIDFlagAbsent signatures present",
			CommitSig{BlockIDFlag: BlockIDFlagAbsent, Signature: []byte{0xAA}},
			true, "signature is present",
		},
		{
			"BlockIDFlagAbsent valid BlockIDFlagAbsent",
			CommitSig{BlockIDFlag: BlockIDFlagAbsent},
			false, "",
		},
		{
			"non-BlockIDFlagAbsent invalid validator address",
			CommitSig{BlockIDFlag: BlockIDFlagCommit, ValidatorAddress: make([]byte, 1)},
			true, "expected ValidatorAddress size",
		},
		{
			"non-BlockIDFlagAbsent invalid signature (zero)",
			CommitSig{
				BlockIDFlag:      BlockIDFlagCommit,
				ValidatorAddress: make([]byte, crypto.AddressSize),
				Signature:        make([]byte, 0),
			},
			true, "signature is missing",
		},
		{
			"non-BlockIDFlagAbsent invalid signature (too large)",
			CommitSig{
				BlockIDFlag:      BlockIDFlagCommit,
				ValidatorAddress: make([]byte, crypto.AddressSize),
				Signature:        make([]byte, MaxSignatureSize+1),
			},
			true, "signature is too big",
		},
		{
			"non-BlockIDFlagAbsent valid",
			CommitSig{
				BlockIDFlag:      BlockIDFlagCommit,
				ValidatorAddress: make([]byte, crypto.AddressSize),
				Signature:        make([]byte, MaxSignatureSize),
			},
			false, "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.cs.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHeader_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		header    Header
		expectErr bool
		errString string
	}{
		{
			"invalid version block",
			Header{Version: tmversion.Consensus{Block: version.BlockProtocol + 1}},
			true, "block protocol is incorrect",
		},
		{
			"invalid chain ID length",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen+1)),
			},
			true, "chainID is too long",
		},
		{
			"invalid height (negative)",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  -1,
			},
			true, "negative Height",
		},
		{
			"invalid height (zero)",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  0,
			},
			true, "zero Height",
		},
		{
			"invalid block ID hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size+1),
				},
			},
			true, "wrong Hash",
		},
		{
			"invalid block ID parts header hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size+1),
					},
				},
			},
			true, "wrong PartSetHeader",
		},
		{
			"invalid last commit hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash: make([]byte, tmhash.Size+1),
			},
			true, "wrong LastCommitHash",
		},
		{
			"invalid data hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash: make([]byte, tmhash.Size),
				DataHash:       make([]byte, tmhash.Size+1),
			},
			true, "wrong DataHash",
		},
		{
			"invalid evidence hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash: make([]byte, tmhash.Size),
				DataHash:       make([]byte, tmhash.Size),
				EvidenceHash:   make([]byte, tmhash.Size+1),
			},
			true, "wrong EvidenceHash",
		},
		{
			"invalid proposer address",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:  make([]byte, tmhash.Size),
				DataHash:        make([]byte, tmhash.Size),
				EvidenceHash:    make([]byte, tmhash.Size),
				ProposerAddress: make([]byte, crypto.AddressSize+1),
			},
			true, "invalid ProposerAddress length",
		},
		{
			"invalid validator hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:  make([]byte, tmhash.Size),
				DataHash:        make([]byte, tmhash.Size),
				EvidenceHash:    make([]byte, tmhash.Size),
				ProposerAddress: make([]byte, crypto.AddressSize),
				ValidatorsHash:  make([]byte, tmhash.Size+1),
			},
			true, "wrong ValidatorsHash",
		},
		{
			"invalid next validator hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:     make([]byte, tmhash.Size),
				DataHash:           make([]byte, tmhash.Size),
				EvidenceHash:       make([]byte, tmhash.Size),
				ProposerAddress:    make([]byte, crypto.AddressSize),
				ValidatorsHash:     make([]byte, tmhash.Size),
				NextValidatorsHash: make([]byte, tmhash.Size+1),
			},
			true, "wrong NextValidatorsHash",
		},
		{
			"invalid consensus hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:     make([]byte, tmhash.Size),
				DataHash:           make([]byte, tmhash.Size),
				EvidenceHash:       make([]byte, tmhash.Size),
				ProposerAddress:    make([]byte, crypto.AddressSize),
				ValidatorsHash:     make([]byte, tmhash.Size),
				NextValidatorsHash: make([]byte, tmhash.Size),
				ConsensusHash:      make([]byte, tmhash.Size+1),
			},
			true, "wrong ConsensusHash",
		},
		{
			"invalid last results hash",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:     make([]byte, tmhash.Size),
				DataHash:           make([]byte, tmhash.Size),
				EvidenceHash:       make([]byte, tmhash.Size),
				ProposerAddress:    make([]byte, crypto.AddressSize),
				ValidatorsHash:     make([]byte, tmhash.Size),
				NextValidatorsHash: make([]byte, tmhash.Size),
				ConsensusHash:      make([]byte, tmhash.Size),
				LastResultsHash:    make([]byte, tmhash.Size+1),
			},
			true, "wrong LastResultsHash",
		},
		{
			"valid header",
			Header{
				Version: tmversion.Consensus{Block: version.BlockProtocol},
				ChainID: string(make([]byte, MaxChainIDLen)),
				Height:  1,
				LastBlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				LastCommitHash:     make([]byte, tmhash.Size),
				DataHash:           make([]byte, tmhash.Size),
				EvidenceHash:       make([]byte, tmhash.Size),
				ProposerAddress:    make([]byte, crypto.AddressSize),
				ValidatorsHash:     make([]byte, tmhash.Size),
				NextValidatorsHash: make([]byte, tmhash.Size),
				ConsensusHash:      make([]byte, tmhash.Size),
				LastResultsHash:    make([]byte, tmhash.Size),
			},
			false, "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.header.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommit_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		commit    *Commit
		expectErr bool
		errString string
	}{
		{
			"invalid height",
			&Commit{Height: -1},
			true, "negative Height",
		},
		{
			"invalid round",
			&Commit{Height: 1, Round: -1},
			true, "negative Round",
		},
		{
			"invalid block ID",
			&Commit{
				Height:     1,
				Round:      1,
				BlockID:    BlockID{},
				HeaderHash: make([]byte, tmhash.Size),
			},
			true, "commit cannot be for nil block",
		},
		{
			"no signatures",
			&Commit{
				Height: 1,
				Round:  1,
				BlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				HeaderHash: make([]byte, tmhash.Size),
			},
			true, "no signatures in commit",
		},
		{
			"invalid signature",
			&Commit{
				Height: 1,
				Round:  1,
				BlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				Signatures: []CommitSig{
					{
						BlockIDFlag:      BlockIDFlagCommit,
						ValidatorAddress: make([]byte, crypto.AddressSize),
						Signature:        make([]byte, MaxSignatureSize+1),
					},
				},
				HeaderHash: make([]byte, tmhash.Size),
			},
			true, "wrong CommitSig",
		},
		{
			"valid commit",
			&Commit{
				Height: 1,
				Round:  1,
				BlockID: BlockID{
					Hash: make([]byte, tmhash.Size),
					PartSetHeader: PartSetHeader{
						Hash: make([]byte, tmhash.Size),
					},
				},
				Signatures: []CommitSig{
					{
						BlockIDFlag:      BlockIDFlagCommit,
						ValidatorAddress: make([]byte, crypto.AddressSize),
						Signature:        make([]byte, MaxSignatureSize),
					},
				},
				HeaderHash: make([]byte, tmhash.Size),
			},
			false, "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.commit.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPutBlock(t *testing.T) {
	ipfsNode, err := coremock.NewMockNode()
	if err != nil {
		t.Error(err)
	}

	ipfsAPI, err := coreapi.NewCoreAPI(ipfsNode)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		name      string
		blockData Data
		expectErr bool
		errString string
	}{
		{"no leaves", generateRandomMsgOnlyData(0), false, ""},
		{"single leaf", generateRandomMsgOnlyData(1), false, ""},
		{"16 leaves", generateRandomMsgOnlyData(16), false, ""},
		{"max square size", generateRandomMsgOnlyData(MaxSquareSize), false, ""},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		tc := tc

		block := &Block{Data: tc.blockData}

		t.Run(tc.name, func(t *testing.T) {
			err = block.PutBlock(ctx, ipfsAPI.Dag().Pinning())
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errString)
				return
			}

			require.NoError(t, err)

			timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			block.fillDataAvailabilityHeader()
			for _, rowRoot := range block.DataAvailabilityHeader.RowsRoots.Bytes() {
				// recreate the cids using only the computed roots
				cid, err := nodes.CidFromNamespacedSha256(rowRoot)
				if err != nil {
					t.Error(err)
				}

				// check if cid was successfully pinned to IPFS
				_, pinned, err := ipfsAPI.Pin().IsPinned(ctx, path.IpldPath(cid))
				if err != nil {
					t.Error(err)
				}
				if !pinned {
					t.Errorf("failure to pin cid %s to IPFS", cid.String())
				}

				// retrieve the data from IPFS
				_, err = ipfsAPI.Dag().Get(timeoutCtx, cid)
				if err != nil {
					t.Errorf("Root not found: %s", cid.String())
				}
			}
		})
	}
}

func generateRandomMsgOnlyData(msgCount int) Data {
	out := make([]Message, msgCount)
	for i, msg := range generateRandNamespacedRawData(msgCount, NamespaceSize, MsgShareSize-2) {
		out[i] = Message{NamespaceID: msg[:NamespaceSize], Data: msg[NamespaceSize:]}
	}
	return Data{
		Messages: Messages{MessagesList: out},
	}
}

// this code is copy pasted from the plugin, and should likely be exported in the plugin instead
func generateRandNamespacedRawData(total int, nidSize int, leafSize int) [][]byte {
	data := make([][]byte, total)
	for i := 0; i < total; i++ {
		nid := make([]byte, nidSize)
		_, err := rand.Read(nid)
		if err != nil {
			panic(err)
		}
		data[i] = nid
	}

	sortByteArrays(data)
	for i := 0; i < total; i++ {
		d := make([]byte, leafSize)
		_, err := rand.Read(d)
		if err != nil {
			panic(err)
		}
		data[i] = append(data[i], d...)
	}

	return data
}

func sortByteArrays(src [][]byte) {
	sort.Slice(src, func(i, j int) bool { return stdbytes.Compare(src[i], src[j]) < 0 })
}
