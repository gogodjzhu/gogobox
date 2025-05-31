package version

import (
	"fmt"

	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: false,
		Short:  "Print the version number of gogobox",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(f.IOStreams.Out, cmd.Root().Annotations["version"])
		},
	}

	return cmd
}
