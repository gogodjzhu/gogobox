package dict_youdao

import (
	"gogobox/internal/config"
	"gogobox/pkg/dict/entity"
	"reflect"
	"testing"
)

func TestDictYoudao_Search(t *testing.T) {
	type fields struct {
		conf *config.DictYoudaoConfig
	}
	type args struct {
		word string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		postFunc func(got *entity.WordItem, args args) bool
	}{
		{
			name: "Test1",
			fields: fields{
				conf: &config.DictYoudaoConfig{},
			},
			args: args{
				word: "note",
			},
			wantErr: false,
			postFunc: func(got *entity.WordItem, args args) bool {
				if got == nil || got.ID == "" || len(got.WordMeanings) == 0 || len(got.WordPhonetics) == 0 {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DictYoudao{
				conf: tt.fields.conf,
			}
			got, err := d.Search(tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.postFunc != nil && !tt.postFunc(got, tt.args) {
				t.Errorf("Search() postFunc failed")
			}
		})
	}
}

func Test_formatWordMeanings(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []*entity.WordMeaning
	}{
		{
			name: "Test1",
			args: args{
				str: "n. meaning1 \nvt.  meaning2 \nvi. mean\ning3 n. \n",
			},
			want: []*entity.WordMeaning{
				{
					PartOfSpeech: "n.",
					Definitions:  "meaning1",
					Examples:     nil,
				},
				{
					PartOfSpeech: "vt.",
					Definitions:  "meaning2",
					Examples:     nil,
				},
				{
					PartOfSpeech: "vi.",
					Definitions:  "mean\ning3",
					Examples:     nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatWordMeanings(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatWordMeanings() = %v, want %v", got, tt.want)
			}
		})
	}
}
