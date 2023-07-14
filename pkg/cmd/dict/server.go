package dict

import (
	"fmt"
	"github.com/spf13/cobra"
	"gogobox/pkg/cmdutil"
	"gogobox/pkg/dict"
	dict_etymonline "gogobox/pkg/dict/etymonline"
	dict_youdao "gogobox/pkg/dict/youdao"
	"net/http"
)

func NewCmdServer(f *cmdutil.Factory) (*cobra.Command, error) {
	var port int
	var root string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the dict server",
		Long:  "Start the dict server, you can specify the port by option",
		RunE: func(cmd *cobra.Command, args []string) error {
			http.HandleFunc("/"+root, func(w http.ResponseWriter, r *http.Request) {
				endpoint := r.URL.Query().Get("endpoint")
				word := r.URL.Query().Get("word")
				var dictionary dict.Dict
				switch endpoint {
				case "youdao":
					dictionary = dict_youdao.NewDictYoudao()
				case "etymonline":
					dictionary = dict_etymonline.NewDictEtymonline()
				default:
					http.Error(w, "unknown dictionary:"+endpoint, http.StatusBadRequest)
				}

				wordInfo, err := dictionary.Search(word)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Write([]byte(wordInfo.RawString()))
			})
			fmt.Fprintln(f.IOStreams.Out, fmt.Sprintf("Server started at 0.0.0.0:%d/%s", port, root))
			return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "The port of the server")
	cmd.Flags().StringVarP(&root, "root", "r", "dict", "The root path of the server")
	return cmd, nil
}
