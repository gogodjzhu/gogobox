package dict_youdao

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
			name: "test-google",
			args: args{
				word: "google",
			},
			want: &dict.WordInfo{
				Word: "google",
				Defines: []*dict.WordDefine{
					{
						Phonetics: []string{
							"英 [ˈɡuːɡl]",
							"美 [ˈɡuːɡl]",
						},
						Definition: "vt. 用搜索引擎搜索（尤指用谷歌搜索引擎）\nn. （Google）谷歌（网络搜索引擎）",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DictYoudao{}
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
