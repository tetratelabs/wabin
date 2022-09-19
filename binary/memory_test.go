package binary

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/wabin/wasm"
)

func TestMemoryType(t *testing.T) {
	zero := uint32(0)

	tests := []struct {
		name     string
		input    *wasm.Memory
		expected []byte
	}{
		{
			name:     "min 0",
			input:    &wasm.Memory{},
			expected: []byte{0, 0},
		},
		{
			name:     "min 0, max 0",
			input:    &wasm.Memory{Max: zero, IsMaxEncoded: true},
			expected: []byte{0x1, 0, 0},
		},
		{
			name:     "min=max",
			input:    &wasm.Memory{Min: 1, Max: 1, IsMaxEncoded: true},
			expected: []byte{0x1, 1, 1},
		},
		{
			name:     "min 0, max largest",
			input:    &wasm.Memory{Max: uint32(65536), IsMaxEncoded: true},
			expected: []byte{0x1, 0, 0x80, 0x80, 0x4},
		},
		{
			name:     "min largest max largest",
			input:    &wasm.Memory{Min: uint32(65536), Max: uint32(65536), IsMaxEncoded: true},
			expected: []byte{0x1, 0x80, 0x80, 0x4, 0x80, 0x80, 0x4},
		},
	}

	for _, tt := range tests {
		tc := tt

		b := encodeMemory(tc.input)
		t.Run(fmt.Sprintf("encode %s", tc.name), func(t *testing.T) {
			require.Equal(t, tc.expected, b)
		})

		t.Run(fmt.Sprintf("decode %s", tc.name), func(t *testing.T) {
			binary, err := decodeMemory(bytes.NewReader(b))
			require.NoError(t, err)
			require.Equal(t, binary, tc.input)
		})
	}
}

func TestDecodeMemoryType_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedErr string
	}{
		{
			name:        "max < min",
			input:       []byte{0x1, 0x80, 0x80, 0x4, 0},
			expectedErr: "min 65536 pages (4 Gi) > max 0 pages (0 Ki)",
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			_, err := decodeMemory(bytes.NewReader(tc.input))
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}
