package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"sort"
	"strings"
)

type RemoveBlock struct {
	From    string `hcl:"from"`
	Destroy bool   `hcl:"destroy"`
}

type RemoveBlocks []RemoveBlock

func (rbs RemoveBlocks) Bytes() []byte {
	f := hclwrite.NewEmptyFile()

	sort.SliceStable(rbs, func(i, j int) bool {
		return rbs[i].From < rbs[j].From
	})
	for i, rb := range rbs {
		sf := strings.Split(string(rb.From), "[")
		block := f.Body().AppendNewBlock("removed", []string{})
		block.Body().SetAttributeTraversal("from", hcl.Traversal{
			hcl.TraverseRoot{
				Name: sf[0],
			},
		})
		lifecycleBlock := block.Body().AppendNewBlock("lifecycle", nil)
		lifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(rb.Destroy))
		if i < len(rbs)-1 {
			f.Body().AppendNewline()
		}
	}

	return f.Bytes()
}
