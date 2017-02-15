package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	if err := os.Setenv("GO_TESTING", "true"); err != nil {
		t.Fatal(err)
	}
	type args struct {
		osArgs []string
		stdin  string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		// {"Check", args{[]string{"", "check"}, "{}"}, false},
		// {"In", args{[]string{"", "in"}, "{}"}, false},
		// {"Out", args{[]string{"", "out"}, "{}"}, false},
		{"Bad json", args{[]string{"", "out"}, `{]`}, true},
		{"Wrong number of args", args{[]string{""}, "{}"}, true},
		{"Bad URI", args{[]string{"", "out"}, `{"source":{"uri":"http://192.168.0.%31/"}}`}, true},
		{"Unknown argument", args{[]string{"", "foo"}, `{}`}, true},
	}
	for _, tt := range tests {
		var stdin *os.File

		os.Args = tt.args.osArgs
		os.Stdin, stdin, _ = os.Pipe()
		fmt.Fprint(stdin, tt.args.stdin)

		assertPanic(t, main, tt.wantPanic)
	}
}

func assertPanic(t *testing.T, f func(), wantPanic bool) {
	defer func() {
		if r := recover(); (r != nil) != wantPanic {
			t.Errorf("Expected panic to be %t but was %t: %v", wantPanic, !wantPanic, r)
		}
	}()
	f()
}
