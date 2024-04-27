package hcl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveBlocks_Bytes(t *testing.T) {
	tests := []struct {
		name     string
		input    RemoveBlocks
		expected string
	}{
		{
			name: "正常系",
			input: RemoveBlocks{
				RemoveBlock{From: "aws_instance.example1", Destroy: true},
				RemoveBlock{From: "module.my_usecase", Destroy: true},
			},
			expected: `removed {
  from = aws_instance.example1
  lifecycle {
    destroy = true
  }
}

removed {
  from = module.my_usecase
  lifecycle {
    destroy = true
  }
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := string(tc.input.Bytes())
			assert.Equal(t, tc.expected, result)
		})
	}
}
