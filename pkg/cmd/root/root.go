package root

import (
	"github.com/spf13/cobra"
	"gogobox/pkg/cmd/dict"
	versionCmd "gogobox/pkg/cmd/version"
	"gogobox/pkg/cmdutil"
)

func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gogobox <command> <subcommand> [flags]",
		Short: "gogobox",
		Long:  `gogobox is a tool collection for bash environments.`,

		Annotations: map[string]string{
			"version": "0.0.1",
			"website": "www.gogobox.xyz",
		},
	}

	//cmd.Flags().Bool("version", false, "Show version")

	cmd.AddCommand(versionCmd.NewCmdVersion(f))
	cmd.AddCommand(dict.NewCmdDict(f))
	cmd.AddCommand(dict.NewNotebookCmd(f))

	return cmd
}
