package dict

type Action string

const (
	Learning Action = "learning"
	Learned  Action = "learned"
)

type Dict interface {
	Search(word string) (*WordInfo, error)
}

type Notebook interface {
	Mark(word string, action Action) error
	Get(word string) (*WordNote, error)
}

type WordInfo struct {
	Word    string
	Defines []*WordDefine
}

type WordDefine struct {
	Phonetics  []string
	Definition string
}

type WordNote struct {
	Word           string `json:"word"`
	LookupTimes    int    `json:"lookup_times"`
	LastLookupTime int64  `json:"last_lookup_time"`
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
