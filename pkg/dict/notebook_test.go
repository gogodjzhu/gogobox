package dict

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func testInit(t *testing.T) {
	err := ioutil.WriteFile("/tmp/test_notebook.yml", []byte(
		`- word: hello
  lookup_times: 2
  last_lookup_time: 1234567890
- word: test
  lookup_times: 1
  last_lookup_time: 1234567890
`), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileNotebook_readNote(t *testing.T) {
	testInit(t)

	type fields struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*WordNote
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				filename: "/tmp/test_notebook.yml",
			},
			want: []*WordNote{
				{
					Word:           "hello",
					LookupTimes:    2,
					LastLookupTime: 1234567890,
				},
				{
					Word:           "test",
					LookupTimes:    1,
					LastLookupTime: 1234567890,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileNotebook{
				filename: tt.fields.filename,
			}
			got, err := f.readNote()
			if (err != nil) != tt.wantErr {
				t.Errorf("readNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readNote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileNotebook_writeNote(t *testing.T) {
	testInit(t)

	type fields struct {
		filename string
	}
	type args struct {
		notes []*WordNote
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				filename: "/tmp/test_notebook.yml",
			},
			args: args{
				notes: []*WordNote{
					{
						Word:           "test",
						LookupTimes:    2,
						LastLookupTime: 1234567890,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileNotebook{
				filename: tt.fields.filename,
			}
			if err := f.writeNote(tt.args.notes); (err != nil) != tt.wantErr {
				t.Errorf("writeNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
