package dict_etymonline

import (
	"github.com/PuerkitoBio/goquery"
	"gogobox/internal/util"
	"gogobox/pkg/dict"
	"net/http"
	"strings"
)

const Host = "https://www.etymonline.com"

type DictEtymonline struct {
}

func NewDictEtymonline() *DictEtymonline {
	return &DictEtymonline{}
}

func (d *DictEtymonline) Search(word string) (*dict.WordInfo, error) {
	url := Host + "/word/" + word
	result, err := util.SendGet(url, nil, func(response *http.Response) interface{} {
		if response.StatusCode != 200 {
			return dict.InvalidWord(word)
		}
		doc, e := goquery.NewDocumentFromReader(response.Body)
		if e != nil {
			return nil
		}
		defines := make([]*dict.WordDefine, 0)
		h := doc.Find("div.ant-col-xs-24").Children().Nodes
		for _, defineNode := range h {
			for _, attribute := range defineNode.Attr {
				if attribute.Key == "class" && strings.HasPrefix(attribute.Val, "word--") {
					subDoc := goquery.NewDocumentFromNode(defineNode)
					var defineWord string
					for _, node := range subDoc.Children().Children().Children().Nodes {
						for _, a := range node.Attr {
							if a.Key == "class" && strings.HasPrefix(a.Val, "word__name--") {
								defineWord = "++++" + goquery.NewDocumentFromNode(node).Text() // ++++ mark font color
								break
							}
						}
					}
					var defineParas string
					for _, definePara := range subDoc.Find("section").Children().Nodes {
						defineParaDoc := goquery.NewDocumentFromNode(definePara)
						if defineParaDoc.Is("blockquote") {
							defineParas += "----" // ---- mark font color
						}
						defineParas += strings.TrimSpace(defineParaDoc.Text()) + "\n"
					}
					defines = append(defines, &dict.WordDefine{
						Phonetics:  []string{},
						Definition: defineWord + "\n" + strings.TrimSpace(defineParas),
					})
					break
				}
			}
		}
		return &dict.WordInfo{
			Word:    word,
			Defines: defines,
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dict.WordInfo), nil
}

func trimFormat(s string, sep string) string {
	var tf string
	for _, str := range strings.Split(s, "\n") {
		tf += strings.TrimSpace(str) + sep
	}
	return tf
}
