package dict_mwebster

import (
	"gogobox/internal/config"
	"gogobox/pkg/dict/entity"
	"reflect"
	"testing"
)

func TestDictMWebster_Search(t *testing.T) {
	type fields struct {
		conf *config.DictMWebsterConfig
	}
	type args struct {
		word string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.WordItem
		wantErr bool
	}{
		{
			name: "Test1",
			fields: fields{
				conf: &config.DictMWebsterConfig{
					Key: "8ed8d1c8-9f76-441c-bb5c-7dca25525212",
				},
			},
			args: args{
				word: "bottle",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DictMWebster{
				conf: tt.fields.conf,
			}
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
