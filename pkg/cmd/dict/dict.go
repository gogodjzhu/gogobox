package dict

import (
	"fmt"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
	"strings"
)

func NewCmdDict(f *cmdutil.Factory) (*cobra.Command, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}

	var currentChapter string
	var endpoint string
	cmd := &cobra.Command{
		Use:   "dict <word>",
		Short: "Look up the word in the dictionary",
		Long:  "Look up the word in the dictionary, you can specify the dictionary by option",
		RunE: func(cmd *cobra.Command, args []string) error {
			dictionary, err := dict.NewDict(cfg.Dict)
			if err != nil {
				return err
			}
			wordItem, err := dictionary.Search(strings.TrimSpace(strings.Join(args, " ")))
			if err != nil {
				return err
			}
			_, _ = fmt.Fprint(f.IOStreams.Out, wordItem.RenderString())
			// add to notebook if notebook is specified
			if currentChapter != "" {
				notebook, err := dict.OpenNotebook(cfg.Notebook)
				if err != nil {
					return err
				}
				if _, err := notebook.Mark(wordItem.Word, dict.Learning); err != nil {
					return err
				}
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if currentChapter != cfg.Notebook.CurrentChapter || endpoint != cfg.Dict.Endpoint {
				cfg.Notebook.CurrentChapter = currentChapter
				cfg.Dict.Endpoint = endpoint
				if err := cfg.Save(); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", cfg.Dict.Endpoint, "Specify the dictionary, youdao or etymonline")
	cmd.Flags().StringVarP(&currentChapter, "chapter", "c", cfg.Notebook.CurrentChapter, "Specify the chapter")
	return cmd, nil
}
