package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	type args struct {
		osArgs []string
		stdin  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Check", args{[]string{"", "check"}, "{}"}},
		{"In", args{[]string{"", "in"}, "{}"}},
		{"Out", args{[]string{"", "out"}, "{}"}},
	}
	for _, tt := range tests {
		var stdin *os.File

		os.Args = tt.args.osArgs
		os.Stdin, stdin, _ = os.Pipe()
		fmt.Fprint(stdin, tt.args.stdin)

		main()
	}
}
