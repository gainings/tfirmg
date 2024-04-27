package hcl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImportBlocksBytes(t *testing.T) {
	testCases := []struct {
		name     string
		blocks   ImportBlocks
		expected string
	}{
		{
			name: "単一のresourceのBlockがsortされて正常に出力される",
			blocks: ImportBlocks{
				ImportBlock{To: "aws_instance.example2", ID: "i-abcd1235"},
				ImportBlock{To: "aws_instance.example1", ID: "i-abcd1234", Provider: "aws"},
			},
			expected: `import {
  to       = aws_instance.example1
  id       = "i-abcd1234"
  provider = "aws"
}

import {
  to = aws_instance.example2
  id = "i-abcd1235"
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := string(tc.blocks.Bytes())
			assert.Equal(t, tc.expected, result)
		})
	}
}
