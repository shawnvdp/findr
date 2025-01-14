package main

import (
	"reflect"
	"testing"
)

func TestScanFileForTerm(t *testing.T) {
	tests := []struct {
		contents []byte
		term     string
		want     []match
	}{
		{
			contents: []byte("this is the contents of a file"),
			term:     "file",
			want:     []match{{Line: "ents of a file", Number: 0}},
		},
		{
			contents: []byte("this is the contents of a file\nthis file also contains a second line"),
			term:     "file",
			want:     []match{{Line: "ents of a file", Number: 0}, {Line: "this file also cont", Number: 1}},
		},
	}
	for _, tt := range tests {
		matches := scanFileForTerm(tt.contents, tt.term)
		if !reflect.DeepEqual(tt.want, matches) {
			t.Fatalf("scanFileForTerm: want %+v, got %+v", tt.want, matches)
		}
	}
}
