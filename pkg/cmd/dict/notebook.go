package dict

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/cmdutil/tui/tui_list"
	"gogobox/pkg/dict"
	dict_youdao "gogobox/pkg/dict/youdao"
)

type NotebookOptions struct {
	Op           string
	NotebookPath string
}

func NewCmdNotebook(f *cmdutil.Factory) (*cobra.Command, error) {
	var opts NotebookOptions

	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}
	opts.NotebookPath = cfg.Dict.NotebookPath

	cmd := &cobra.Command{
		Use:   "notebook <word>",
		Short: "Learning words in notebook",
		RunE: func(cmd *cobra.Command, args []string) error {
			notebook, err := dict.NewFileNotebook(opts.NotebookPath)
			var model tea.Model
			switch opts.Op {
			case "review":
				notes, err := notebook.List()
				if err != nil {
					return err
				}
				options := make([]tui_list.Option, len(notes))
				for i, note := range notes {
					options[i] = tui_list.NewOption(note.Word,
						fmt.Sprintf("lookupTimes:%d", note.LookupTimes))
				}
				model = tui_list.NewModel("Words review", options, []tui_list.CallbackFunc{
					{
						Keys:            []string{"enter"},
						FullDescription: "look up selected word",
						Callback: func(option tui_list.Option) {
							dictionary := dict_youdao.NewDictYoudao()
							wordInfo, _ := dictionary.Search(option.Title())
							fmt.Fprintln(f.IOStreams.Out, wordInfo.RenderString())
							if len(wordInfo.Defines) > 0 {
								if err := notebook.Mark(wordInfo.Word, dict.Learning); err != nil {
									fmt.Fprintln(f.IOStreams.Out, "[Err] mark word failed")
								}
							}
						},
					},
					{
						Keys:             []string{"d"},
						ShortDescription: "delete",
						FullDescription:  "delete selected word",
						Callback: func(option tui_list.Option) {
							if err := notebook.Mark(option.Title(), dict.Delete); err != nil {
								fmt.Fprintln(f.IOStreams.Out, "[Err] mark word failed")
							} else {
								fmt.Fprintln(f.IOStreams.Out, "delete word["+option.Title()+"] success")
							}
						},
					},
				})
			case "exam":
				fmt.Errorf("not implemented")
			default:
				return fmt.Errorf("unknown operation: %s", opts.Op)
			}
			_, err = tea.NewProgram(model).Run()
			return err
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Dict.NotebookPath == opts.NotebookPath {
				return nil
			}
			cfg.Dict.NotebookPath = opts.NotebookPath
			return cfg.Save()
		},
	}

	cmd.Flags().StringVarP(&opts.Op, "operation", "o", "review", "Specify operation, exam or review")
	cmd.Flags().StringVarP(&opts.NotebookPath, "notebookPath", "n", opts.NotebookPath, "Specify notebook path")
	return cmd, nil
}
