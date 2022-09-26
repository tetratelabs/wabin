package binary

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/wabin/wasm"
)

func TestEncodeCustom(t *testing.T) {
	tests := []struct {
		name     string
		custom   *wasm.CustomSection
		expected []byte
	}{
		{
			name: "custom section",
			custom: &wasm.CustomSection{
				Name: "test",
				Data: []byte("12345"),
			},
			expected: []byte{
				4, 't', 'e', 's', 't',
				'1', '2', '3', '4', '5',
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			bytes := encodeCustomSection(tc.custom)
			require.Equal(t, tc.expected, bytes)
		})
	}
}
