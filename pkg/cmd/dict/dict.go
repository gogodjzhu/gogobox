package dict

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
	dict_etymonline "gogobox/pkg/dict/etymonline"
	dict_youdao "gogobox/pkg/dict/youdao"
	"strings"
)

type Options struct {
	Endpoint string
}

func NewCmdDict(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		Endpoint: "youdao",
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
			}

			wordInfo, err := dictionary.Search(strings.Join(args, " "))
			if err != nil {
				return err
			}
			red := color.New(color.FgRed).SprintFunc()
			gray := color.New(color.FgHiBlack).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			fmt.Fprintln(f.IOStreams.Out, red(wordInfo.Word))
			for _, define := range wordInfo.Defines {
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
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", "youdao", "Specify the dictionary, youdao or etymonline")
	return cmd
}
