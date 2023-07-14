package dict

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
	dict_etymonline "gogobox/pkg/dict/etymonline"
	dict_youdao "gogobox/pkg/dict/youdao"
	"strings"
)

type Options struct {
	Endpoint     string
	NotebookPath string
}

func NewCmdDict(f *cmdutil.Factory) (*cobra.Command, error) {
	var opts Options
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}
	opts.NotebookPath = cfg.Dict.NotebookPath
	opts.Endpoint = cfg.Dict.Endpoint

	cmd := &cobra.Command{
		Use:   "dict <word>",
		Short: "Look up the word in the dictionary",
		Long:  "Look up the word in the dictionary, you can specify the dictionary by option",
		RunE: func(cmd *cobra.Command, args []string) error {
			var dictionary dict.Dict
			switch opts.Endpoint {
			case "youdao":
				dictionary = dict_youdao.NewDictYoudao()
			case "etymonline":
				dictionary = dict_etymonline.NewDictEtymonline()
			default:
				return errors.New("unknown dictionary:" + opts.Endpoint)
			}

			wordInfo, err := dictionary.Search(strings.Join(args, " "))
			if err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, wordInfo.RenderString())
			if len(wordInfo.Defines) > 0 {
				notebook, err := dict.NewFileNotebook(opts.NotebookPath)
				if err != nil {
					return err
				}
				if err := notebook.Mark(wordInfo.Word, dict.Learning); err != nil {
					return err
				}
			}
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Endpoint != cfg.Dict.Endpoint || opts.NotebookPath != cfg.Dict.NotebookPath {
				cfg.Dict.Endpoint = opts.Endpoint
				cfg.Dict.NotebookPath = opts.NotebookPath
				if err := cfg.Save(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", opts.Endpoint, "Specify the dictionary, youdao or etymonline")
	cmd.Flags().StringVarP(&opts.NotebookPath, "notebook", "n", opts.NotebookPath, "Specify the notebook path")
	return cmd, nil
}
