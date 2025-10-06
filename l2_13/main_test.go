package main

import (
	"reflect"
	"testing"
)

func Test_flags_parseFParameter(t *testing.T) {
	type fields struct {
		f string
	}
	//type flags struct {
	//	f string
	//}
	tests := []struct {
		name    string
		action  string
		fields  fields
		want    []int
		wantErr bool
	}{
		{
			name:   "test1",
			action: "simple check",
			fields: fields{
				f: "1,3-5",
			},
			want:    []int{1, 3, 4, 5},
			wantErr: false,
		},
		{
			name:   "test2",
			action: "check collision",
			fields: fields{
				f: "1-4,3-5",
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name:   "test3",
			action: "sentence with only comma",
			fields: fields{
				f: "1,3,4,5,6,7",
			},
			want:    []int{1, 3, 4, 5, 6, 7},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &flags{
				f: tt.fields.f,
			}
			got, err := f.parseFParameter()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFParameter() got = %v, want %v", got, tt.want)
			} else {
				t.Log("OK!; parseFParameter() got = ", got)
			}
		})
	}
}
