package models_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestLog_UnmarshalEmptyTxHash(t *testing.T) {
	t.Parallel()

	input := `{
		"transactionHash": null,
		"transactionIndex": "0x3",
		"address": "0x1aee7c03606fca5035d204c3818d0660bb230e44",
		"blockNumber": "0x8bf99b",
		"topics": ["0xdeadbeefdeadbeedeadbeedeadbeefffdeadbeefdeadbeedeadbeedeadbeefff"],
		"blockHash": "0xdb777676330c067e3c3a6dbfc2d51282cac5bcc1b7a884dd8d85ba72ca1f147e",
		"data": "0xdeadbeef",
		"logIndex": "0x5",
		"transactionLogIndex": "0x3"
	}`

	var log models.Log
	err := json.Unmarshal([]byte(input), &log)
	assert.NoError(t, err)
}

func TestModels_HexToFunctionSelector(t *testing.T) {
	t.Parallel()
	fid := models.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_HexToFunctionSelectorOverflow(t *testing.T) {
	t.Parallel()
	fid := models.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSON(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONLiteral(t *testing.T) {
	t.Parallel()
	literalSelectorBytes := []byte(`"setBytes(bytes)"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(literalSelectorBytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xda359dc8", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONError(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Error(t, err)
}

func TestModels_Header_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		path       string
		wantNumber hexutil.Big
		wantHash   string
	}{
		{
			"parity",
			"../../internal/fixtures/eth/subscription_new_heads_parity.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5114",
		},
		{
			"geth",
			"../../internal/fixtures/eth/subscription_new_heads_geth.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5fff",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var header models.BlockHeader

			data := cltest.MustReadFile(t, test.path)
			value := gjson.Get(string(data), "params.result")
			assert.NoError(t, json.Unmarshal([]byte(value.String()), &header))

			assert.Equal(t, test.wantNumber, header.Number)
			assert.Equal(t, test.wantHash, header.Hash().String())
		})
	}
}

func TestHead_NewHead(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input *big.Int
		want  string
	}{
		{big.NewInt(0), "0"},
		{big.NewInt(0xf), "f"},
		{big.NewInt(0x10), "10"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			num := models.NewHead(test.input, cltest.NewHash())
			assert.Equal(t, test.want, fmt.Sprintf("%x", num.ToInt()))
		})
	}
}

func TestHead_GreaterThan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		left    *models.Head
		right   *models.Head
		greater bool
	}{
		{"nil nil", nil, nil, false},
		{"present nil", cltest.Head(1), nil, true},
		{"nil present", nil, cltest.Head(1), false},
		{"less", cltest.Head(1), cltest.Head(2), false},
		{"equal", cltest.Head(2), cltest.Head(2), false},
		{"greater", cltest.Head(2), cltest.Head(1), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.greater, test.left.GreaterThan(test.right))
		})
	}
}

func TestHead_NextInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		bn   *models.Head
		want *big.Int
	}{
		{"nil", nil, nil},
		{"one", cltest.Head(1), big.NewInt(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.bn.NextInt())
		})
	}
}
