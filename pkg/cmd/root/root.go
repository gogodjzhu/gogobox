package root

import (
	"github.com/spf13/cobra"
	"gogobox/pkg/cmd/dict"
	versionCmd "gogobox/pkg/cmd/version"
	"gogobox/pkg/cmdutil"
)

func NewCmdRoot(f *cmdutil.Factory) (*cobra.Command, error) {
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

	if cmdDict, err := dict.NewCmdDict(f); err != nil {
		return nil, err
	} else {
		cmd.AddCommand(cmdDict)
	}

	if cmdNotebook, err := dict.NewCmdNotebook(f); err != nil {
		return nil, err
	} else {
		cmd.AddCommand(cmdNotebook)
	}

	if cmdServer, err := dict.NewCmdServer(f); err != nil {
		return nil, err
	} else {
		cmd.AddCommand(cmdServer)
	}

	return cmd, nil
}
