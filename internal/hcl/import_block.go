package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"sort"
)

type ImportBlock struct {
	To       string `hcl:"to"`
	ID       string `hcl:"id"`
	Provider string `hcl:"provider"`
	ForEach  string `hcl:"for_each"`
}

type ImportBlocks []ImportBlock

func (ibs ImportBlocks) Bytes() []byte {
	f := hclwrite.NewEmptyFile()

	sort.SliceStable(ibs, func(i, j int) bool {
		return ibs[i].To < ibs[j].To
	})

	for i, ib := range ibs {
		block := f.Body().AppendNewBlock("import", []string{})
		block.Body().SetAttributeTraversal("to", hcl.Traversal{
			hcl.TraverseRoot{
				Name: ib.To,
			},
		})
		block.Body().SetAttributeValue("id", cty.StringVal(ib.ID))
		if ib.Provider != "" {
			block.Body().SetAttributeValue("provider", cty.StringVal(ib.Provider))
		}
		if ib.ForEach != "" {
			block.Body().SetAttributeValue("for_each", cty.StringVal(ib.ForEach))
		}
		if i < len(ibs)-1 {
			f.Body().AppendNewline()
		}
	}

	return f.Bytes()
}
