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

func NewCmdNotebook(f *cmdutil.Factory) (*cobra.Command, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}

	var op string
	var chapter string
	cmd := &cobra.Command{
		Use:   "chapter <word>",
		Short: "Learning words in chapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			notebook, err := dict.OpenNotebook(cfg.Notebook)
			if err != nil {
				return err
			}
			var model tea.Model
			switch op {
			case "review":
				notes, err := notebook.ListNotes()
				if err != nil {
					return err
				}
				initOptions := make([]tui_list.OptionEntity, len(notes))
				for i, note := range notes {
					initOptions[i] = tui_list.NewOption(&wordItemOptions{
						item:  note.WordItemId,
						title: note.Word,
						hint:  fmt.Sprintf("lookupTimes:%d", note.LookupTimes),
					})
				}
				model = tui_list.NewApp("Words review", initOptions, []tui_list.CallbackFunc{
					{
						Keys:            []string{"enter"},
						FullDescription: "look up selected word",
						Callback: func(selectedOption tui_list.OptionEntity) []tui_list.OptionEntity {
							dictionary, err := dict_youdao.NewDictYoudao(cfg.Dict.YoudaoConfig)
							if err != nil {
								_, _ = fmt.Fprintln(f.IOStreams.Out, "[Err] init dictionary failed")
								return nil
							}
							words, err := notebook.ListNotes()
							if err != nil {
								_, _ = fmt.Fprintln(f.IOStreams.Out, "[Err] list words failed")
								return nil
							}
							updateOptions := make([]tui_list.OptionEntity, len(words))
							for i, word := range words {
								if word.WordItemId != selectedOption.Entity().(string) {
									updateOptions[i] = tui_list.NewOption(&wordItemOptions{
										item:  word.WordItemId,
										title: word.Word,
										hint:  fmt.Sprintf("lookupTimes:%d", word.LookupTimes),
									})
								} else {
									wordItem, _ := dictionary.Search(word.Word)
									updateOptions[i] = tui_list.NewOption(&wordItemOptions{
										item:  word.WordItemId,
										title: wordItem.Word,
										hint:  wordItem.RawString(),
									})
								}
							}
							return updateOptions
						},
					},
					{
						Keys:             []string{"x"},
						ShortDescription: "delete",
						FullDescription:  "delete selected word",
						Callback: func(selectedOption tui_list.OptionEntity) []tui_list.OptionEntity {
							if _, err := notebook.Mark(selectedOption.Title(), dict.Delete); err != nil {
								fmt.Fprintln(f.IOStreams.Out, "[Err] mark word failed")
							}
							words, err := notebook.ListNotes()
							if err != nil {
								fmt.Fprintln(f.IOStreams.Out, "[Err] list words failed")
								return nil
							}
							updateOptions := make([]tui_list.OptionEntity, len(words))
							for i, word := range words {
								updateOptions[i] = tui_list.NewOption(&wordItemOptions{
									item:  word.WordItemId,
									title: word.Word,
									hint:  fmt.Sprintf("lookupTimes:%d", word.LookupTimes),
								})
							}
							return updateOptions
						},
					},
				})
			case "exam":
				_ = fmt.Errorf("not implemented")
			default:
				return fmt.Errorf("unknown operation: %s", op)
			}
			_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
			return err
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if chapter == cfg.Notebook.CurrentChapter {
				return nil
			}
			cfg.Notebook.CurrentChapter = chapter
			return cfg.Save()
		},
	}

	cmd.Flags().StringVarP(&op, "operation", "o", "review", "Specify operation, exam or review")
	cmd.Flags().StringVarP(&chapter, "chapter", "c", cfg.Notebook.CurrentChapter, "Specify chapter name")
	return cmd, nil
}

type wordItemOptions struct {
	item  string
	title string
	hint  string
}

func (w *wordItemOptions) Entity() interface{} {
	return w.item
}

func (w *wordItemOptions) Title() string {
	return w.title
}

func (w *wordItemOptions) Description() string {
	return w.hint
}
