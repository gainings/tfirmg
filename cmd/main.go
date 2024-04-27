package cmd

import (
	"github.com/gainings/tfirg/internal/usecase"
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

	BaseCmd.AddCommand(
		rc,
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
