package hcl

import (
	"fmt"
	"github.com/spf13/afero"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Hcl struct {
	Blocks hclsyntax.Blocks
}

var aferoFS = afero.NewOsFs()

func ParseHCL(dirPath string) (Hcl, error) {
	parser := hclparse.NewParser()
	var blocks hclsyntax.Blocks
	afero.Walk(aferoFS, dirPath, func(path string, info fs.FileInfo, err error) error {
		if info != nil && info.IsDir() && path != dirPath {
			return filepath.SkipDir // skip subdirectory
		}
		if !strings.HasSuffix(path, ".tf") {
			return nil
		}
		c, err := afero.ReadFile(aferoFS, path)
		if err != nil {
			return fmt.Errorf("Failed to read file: %s", err)
		}
		file, diags := parser.ParseHCL(c, path)
		if diags.HasErrors() {
			return fmt.Errorf("Failed to parse file: %s", diags.Error())
		}
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			return fmt.Errorf("File body is not *hclsyntax.Body type")
		}
		blocks = append(blocks, body.Blocks...)
		return nil
	})
	return Hcl{Blocks: blocks}, nil
}

type BackendInfo struct {
	Type       string
	Attributes map[string]string
}

func (h Hcl) BackendURL() (*BackendInfo, error) {
	for _, block := range h.Blocks {
		if block.Type == "terraform" {
			for _, innerBlock := range block.Body.Blocks {
				if innerBlock.Type == "backend" {
					backendInfo := &BackendInfo{
						Type:       innerBlock.Labels[0],
						Attributes: make(map[string]string),
					}
					for _, attr := range innerBlock.Body.Attributes {
						val, diags := attr.Expr.Value(nil)
						if diags.HasErrors() {
							return nil, fmt.Errorf("error reading attribute %s: %s", attr.Name, diags.Error())
						}
						backendInfo.Attributes[attr.Name] = fmt.Sprintf("%v", val.AsString())
					}
					return backendInfo, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("no backend block found")
}

// ResourceNameMap is return map store terraform resource name
func (h Hcl) ResourceNameMap() map[string]struct{} {
	resources := make(map[string]struct{})
	for _, block := range h.Blocks {
		if block.Type == "resource" {
			resources[fmt.Sprintf("%s.%s", block.Labels[0], block.Labels[1])] = struct{}{}
		} else if block.Type == "module" {
			resources[fmt.Sprintf("module.%s", block.Labels[0])] = struct{}{}
		}
	}
	return resources
}
