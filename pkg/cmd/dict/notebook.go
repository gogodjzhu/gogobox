package dict

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gogobox/internal/config"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
)

type NotebookOptions struct {
	Op     string
	Config *config.DictConfig
}

func NewNotebookCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &NotebookOptions{
		Op: "op",
	}
	if cfg, err := f.Config(); err != nil {
		fmt.Println(f.IOStreams.Out, "[Err] read config failed")
		return nil
	} else {
		opts.Config = cfg.Dict
	}
	notebook, err := dict.NewFileNotebook(opts.Config.NotebookPath)

	if err != nil {
		fmt.Println(f.IOStreams.Out, "[Err] create notebook failed")
	}

	cmd := &cobra.Command{
		Use:   "notebook <word>",
		Short: "Learning words in notebook",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch opts.Op {
			case "learn":
				return fmt.Errorf("not implemented")
			case "review":
				word, err := notebook.Review()
				if err != nil {
					return err
				}
				red := color.New(color.FgRed).SprintFunc()
				fmt.Fprintln(f.IOStreams.Out, red(word.Word))
			default:
				return fmt.Errorf("unknown operation: %s", opts.Op)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.Op, "op", "o", "review", "Specify operation, learn or review")
	return cmd
}
