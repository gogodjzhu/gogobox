package dict

import (
	"github.com/fatih/color"
	"strings"
)

type Action string

const (
	Learning Action = "learning"
	Learned  Action = "learned"
	Delete   Action = "delete"
)

type Dict interface {
	Search(word string) (*WordInfo, error)
}

type Notebook interface {
	Mark(word string, action Action) error
	Get(word string) (*WordNote, error)
	List() ([]*WordNote, error)
}

type WordInfo struct {
	Word    string        `yaml:"word"`
	Defines []*WordDefine `yaml:"defines"`
}

func (f *WordInfo) RenderString() string {
	red := color.New(color.FgRed).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

	var str string
	str += red(f.Word) + "\n"
	if len(f.Defines) > 0 {
		for _, define := range f.Defines {
			str += green(strings.Join(define.Phonetics, " ")) + "\n"
			for _, s := range strings.Split(define.Definition, "\n") {
				switch {
				case strings.HasPrefix(s, "----"):
					str += gray(s[4:]) + "\n"
				case strings.HasPrefix(s, "++++"):
					str += cyan(s[4:]) + "\n"
				default:
					str += s + "\n"
				}
			}
		}
	}
	return str
}

type WordDefine struct {
	Phonetics  []string `yaml:"phonetics"`
	Definition string   `yaml:"definition"`
}

type WordNote struct {
	Word           string `yaml:"word"`
	LookupTimes    int    `yaml:"lookup_times"`
	LastLookupTime int64  `yaml:"last_lookup_time"`
}

func InvalidWord(word string) *WordInfo {
	return &WordInfo{
		Word: word,
		Defines: []*WordDefine{
			{
				Definition: "Word not found: " + word,
			},
		},
	}
}
