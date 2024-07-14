package usecase

import (
	"context"
	"github.com/gainings/tfirmg/internal/hcl"
	"github.com/gainings/tfirmg/internal/model/resource"
	"github.com/gainings/tfirmg/internal/rules"
	_ "github.com/gainings/tfirmg/internal/rules/providers/aws"
	"github.com/gainings/tfirmg/internal/tfstate"
	"github.com/spf13/cobra"
	"log/slog"
)

// Root is command root to generate terraform code for root resources
func Root(command *cobra.Command, args []string) error {
	ctx := context.Background()

	srcDir, err := command.Flags().GetString("src-dir")
	if err != nil {
		return err
	}

	dstDir, err := command.Flags().GetString("dst-dir")
	if err != nil {
		return err
	}

	srcTFStatePath, err := command.Flags().GetString("src-tfstate-path")
	if err != nil {
		return err
	}

	rc := resource.NewResourceCreator(rules.NewRules())
	transformer := tfstate.NewTransformer(rc)
	ru := rootUsecase{
		transformer: transformer,
		options: rootUsecaseOpt{
			srcDir:         srcDir,
			dstDir:         dstDir,
			srcTFStatePath: srcTFStatePath,
		},
	}

	return ru.run(ctx)
}

type rootUsecase struct {
	transformer tfstate.Transformer
	options     rootUsecaseOpt
}

type rootUsecaseOpt struct {
	srcDir         string
	dstDir         string
	srcTFStatePath string
}

func (ru rootUsecase) run(ctx context.Context) error {
	srcHCLBlocks, err := hcl.ParseHCL(ru.options.srcDir)
	if err != nil {
		return err
	}
	srcResourceMap := srcHCLBlocks.ResourceNameMap()

	dstHCLBlocks, err := hcl.ParseHCL(ru.options.dstDir)
	if err != nil {
		return err
	}
	dstResourceMap := dstHCLBlocks.ResourceNameMap()

	srcTfstate, err := tfstate.LoadTFState(ctx, ru.options.srcTFStatePath)
	if err != nil {
		return err
	}
	srcTFStateResources := ru.transformer.TransformToResources(srcTfstate)

	notFoundResources, notFoundResourceInModules := findMissingResources(srcResourceMap, srcTFStateResources)

	var onlyCode []resource.Resource
	for key, _ := range dstResourceMap {
		r, exist := notFoundResources[key]
		if exist {
			onlyCode = append(onlyCode, r)
			continue
		}
		rs, existInModule := notFoundResourceInModules[key]
		if existInModule {
			for _, r := range rs {
				onlyCode = append(onlyCode, r)
			}
			continue
		}
	}

	ibs, rbs := generateHCLBlocks(onlyCode)

	if err := writeToFile(ru.options.dstDir, "import.tf", ibs.Bytes()); err != nil {
		return err
	}

	if err := writeToFile(ru.options.srcDir, "removed.tf", rbs.Bytes()); err != nil {
		return err
	}

	return nil
}

func findMissingResources(resourceNameMap map[string]struct{}, srcTfstateResources resource.Resources) (map[string]resource.Resource, map[string][]resource.Resource) {
	notFoundResources := make(map[string]resource.Resource)
	notFoundResourceInModules := make(map[string][]resource.Resource)

	slog.Debug("---Find missing resources---")
	for _, r := range srcTfstateResources {
		if r.Module != nil {
			if _, ok := resourceNameMap[r.Module.Name]; !ok {
				slog.Debug("This resource removed from terraform code", "target resource", r.Name, "index", r.IndexKey, "module", r.Module.Name)
				notFoundResourceInModules[r.Module.Name] = append(notFoundResourceInModules[r.Module.Name], r)
			}
		} else {
			if _, ok := resourceNameMap[r.Name]; !ok {
				slog.Debug("This resource removed from terraform code", "target resource", r.Name, "index", r.IndexKey)
				notFoundResources[r.Name] = r
			}
		}
	}

	return notFoundResources, notFoundResourceInModules
}

func generateHCLBlocks(onlyCode []resource.Resource) (hcl.ImportBlocks, hcl.RemoveBlocks) {
	var ibs hcl.ImportBlocks
	var rbs hcl.RemoveBlocks

	slog.Debug("---Generate HCL Blocks---")
	for _, r := range onlyCode {
		slog.Debug("Generate Import / Remved Block", "target resource", r.Name, "index", r.IndexKey)
		ib := hcl.ImportBlock{
			To: r.Address.String(),
			ID: r.ID.String(),
		}
		ibs = append(ibs, ib)

		rb := hcl.RemoveBlock{
			From:    r.Address.String(),
			Destroy: false,
		}
		rbs = append(rbs, rb)
	}

	return ibs, rbs
}
