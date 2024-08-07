package usecase

import (
	"context"
	"fmt"
	"github.com/gainings/tfirmg/internal/hcl"
	"github.com/gainings/tfirmg/internal/model/resource"
	"github.com/gainings/tfirmg/internal/modulemeta"
	"github.com/gainings/tfirmg/internal/rules"
	_ "github.com/gainings/tfirmg/internal/rules/providers/aws"
	"github.com/gainings/tfirmg/internal/tfstate"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path"
	"strings"
)

// Module is command module to generate terraform code for module resources
func Module(command *cobra.Command, args []string) error {
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

	srcModule, err := command.Flags().GetString("src-module")
	if err != nil {
		return err
	}

	dstModule, err := command.Flags().GetString("dst-module")
	if err != nil {
		return err
	}

	rc := resource.NewResourceCreator(rules.NewRules())
	transformer := tfstate.NewTransformer(rc)
	mu := moduleUsecase{
		transformer: transformer,
		options: moduleUsecaseOpt{
			srcDir:         srcDir,
			dstDir:         dstDir,
			srcTFStatePath: srcTFStatePath,
			srcModule:      srcModule,
			dstModule:      dstModule,
		},
	}

	return mu.run(ctx)
}

type moduleUsecase struct {
	transformer tfstate.Transformer
	options     moduleUsecaseOpt
}

type moduleUsecaseOpt struct {
	srcDir         string
	dstDir         string
	srcTFStatePath string
	srcModule      string
	dstModule      string
}

func (mu moduleUsecase) run(ctx context.Context) error {
	srcModuleJson, err := os.Open(path.Join(mu.options.srcDir, ".terraform/modules/modules.json"))
	if err != nil {
		return fmt.Errorf("failed open file. you must execute terraform init in src root: %s", err.Error())
	}
	srcm, err := modulemeta.Decode(srcModuleJson)
	if err != nil {
		return err
	}
	srcmm := srcm.GetModuleMap()

	srcModuleHCLBlocks, err := hcl.ParseHCL(srcmm[mu.options.srcModule].Dir)
	if err != nil {
		return err
	}
	srcResourceMap := srcModuleHCLBlocks.ResourceNameMap()

	dstModuleJson, err := os.Open(path.Join(mu.options.dstDir, ".terraform/modules/modules.json"))
	if err != nil {
		return fmt.Errorf("failed open file. you must execute terraform init in dst root: %s", err.Error())
	}
	dstm, err := modulemeta.Decode(dstModuleJson)
	if err != nil {
		return err
	}
	dstmm := dstm.GetModuleMap()

	dstModuleHCLBlocks, err := hcl.ParseHCL(dstmm[mu.options.dstModule].Dir)
	if err != nil {
		return err
	}
	dstResourceMap := dstModuleHCLBlocks.ResourceNameMap()

	srcTfstate, err := tfstate.LoadTFState(ctx, mu.options.srcTFStatePath)
	if err != nil {
		return err
	}
	srcTFStateResources := mu.transformer.TransformToResources(srcTfstate)

	notFoundResources, notFoundResourcesInModules := mu.findMissingResourcesInModules(srcResourceMap, srcTFStateResources)

	var onlyCode []resource.Resource
	for key, _ := range dstResourceMap {
		r, exist := notFoundResources[key]
		if exist {
			for _, i := range r {
				onlyCode = append(onlyCode, i)
				continue
			}
		}

		rs, existInModule := notFoundResourcesInModules[key]
		if existInModule {
			for _, r := range rs {
				onlyCode = append(onlyCode, r...)
			}
			continue
		}
	}

	if mu.options.srcDir == mu.options.dstDir {
		mbs := mu.generateMovedBlocks(onlyCode)

		if err := writeToFile(mu.options.srcDir, "moved.tf", mbs.Bytes()); err != nil {
			return err
		}

	} else {
		ibs, rbs := mu.generateImportRemovedBlocks(onlyCode)

		if err := writeToFile(mu.options.dstDir, "import.tf", ibs.Bytes()); err != nil {
			return err
		}

		if err := writeToFile(mu.options.srcDir, "removed.tf", rbs.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// findMissingResourceInModules return
func (mu moduleUsecase) findMissingResourcesInModules(resourceNameMap map[string]struct{}, srcTfstateResources resource.Resources) (map[string]map[string]resource.Resource, map[string]map[string]resource.Resources) {
	notFoundResources := make(map[string]map[string]resource.Resource)
	notFoundResourceInModules := make(map[string]map[string]resource.Resources)

	slog.Debug("---Find missing resource in modules---")
	for _, r := range srcTfstateResources {
		if r.Module == nil {
			continue
		}

		srcMn := fmt.Sprintf("module.%s", mu.options.srcModule)
		if r.Module.Name == srcMn {
			if _, ok := resourceNameMap[r.Name]; !ok {
				k := strings.TrimPrefix(r.Name, mu.options.srcModule)
				if _, exists := notFoundResources[k]; !exists {
					notFoundResources[k] = make(map[string]resource.Resource)
				}
				slog.Debug("This resource removed from terraform code", "target resource", r.Name, "index", r.IndexKey, "module", r.Module.Name)
				notFoundResources[k][r.IndexKey] = r
			}
		} else {
			if _, ok := resourceNameMap[r.Name]; !ok {
				k := strings.Join([]string{strings.TrimPrefix(r.Module.Name, fmt.Sprintf("module.%s.", mu.options.srcModule)), r.Name}, ".")
				if _, exists := notFoundResourceInModules[k]; !exists {
					notFoundResourceInModules[k] = make(map[string]resource.Resources)
				}
				slog.Debug("This resource removed from terraform code", "target resource", r.Name, "index", r.IndexKey, "module", r.Module.Name)
				notFoundResourceInModules[k][r.IndexKey] = append(notFoundResourceInModules[k][r.IndexKey], r)
			}
		}
	}

	return notFoundResources, notFoundResourceInModules
}
func (mu moduleUsecase) generateMovedBlocks(onlyCode []resource.Resource) hcl.MovedBlocks {
	var mbs hcl.MovedBlocks

	for _, r := range onlyCode {
		parts := strings.Split(r.Address.String(), ".")
		parts[1] = mu.options.dstModule

		slog.Debug("Generate Moved Block", "target resource", r.Name, "index", r.IndexKey, "module", r.Module.Name)
		mb := hcl.MovedBlock{
			From: r.Address.String(),
			To:   strings.Join(parts, "."),
		}
		mbs = append(mbs, mb)
	}

	return mbs
}

func (mu moduleUsecase) generateImportRemovedBlocks(onlyCode []resource.Resource) (hcl.ImportBlocks, hcl.RemoveBlocks) {
	var ibs hcl.ImportBlocks
	var rbs hcl.RemoveBlocks

	slog.Debug("---Generate HCL Blocks---")
	for _, r := range onlyCode {
		parts := strings.Split(r.Address.String(), ".")
		dstParts := strings.Split(strings.Join([]string{"module", mu.options.dstModule}, "."), ".")
		for i := 0; i < len(dstParts); i++ {
			if i < len(parts) {
				parts[i] = dstParts[i]
			}
		}

		slog.Debug("Generate Import / Remved Block", "target resource", r.Name, "index", r.IndexKey, "module", r.Module.Name)
		ib := hcl.ImportBlock{
			To: strings.Join(parts, "."),
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
