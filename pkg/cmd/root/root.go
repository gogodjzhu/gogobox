package root

import (
	"github.com/gogodjzhu/gogobox/pkg/cmd/minio"
	"github.com/gogodjzhu/gogobox/pkg/cmd/timefmt"
	"github.com/gogodjzhu/gogobox/pkg/cmd/version"
	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/spf13/cobra"
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

	cmd.AddCommand(version.NewCmdVersion(f))
	cmd.AddCommand(minio.NewCmdMinIO(f))
	cmd.AddCommand(timefmt.NewCmdTimeFmt(f))

	return cmd, nil
}
