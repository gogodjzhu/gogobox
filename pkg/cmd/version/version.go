package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
)

func NewCmdVersion(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(f.IOStreams.Out, cmd.Root().Annotations["version"])
		},
	}

	return cmd
}
