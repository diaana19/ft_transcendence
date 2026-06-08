package services

import "testing"

func TestLevel(t *testing.T) {
	cases := []struct {
		total int64
		want  int
	}{
		{total: 0, want: 0},
		{total: -5, want: 0},
		{total: 1, want: 0},
		{total: 2, want: 1},
		{total: 3, want: 1},
		{total: 4, want: 2},
		{total: 7, want: 2},
		{total: 8, want: 3},
		{total: 131071, want: 16},
		{total: 131072, want: 17},
		{total: 134254, want: 17},
		{total: 262143, want: 17},
		{total: 262144, want: 18},
	}
	for _, c := range cases {
		if got := Level(c.total); got != c.want {
			t.Errorf("Level(%d) = %d, want %d", c.total, got, c.want)
		}
	}
}
