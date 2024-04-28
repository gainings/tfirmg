package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"sort"
)

type MovedBlock struct {
	From string `hcl:"from"`
	To   string `hcl:"to"`
}

type MovedBlocks []MovedBlock

func (mbs MovedBlocks) Bytes() []byte {
	f := hclwrite.NewEmptyFile()

	sort.SliceStable(mbs, func(i, j int) bool {
		return mbs[i].From < mbs[j].From
	})
	for i, mb := range mbs {
		block := f.Body().AppendNewBlock("moved", []string{})
		block.Body().SetAttributeTraversal("from", hcl.Traversal{
			hcl.TraverseRoot{
				Name: mb.From,
			},
		})
		block.Body().SetAttributeTraversal("to", hcl.Traversal{
			hcl.TraverseRoot{
				Name: mb.To,
			},
		})
		if i < len(mbs)-1 {
			f.Body().AppendNewline()
		}
	}

	return f.Bytes()
}
