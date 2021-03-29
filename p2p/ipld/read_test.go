package ipld

import (
	"bytes"
	"context"
	"crypto/sha256"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core/coreapi"

	coremock "github.com/ipfs/go-ipfs/core/mock"
	format "github.com/ipfs/go-ipld-format"
	"github.com/lazyledger/lazyledger-core/p2p/ipld/plugin/nodes"
	"github.com/lazyledger/lazyledger-core/types"
	"github.com/lazyledger/nmt"
	"github.com/lazyledger/rsmt2d"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeafPath(t *testing.T) {
	type test struct {
		name         string
		index, total uint32
		expected     []string
	}

	// test cases
	tests := []test{
		{"nil", 0, 0, []string(nil)},
		{"0 index 16 total leaves", 0, 16, strings.Split("0/0/0/0", "/")},
		{"1 index 16 total leaves", 1, 16, strings.Split("0/0/0/1", "/")},
		{"9 index 16 total leaves", 9, 16, strings.Split("1/0/0/1", "/")},
		{"15 index 16 total leaves", 15, 16, strings.Split("1/1/1/1", "/")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := leafPath(tt.index, tt.total)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.expected, result)
		},
		)
	}
}

func TestNextPowerOf2(t *testing.T) {
	type test struct {
		input    uint32
		expected uint32
	}
	tests := []test{
		{
			input:    2,
			expected: 2,
		},
		{
			input:    11,
			expected: 8,
		},
		{
			input:    511,
			expected: 256,
		},
		{
			input:    1,
			expected: 1,
		},
		{
			input:    0,
			expected: 0,
		},
	}
	for _, tt := range tests {
		res := nextPowerOf2(tt.input)
		assert.Equal(t, tt.expected, res)
	}
}

func TestGetLeafData(t *testing.T) {
	type test struct {
		name    string
		timeout time.Duration
		rootCid cid.Cid
		leaves  [][]byte
	}

	// create a mock node
	ipfsNode, err := coremock.NewMockNode()
	if err != nil {
		t.Error(err)
	}

	// issue a new API object
	ipfsAPI, err := coreapi.NewCoreAPI(ipfsNode)
	if err != nil {
		t.Error(err)
	}

	// create the context and batch needed for node collection from the tree
	ctx := context.Background()
	batch := format.NewBatch(ctx, ipfsAPI.Dag().Pinning())

	// generate random data for the nmt
	data := generateRandNamespacedRawData(16, types.NamespaceSize, types.ShareSize)

	// create a random tree
	tree, err := createNmtTree(ctx, batch, data)
	if err != nil {
		t.Error(err)
	}

	// calculate the root
	root := tree.Root()

	// commit the data to IPFS
	err = batch.Commit()
	if err != nil {
		t.Error(err)
	}

	// compute the root and create a cid for the root hash
	rootCid, err := nodes.CidFromNamespacedSha256(root.Bytes())
	if err != nil {
		t.Error(err)
	}

	// test cases
	tests := []test{
		{"16 leaves", time.Second, rootCid, data},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			for i, leaf := range tt.leaves {
				data, err := GetLeafData(ctx, tt.rootCid, uint32(i), uint32(len(tt.leaves)), ipfsAPI)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, leaf, data)
			}
		},
		)
	}
}

func TestBlockRecovery(t *testing.T) {
	// adjustedLeafSize describes the size of a leaf that will not get split
	adjustedLeafSize := types.MsgShareSize

	originalSquareWidth := 2
	sharecount := originalSquareWidth * originalSquareWidth
	extendedSquareWidth := originalSquareWidth * originalSquareWidth
	extendedShareCount := extendedSquareWidth * extendedSquareWidth

	// generate test data
	quarterShares := generateRandNamespacedRawData(sharecount, types.NamespaceSize, adjustedLeafSize)
	allShares := generateRandNamespacedRawData(sharecount, types.NamespaceSize, adjustedLeafSize)

	testCases := []struct {
		name string
		// blockData types.Data
		shares    [][]byte
		expectErr bool
		errString string
		d         int // number of shares to delete
	}{
		// missing more shares causes RepairExtendedDataSquare to hang see
		// https://github.com/lazyledger/rsmt2d/issues/21
		{"missing 1/4 shares", quarterShares, false, "", extendedShareCount / 4},
		{"missing all but one shares", allShares, true, "failed to solve data square", extendedShareCount - 1},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			squareSize := uint64(math.Sqrt(float64(len(tc.shares))))

			// create trees for creating roots
			tree := NewErasuredNamespacedMerkleTree(squareSize)
			recoverTree := NewErasuredNamespacedMerkleTree(squareSize)

			eds, err := rsmt2d.ComputeExtendedDataSquare(tc.shares, rsmt2d.RSGF8, tree.Constructor)
			if err != nil {
				t.Error(err)
			}

			// calculate roots using the first complete square
			rowRoots := eds.RowRoots()
			colRoots := eds.ColumnRoots()

			flat := flatten(eds)

			// recover a partially complete square
			reds, err := rsmt2d.RepairExtendedDataSquare(
				rowRoots,
				colRoots,
				removeRandShares(flat, tc.d),
				rsmt2d.RSGF8,
				recoverTree.Constructor,
			)

			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errString)
				return
			}

			require.NoError(t, err)

			// check that the squares are equal
			assert.Equal(t, flatten(eds), flatten(reds))
		})
	}
}

func flatten(eds *rsmt2d.ExtendedDataSquare) [][]byte {
	out := make([][]byte, eds.Width()*eds.Width())
	count := 0
	for i := uint(0); i < eds.Width(); i++ {
		for _, share := range eds.Row(i) {
			out[count] = share
			count++
		}
	}
	return out
}

// nmtcommitment generates the nmt root of some namespaced data
func createNmtTree(
	ctx context.Context,
	batch *format.Batch,
	namespacedData [][]byte,
) (*nmt.NamespacedMerkleTree, error) {
	na := nodes.NewNmtNodeAdder(ctx, batch)
	tree := nmt.New(sha256.New(), nmt.NamespaceIDSize(types.NamespaceSize), nmt.NodeVisitor(na.Visit))
	for _, leaf := range namespacedData {
		err := tree.Push(leaf[:types.NamespaceSize], leaf[types.NamespaceSize:])
		if err != nil {
			return tree, err
		}
	}

	return tree, nil
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
	sort.Slice(src, func(i, j int) bool { return bytes.Compare(src[i], src[j]) < 0 })
}

// removes d shares from data
func removeRandShares(data [][]byte, d int) [][]byte {
	count := len(data)
	// remove shares randomly
	for i := 0; i < d; {
		ind := rand.Intn(count)
		if len(data[ind]) == 0 {
			continue
		}
		data[ind] = nil
		i++
	}
	return data
}
