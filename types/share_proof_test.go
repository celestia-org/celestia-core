package types

import (
	"testing"

	"github.com/cometbft/cometbft/pkg/consts"
	"github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/assert"
)

func TestShareProofValidate(t *testing.T) {
	type testCase struct {
		name    string
		sp      ShareProof
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "empty share proof returns error",
			sp:      ShareProof{},
			wantErr: true,
		},
		{
			name:    "valid share proof returns no error",
			sp:      validShareProof(),
			wantErr: false,
		},
		{
			name:    "share proof with mismatched number of share proofs returns error",
			sp:      mismatchedShareProofs(),
			wantErr: true,
		},
		{
			name:    "share proof with mismatched number of shares returns error",
			sp:      mismatchedShares(),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.sp.Validate()
			if tc.wantErr {
				assert.Error(t, got)
				return
			}
			assert.NoError(t, got)
		})
	}
}

func TestShareProofVerify(t *testing.T) {
	type testCase struct {
		name       string
		sp         ShareProof
		root       []byte
		isVerified bool
	}

	testCases := []testCase{
		{
			name:       "empty share proof returns false",
			sp:         ShareProof{},
			root:       root,
			isVerified: false,
		},
		{
			name:       "valid share proof returns true",
			sp:         validShareProof(),
			root:       root,
			isVerified: true,
		},
		{
			name:       "share proof with mismatched number of share proofs returns false",
			sp:         mismatchedShareProofs(),
			root:       root,
			isVerified: false,
		},
		{
			name:       "share proof with mismatched number of shares returns false",
			sp:         mismatchedShares(),
			root:       root,
			isVerified: false,
		},
		{
			name:       "valid share proof with incorrect root returns false",
			sp:         validShareProof(),
			root:       incorrectRoot,
			isVerified: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isVerified, tc.sp.VerifyProof(tc.root))
		})
	}
}

func mismatchedShareProofs() ShareProof {
	sp := validShareProof()
	sp.ShareProofs = []*types.NMTProof{}
	return sp
}

func mismatchedShares() ShareProof {
	sp := validShareProof()
	sp.Data = [][]byte{}
	return sp
}

// validShareProof returns a valid ShareProof for a single share. This test data
// was copied from celestia-app's pkg/proof/proof_test.go
// TestNewShareInclusionProof: "1 transaction share"
func validShareProof() ShareProof {
	return ShareProof{
		Data: [][]uint8{{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x62, 0xc, 0x0, 0x0, 0x0, 0x2a, 0xf4, 0x3, 0xff, 0xe8, 0x78, 0x6c, 0x48, 0x84, 0x9, 0x5, 0x5, 0x79, 0x8f, 0x29, 0x67, 0xa2, 0xe1, 0x8d, 0x2f, 0xdc, 0xf2, 0x60, 0xe4, 0x62, 0x71, 0xf9, 0xae, 0x92, 0x83, 0x3a, 0x7f, 0xf3, 0xc6, 0x14, 0xb4, 0x17, 0xfc, 0x64, 0x4b, 0x89, 0x18, 0x5e, 0x22, 0x4b, 0x0, 0x82, 0xeb, 0x67, 0x5b, 0x51, 0x43, 0x4e, 0xc3, 0x42, 0x48, 0xc1, 0xfd, 0x88, 0x71, 0xcb, 0xee, 0xf3, 0x92, 0x20, 0x9c, 0x15, 0xc0, 0x4f, 0x11, 0xa4, 0x5, 0xd0, 0xdf, 0xb8, 0x25, 0x60, 0x58, 0xae, 0x2, 0x2d, 0x78, 0xf8, 0x1f, 0x67, 0xeb, 0x88, 0x58, 0x5d, 0x5a, 0x4a, 0x74, 0xe7, 0xdf, 0x38, 0x6a, 0xa4, 0x3f, 0x62, 0xd6, 0x3d, 0x17, 0xd2, 0x7e, 0x92, 0x9c, 0x4a, 0xd0, 0x2b, 0x55, 0x49, 0x3b, 0xa7, 0x5a, 0x29, 0xd5, 0x6b, 0x91, 0xde, 0xfe, 0x5b, 0x39, 0x88, 0xc5, 0xbb, 0x91, 0x16, 0xf6, 0x47, 0xec, 0x8, 0x3, 0x2a, 0x1e, 0x6e, 0x4b, 0x27, 0x34, 0x90, 0x38, 0x46, 0x6e, 0xce, 0x35, 0xdf, 0xd6, 0x1e, 0x1a, 0xf2, 0xf0, 0x6e, 0xa0, 0xfe, 0x84, 0x51, 0xf2, 0xc1, 0x32, 0xd, 0x89, 0x17, 0x5f, 0x4c, 0xab, 0x81, 0xd4, 0x44, 0x5a, 0x55, 0xdb, 0xe5, 0xa7, 0x3c, 0x42, 0xb6, 0xb3, 0x20, 0xc4, 0x81, 0x75, 0x8, 0x5e, 0x39, 0x21, 0x51, 0x4c, 0x93, 0x2c, 0x7c, 0xb3, 0xd0, 0x37, 0xf9, 0x6a, 0xab, 0x93, 0xf0, 0x3f, 0xa2, 0x44, 0x1f, 0x63, 0xae, 0x96, 0x4e, 0x26, 0x7a, 0x1f, 0x18, 0x5b, 0x28, 0x4d, 0x24, 0xe8, 0x98, 0x56, 0xbf, 0x98, 0x44, 0x23, 0x17, 0x85, 0x22, 0x38, 0x56, 0xeb, 0xf3, 0x4e, 0x87, 0x1e, 0xc1, 0x51, 0x6, 0x71, 0xa7, 0xa9, 0x45, 0xef, 0xc7, 0x89, 0x5c, 0xed, 0x68, 0xbd, 0x43, 0x2f, 0xe6, 0xf1, 0x56, 0xef, 0xf, 0x4f, 0x57, 0xaa, 0x8c, 0x5c, 0xbd, 0x21, 0xb4, 0xaa, 0x15, 0x71, 0x6a, 0xdc, 0x12, 0xda, 0xee, 0xd9, 0x19, 0xbc, 0x17, 0xa2, 0x49, 0xd6, 0xbe, 0xd2, 0xc6, 0x6a, 0xbc, 0x53, 0xe4, 0x28, 0xd4, 0xeb, 0xe9, 0x9b, 0xd6, 0x85, 0x89, 0xb9, 0xe8, 0xa2, 0x70, 0x40, 0xad, 0xb1, 0x1a, 0xa0, 0xb1, 0xb5, 0xee, 0xde, 0x6d, 0xa9, 0x2a, 0x4b, 0x6, 0xd1, 0xfa, 0x67, 0x13, 0xac, 0x7d, 0x9a, 0x81, 0xc6, 0xef, 0x78, 0x42, 0x18, 0xf, 0x7b, 0xaf, 0x50, 0xa7, 0xdb, 0xb6, 0xde, 0xab, 0x3, 0xdc, 0x5, 0x14, 0x5f, 0x9, 0xdb, 0x81, 0xe3, 0x72, 0x2, 0x61, 0x23, 0x77, 0x12, 0x82, 0xfc, 0x9, 0x43, 0xfb, 0xd6, 0x38, 0x53, 0xfd, 0x77, 0xe, 0x17, 0xcc, 0x93, 0x5e, 0x4e, 0x60, 0x87, 0xda, 0xbd, 0xfc, 0x86, 0xdd, 0xb1, 0xd6, 0x74, 0x41, 0x71, 0x24, 0xda, 0x1, 0x3f, 0x11, 0x17, 0x9e, 0x54, 0x66, 0xb6, 0xc4, 0x9a, 0xb8, 0x59, 0xb9, 0x13, 0x4e, 0xed, 0x8, 0xe5, 0x99, 0x27, 0xa0, 0x6b, 0x1, 0x6c, 0x8a, 0xbf, 0x20, 0x3d, 0x75, 0xd5, 0x7e, 0xea, 0xe0, 0xef, 0x7f, 0xfe, 0xa8, 0xaf, 0x76, 0xad, 0x30, 0x55, 0x65, 0x9d, 0xbe, 0x30, 0x32, 0x9f, 0x3b, 0xb7, 0xa1, 0x5c, 0x98, 0xef, 0xe1, 0xe4, 0x33, 0x1a, 0x56, 0x5a, 0x22, 0xd1, 0x38, 0x9b, 0xee, 0xfa, 0x11, 0x6f, 0xa7, 0xd7, 0x6, 0x17, 0xdc, 0xc6, 0x4d, 0xbd, 0x3f, 0x3c, 0xe6, 0xac, 0x54, 0x70, 0xda, 0x11, 0xdb, 0x87, 0xe2, 0xc2, 0x26, 0x7e, 0x48, 0x3b, 0xda, 0xf4, 0x98, 0x3c, 0x51}},
		ShareProofs: []*types.NMTProof{
			{
				Start:    0,
				End:      1,
				Nodes:    [][]uint8{{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x27, 0x3a, 0x5f, 0x16, 0x36, 0xa3, 0xce, 0x1c, 0x17, 0x58, 0x7e, 0xb8, 0xaa, 0xc8, 0x5e, 0x58, 0x9e, 0xa9, 0x36, 0x3c, 0x3d, 0x5c, 0xb5, 0xc2, 0xf0, 0x26, 0x1a, 0x9a, 0x13, 0xcd, 0x59, 0xb2}, {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x55, 0xe5, 0x43, 0x2e, 0xa2, 0x32, 0x84, 0x75, 0x8a, 0x88, 0x8d, 0x7c, 0x27, 0xdc, 0x2e, 0x13, 0x1e, 0x44, 0xc4, 0xe7, 0x51, 0x64, 0xe5, 0xe4, 0xf4, 0x7d, 0x4, 0xb8, 0x10, 0x3b, 0x72, 0xa5}, {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x4d, 0xeb, 0x2a, 0x3c, 0x56, 0x98, 0x49, 0xdb, 0x61, 0x54, 0x12, 0xee, 0xb, 0xeb, 0x29, 0xf8, 0xc9, 0x71, 0x9c, 0xf7, 0x28, 0xbb, 0x7a, 0x85, 0x70, 0xa1, 0x81, 0xc8, 0x5f, 0x6a, 0x63, 0x59}, {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xf0, 0xb5, 0x59, 0x71, 0xba, 0x6a, 0xf, 0xd1, 0xf, 0x2e, 0x79, 0xd4, 0xdc, 0xfb, 0x93, 0x94, 0x58, 0x3d, 0xd9, 0xef, 0xe2, 0x2b, 0xd4, 0xe3, 0x71, 0xbd, 0xd4, 0xd9, 0xc2, 0xc4, 0xef, 0xd1}, {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x5, 0x8f, 0xf0, 0x4e, 0x81, 0x8e, 0xc7, 0x2f, 0x35, 0xec, 0x9, 0xdf, 0xf1, 0x41, 0xd5, 0x5a, 0x2f, 0xa3, 0xa0, 0xe5, 0x8d, 0x83, 0x70, 0xf2, 0x11, 0xea, 0xc2, 0xa3, 0x4a, 0x7a, 0xc5, 0x17}, {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x96, 0x6d, 0x3f, 0x7b, 0xf5, 0xef, 0x38, 0x4b, 0xa5, 0x38, 0x98, 0x7e, 0x3b, 0x4e, 0x12, 0x21, 0xcb, 0xd7, 0xff, 0xd6, 0xf3, 0x7d, 0xf, 0x8a, 0x57, 0xfe, 0x5, 0x5, 0xb6, 0x62, 0xa6, 0xae}},
				LeafHash: []uint8(nil),
			},
		},
		NamespaceID:      consts.TxNamespaceID,
		RowProof:         validRowProof(),
		NamespaceVersion: uint32(0),
	}
}
