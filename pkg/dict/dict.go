package dict

type Dict interface {
	Search(word string) (*WordInfo, error)
}

type WordInfo struct {
	Word    string
	Defines []*WordDefine
}

type WordDefine struct {
	Phonetics  []string
	Definition string
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
