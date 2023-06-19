package util

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_sendGet(t *testing.T) {
	type args struct {
		url    string
		header map[string]string
		wrap   func(resp *http.Response) interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "baidu",
			args: args{
				url: "https://www.baidu.com",
				wrap: func(resp *http.Response) interface{} {
					return resp.StatusCode
				},
			},
			want:    200,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendGet(tt.args.url, tt.args.header, tt.args.wrap)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
