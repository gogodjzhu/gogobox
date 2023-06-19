package dict_etymonline

import (
	"gogobox/pkg/dict"
	"reflect"
	"testing"
)

func TestDictYoudao_Search(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name    string
		args    args
		want    *dict.WordInfo
		wantErr bool
	}{
		{
			name: "test-invalid",
			args: args{
				word: "inooo",
			},
			want:    dict.InvalidWord("inooo"),
			wantErr: false,
		},
		{
			name: "test-bing",
			args: args{
				word: "bing",
			},
			want: &dict.WordInfo{
				Word: "bing",
				Defines: []*dict.WordDefine{
					{
						Phonetics:  []string{},
						Definition: "++++bing (n.)\n\"heap or pile,\" 1510s, from Old Norse bingr \"heap.\" Also used from early 14c. as a word for bin, perhaps from notion of \"place where things are piled.\"",
					},
					{
						Phonetics:  []string{},
						Definition: "++++Bing (adj.)\nin reference to a a dark red type of cherry widely grown in the U.S., 1889, said to have been developed 1870s and named for Ah Bing, Chinese orchard foreman for Oregon fruit-grower Seth Lewelling.",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DictEtymonline{}
			got, err := d.Search(tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
