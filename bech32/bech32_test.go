package bech32_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/tendermint/tmlibs/bech32"
)

func TestEncodeAndDecode(t *testing.T) {

	sum := sha256.Sum256([]byte("hello world\n"))

	bech, err := bech32.ConvertAndEncode("shasum", sum[:])

	if err != nil {
		t.Error(err)
	}
	hrp, data, err := bech32.DecodeAndConvert(bech)

	if err != nil {
		t.Error(err)
	}
	if hrp != "shasum" {
		t.Error("Invalid hrp")
	}
	if bytes.Compare(data, sum[:]) != 0 {
		t.Error("Invalid decode")
	}
}
