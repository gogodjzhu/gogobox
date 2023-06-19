package dict_youdao

import (
	"github.com/PuerkitoBio/goquery"
	"gogobox/internal/util"
	"gogobox/pkg/dict"
	"net/http"
	"strings"
)

const Host = "https://dict.youdao.com"

type DictYoudao struct {
}

func NewDictYoudao() *DictYoudao {
	return &DictYoudao{}
}

func (d *DictYoudao) Search(word string) (*dict.WordInfo, error) {
	url := Host + "/search?q=" + word
	result, err := util.SendGet(url, nil, func(response *http.Response) interface{} {
		doc, e := goquery.NewDocumentFromReader(response.Body)
		if e != nil {
			return nil
		}
		keyword := strings.TrimSpace(doc.Find("span.keyword").Text())
		if len(strings.TrimSpace(keyword)) == 0 {
			return dict.InvalidWord(word)
		}
		trans := strings.TrimSpace(doc.Find("#phrsListTab > div.trans-container > ul").Text())
		enPhonetic := strings.TrimSpace(doc.Find("#phrsListTab > h2 > div > span:nth-child(1)").Text())
		usPhonetic := strings.TrimSpace(doc.Find("#phrsListTab > h2 > div > span:nth-child(2)").Text())
		return &dict.WordInfo{
			Word: word,
			Defines: []*dict.WordDefine{
				{
					Phonetics: []string{
						trimFormat(enPhonetic, " "),
						trimFormat(usPhonetic, " "),
					},
					Definition: trimFormat(trans, "\n"),
				},
			},
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
	return strings.TrimSpace(tf)
}
