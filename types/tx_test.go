package types

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tmrand "github.com/tendermint/tendermint/libs/rand"
	ctest "github.com/tendermint/tendermint/libs/test"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func makeTxs(cnt, size int) Txs {
	txs := make(Txs, cnt)
	for i := 0; i < cnt; i++ {
		txs[i] = tmrand.Bytes(size)
	}
	return txs
}

func randInt(low, high int) int {
	off := tmrand.Int() % (high - low)
	return low + off
}

func TestTxIndex(t *testing.T) {
	for i := 0; i < 20; i++ {
		txs := makeTxs(15, 60)
		for j := 0; j < len(txs); j++ {
			tx := txs[j]
			idx := txs.Index(tx)
			assert.Equal(t, j, idx)
		}
		assert.Equal(t, -1, txs.Index(nil))
		assert.Equal(t, -1, txs.Index(Tx("foodnwkf")))
	}
}

func TestTxIndexByHash(t *testing.T) {
	for i := 0; i < 20; i++ {
		txs := makeTxs(15, 60)
		for j := 0; j < len(txs); j++ {
			tx := txs[j]
			idx := txs.IndexByHash(tx.Hash())
			assert.Equal(t, j, idx)
		}
		assert.Equal(t, -1, txs.IndexByHash(nil))
		assert.Equal(t, -1, txs.IndexByHash(Tx("foodnwkf").Hash()))
	}
}

func TestValidTxProof(t *testing.T) {
	cases := []struct {
		txs Txs
	}{
		{Txs{{1, 4, 34, 87, 163, 1}}},
		{Txs{{5, 56, 165, 2}, {4, 77}}},
		{Txs{Tx("foo"), Tx("bar"), Tx("baz")}},
		{makeTxs(20, 5)},
		{makeTxs(7, 81)},
		{makeTxs(61, 15)},
	}

	for h, tc := range cases {
		txs := tc.txs
		root := txs.Hash()
		// make sure valid proof for every tx
		for i := range txs {
			tx := []byte(txs[i])
			proof := txs.Proof(i)
			assert.EqualValues(t, i, proof.Proof.Index, "%d: %d", h, i)
			assert.EqualValues(t, len(txs), proof.Proof.Total, "%d: %d", h, i)
			assert.EqualValues(t, root, proof.RootHash, "%d: %d", h, i)
			assert.EqualValues(t, tx, proof.Data, "%d: %d", h, i)
			assert.EqualValues(t, txs[i].Hash(), proof.Leaf(), "%d: %d", h, i)
			assert.Nil(t, proof.Validate(root), "%d: %d", h, i)
			assert.NotNil(t, proof.Validate([]byte("foobar")), "%d: %d", h, i)

			// read-write must also work
			var (
				p2  TxProof
				pb2 tmproto.TxProof
			)
			pbProof := proof.ToProto()
			bin, err := pbProof.Marshal()
			require.NoError(t, err)

			err = pb2.Unmarshal(bin)
			require.NoError(t, err)

			p2, err = TxProofFromProto(pb2)
			if assert.Nil(t, err, "%d: %d: %+v", h, i, err) {
				assert.Nil(t, p2.Validate(root), "%d: %d", h, i)
			}
		}
	}
}

func TestTxProofUnchangable(t *testing.T) {
	// run the other test a bunch...
	for i := 0; i < 40; i++ {
		testTxProofUnchangable(t)
	}
}

func testTxProofUnchangable(t *testing.T) {
	// make some proof
	txs := makeTxs(randInt(2, 100), randInt(16, 128))
	root := txs.Hash()
	i := randInt(0, len(txs)-1)
	proof := txs.Proof(i)

	// make sure it is valid to start with
	assert.Nil(t, proof.Validate(root))
	pbProof := proof.ToProto()
	bin, err := pbProof.Marshal()
	require.NoError(t, err)

	// try mutating the data and make sure nothing breaks
	for j := 0; j < 500; j++ {
		bad := ctest.MutateByteSlice(bin)
		if !bytes.Equal(bad, bin) {
			assertBadProof(t, root, bad, proof)
		}
	}
}

// This makes sure that the proof doesn't deserialize into something valid.
func assertBadProof(t *testing.T, root []byte, bad []byte, good TxProof) {

	var (
		proof   TxProof
		pbProof tmproto.TxProof
	)
	err := pbProof.Unmarshal(bad)
	if err == nil {
		proof, err = TxProofFromProto(pbProof)
		if err == nil {
			err = proof.Validate(root)
			if err == nil {
				// XXX Fix simple merkle proofs so the following is *not* OK.
				// This can happen if we have a slightly different total (where the
				// path ends up the same). If it is something else, we have a real
				// problem.
				assert.NotEqual(t, proof.Proof.Total, good.Proof.Total, "bad: %#v\ngood: %#v", proof, good)
			}
		}
	}
}

func TestUnwrapMalleatedTx(t *testing.T) {
	// perform a simple test for being unable to decode a non
	// malleated transaction
	tx := Tx{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_, _, ok := UnwrapMalleatedTx(tx)
	require.False(t, ok)

	// create a proto message that used to be decoded when it shouldn't have
	randomBlock := MakeBlock(
		1,
		[]Tx{tx},
		nil,
		nil,
		[]Message{
			{
				NamespaceID: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				Data:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		&Commit{},
	)

	protoB, err := randomBlock.ToProto()
	require.NoError(t, err)

	rawBlock, err := protoB.Marshal()
	require.NoError(t, err)

	// due to protobuf not actually requiring type compatibility
	// we need to make sure that there is some check
	_, _, ok = UnwrapMalleatedTx(rawBlock)
	require.False(t, ok)

	pHash := sha256.Sum256(rawBlock)
	MalleatedTx, err := WrapMalleatedTx(pHash[:], rawBlock)
	require.NoError(t, err)

	// finally, ensure that the unwrapped bytes are identical to the input
	unwrappedHash, unwrapped, ok := UnwrapMalleatedTx(MalleatedTx)
	require.True(t, ok)
	assert.Equal(t, 32, len(unwrappedHash))
	require.Equal(t, rawBlock, []byte(unwrapped))
}
