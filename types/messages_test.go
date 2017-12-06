package types

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	b, err := json.Marshal(&ResponseDeliverTx{})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(b), "code"))
}

func TestWriteReadMessage(t *testing.T) {
	cases := []proto.Message{
		&Header{
			NumTxs: 4,
		},
		// TODO: add the rest
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		err := WriteMessage(c, buf)
		assert.Nil(t, err)

		msg := new(Header)
		err = ReadMessage(buf, msg)
		assert.Nil(t, err)

		assert.Equal(t, c, msg)
	}
}
