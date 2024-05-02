package hcl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMovedBlocksBytes(t *testing.T) {
	tests := []struct {
		name     string
		mbs      MovedBlocks
		expected string
	}{
		{
			name: "multiple blocks sorted",
			mbs: MovedBlocks{
				{From: "module.alpha", To: "module.omega"},
				{From: "module.beta", To: "module.gamma"},
			},
			expected: `moved {
  from = module.alpha
  to   = module.omega
}

moved {
  from = module.beta
  to   = module.gamma
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := string(tc.mbs.Bytes())
			assert.Equal(t, tc.expected, output)
		})
	}
}
