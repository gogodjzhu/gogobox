package util

import (
	"reflect"
	"testing"
)

func Test_splitWorker(t *testing.T) {
	type args struct {
		str            string
		separatorChars []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test-split-simple",
			args: args{
				str:            "hello world!hello world!",
				separatorChars: []string{" ", "!"},
			},
			want: []string{"hello", " ", "world", "!", "hello", " ", "world", "!"},
		},
		{
			name: "test-split-conjunction",
			args: args{
				str:            "hello###world",
				separatorChars: []string{"#"},
			},
			want: []string{"hello", "#", "#", "#", "world"},
		},
		{
			name: "test-split-phrase",
			args: args{
				str:            "hellotheworld",
				separatorChars: []string{"the"},
			},
			want: []string{"hello", "the", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitWorker(tt.args.str, tt.args.separatorChars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}
