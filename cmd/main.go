package cmd

import (
	"github.com/gainings/tfirmg/internal/usecase"
	"github.com/spf13/cobra"
)

// RootCmd is cmd that handling terraform root resource
var BaseCmd = &cobra.Command{
	Long: "this tool provide any tool for develop utilities",
}

var version string

var (
	port int
)

func init() {
	cobra.OnInitialize()
	BaseCmd.Run = func(cmd *cobra.Command, args []string) {
		_ = BaseCmd.Help()
	}
	rc := runRootCmd()
	rc.Flags().String("src-dir", "", "src directory that exists terraform files")
	rc.Flags().String("dst-dir", "", "dst directory that  terraform files")
	rc.Flags().String("src-tfstate-path", "", "tfstate path")

	mc := runModuleCmd()
	mc.Flags().String("src-dir", "", "src directory that exists terraform files")
	mc.Flags().String("dst-dir", "", "dst directory that  terraform files")
	mc.Flags().String("src-tfstate-path", "", "tfstate path")
	mc.Flags().String("src-module", "", "src directory that exists terraform files")
	mc.Flags().String("dst-module", "", "dst directory that  terraform files")

	BaseCmd.AddCommand(
		rc,
		mc,
	)
}
func runRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate import and moved block",
		RunE:  usecase.Root,
	}
	cmd.Flags().SetInterspersed(false)
	return cmd
}
func runModuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Generate import and moved block",
		RunE:  usecase.Module,
	}
	cmd.Flags().SetInterspersed(false)
	return cmd
}
