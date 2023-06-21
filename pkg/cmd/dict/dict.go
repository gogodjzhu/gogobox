package dict

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gogobox/internal/config"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
	dict_etymonline "gogobox/pkg/dict/etymonline"
	dict_youdao "gogobox/pkg/dict/youdao"
	"strings"
)

type Options struct {
	Endpoint string
	Config   *config.DictConfig
}

func NewCmdDict(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		Endpoint: "youdao",
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
		Use:   "dict <word>",
		Short: "Look up the word in the dictionary",
		Long:  "Look up the word in the dictionary, you can sepecify the dictionary by option",
		RunE: func(cmd *cobra.Command, args []string) error {
			var dictionary dict.Dict
			switch opts.Endpoint {
			case "youdao":
				dictionary = dict_youdao.NewDictYoudao()
			case "etymonline":
				dictionary = dict_etymonline.NewDictEtymonline()
			default:
				return fmt.Errorf("unknown dictionary: %s", opts.Endpoint)
			}

			wordInfo, err := dictionary.Search(strings.Join(args, " "))
			if err != nil {
				return err
			}
			red := color.New(color.FgRed).SprintFunc()
			gray := color.New(color.FgHiBlack).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			green := color.New(color.FgHiGreen).SprintFunc()
			fmt.Fprintln(f.IOStreams.Out, red(wordInfo.Word))
			if len(wordInfo.Defines) > 0 {
				for _, define := range wordInfo.Defines {
					fmt.Fprintln(f.IOStreams.Out, green(strings.Join(define.Phonetics, " ")))
					for _, s := range strings.Split(define.Definition, "\n") {
						switch {
						case strings.HasPrefix(s, "----"):
							fmt.Fprintln(f.IOStreams.Out, gray(s[4:]))
						case strings.HasPrefix(s, "++++"):
							fmt.Fprintln(f.IOStreams.Out, cyan(s[4:]))
						default:
							fmt.Fprintln(f.IOStreams.Out, s)
						}
					}
				}
				if err := notebook.Mark(wordInfo.Word, dict.Learning); err != nil {
					fmt.Fprintln(f.IOStreams.Out, "[Err] mark word failed")
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", "youdao", "Specify the dictionary, youdao or etymonline")
	return cmd
}
