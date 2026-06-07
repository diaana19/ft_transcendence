package utils

import (
	"reflect"
	"testing"
)

func TestExtractHashtags(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"none", "just plain text", nil},
		{"basic", "love #golang and #redis", []string{"#golang", "#redis"}},
		{"lowercased and deduped", "#Go #go #GO done", []string{"#go"}},
		{"underscores and digits", "#web3 #snake_case ok", []string{"#web3", "#snake_case"}},
		{"stops at non-word", "#tag, end. #other!", []string{"#tag", "#other"}},
		{"order preserved", "#b #a #b", []string{"#b", "#a"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ExtractHashtags(c.in); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("ExtractHashtags(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}
