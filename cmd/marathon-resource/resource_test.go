package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_parsePayload(t *testing.T) {
	type args struct {
		p    params
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Reads file with no replacements", args{params{AppJSON: "app.json"}, "fixtures"}, []byte{123, 10, 32, 32, 32, 32, 34, 102, 111, 111, 34, 58, 32, 34, 98, 97, 114, 34, 10, 125, 10}, false},
		{"Reads file with replacements", args{params{AppJSON: "app_template.json", Replacements: []metadata{{"foo", "bar"}}}, "fixtures"}, []byte{123, 10, 32, 32, 32, 32, 34, 102, 111, 111, 34, 58, 32, 34, 98, 97, 114, 34, 10, 125, 10}, false},
		{"Reads file with bad tmpl", args{params{AppJSON: "app_template_bad.json", Replacements: []metadata{{"foo", "bar"}}}, "fixtures"}, nil, true},
	}
	for _, tt := range tests {
		got, err := parsePayload(tt.args.p, tt.args.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. parsePayload() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			p, _ := ioutil.ReadAll(got)
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("%q. parsePayload() = %v, want %v", tt.name, string(p), string(tt.want))
			}
		}
	}
}
