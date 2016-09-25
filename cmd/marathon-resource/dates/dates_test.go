package dates

import (
	"reflect"
	"testing"
)

func Test_maxDateString(t *testing.T) {
	type args struct {
		stringSlice []string
		needle      string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Max string",
			args{[]string{"2015-04-11T09:31:50.021Z", "2014-03-01T23:42:20.938Z", "2015-02-11T09:31:50.021Z"}, "2015-02-11T09:31:50.021Z"},
			[]string{"2015-02-11T09:31:50.021Z", "2015-04-11T09:31:50.021Z"},
			false,
		},
		{
			"current version doesn't exist",
			args{[]string{"2015-04-11T09:31:50.021Z", "2014-03-01T23:42:20.938Z"}, "2015-02-11T09:31:50.021Z"},
			[]string{"2015-04-11T09:31:50.021Z"},
			false,
		},
		{
			"Bad string",
			args{[]string{"2015-02-11T09:31:50.021Z", "hello"}, "2015-02-11T09:31:50.021Z"},
			nil,
			true,
		},
		{
			"Bad string2",
			args{[]string{"2015-02-11T09:31:50.021Z"}, "bar"},
			nil,
			true,
		},
		{
			"Out of range",
			args{[]string{"2015-02-11T09:31:50.021Z"}, "2015-04-11T09:31:50.021Z"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		got, err := NewerTimestamps(tt.args.stringSlice, tt.args.needle)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewerTimestamps(%v) error = %v, wantErr %v", tt.name, tt.args.stringSlice, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewerTimestamps(%v) = %v, want %v", tt.name, tt.args.stringSlice, got, tt.want)
		}
	}
}
