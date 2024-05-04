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
	rc.Flags().String("dst-dir", "", "dst directory that terraform files")
	rc.Flags().String("src-tfstate-path", "", "tfstate path")

	mc := runModuleCmd()
	mc.Flags().String("src-dir", "", "src directory that exists terraform files")
	mc.Flags().String("dst-dir", "", "dst directory that terraform files")
	mc.Flags().String("src-tfstate-path", "", "tfstate path")
	mc.Flags().String("src-module", "", "src module name that resource moved")
	mc.Flags().String("dst-module", "", "dst module name that will be moved")

	BaseCmd.AddCommand(
		rc,
		mc,
	)
}
func runRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "root",
		Short: "Generate import and moved block for root module",
		RunE:  usecase.Root,
	}
	cmd.Flags().SetInterspersed(false)
	return cmd
}
func runModuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Generate import, removed and moved block for specific modules",
		RunE:  usecase.Module,
	}
	cmd.Flags().SetInterspersed(false)
	return cmd
}
