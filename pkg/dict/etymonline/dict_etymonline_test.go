package dict_etymonline

import (
	"gogobox/internal/config"
	"gogobox/pkg/dict/entity"
	"testing"
)

func TestDictEtymonline_Search(t *testing.T) {
	type fields struct {
		conf *config.DictEtymonineConfig
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
				conf: &config.DictEtymonineConfig{},
			},
			args: args{
				word: "note",
			},
			wantErr: false,
			postFunc: func(got *entity.WordItem, args args) bool {
				if got == nil || got.ID == "" || len(got.WordMeanings) == 0 /*|| len(got.WordPhonetics) == 0*/ {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DictEtymonline{
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
