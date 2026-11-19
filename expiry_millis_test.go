package redsync

import (
	"testing"
	"time"
)

func TestExpiryMillis(t *testing.T) {
	tests := []struct {
		in   time.Duration
		want int
	}{
		{0, 0},
		{-time.Second, 0},
		{500 * time.Microsecond, 1},
		{999 * time.Microsecond, 1},
		{time.Millisecond, 1},
		{1500 * time.Microsecond, 1},
		{2 * time.Millisecond, 2},
		{time.Second, 1000},
	}
	for _, tc := range tests {
		if got := expiryMillis(tc.in); got != tc.want {
			t.Fatalf("expiryMillis(%v)=%d, want %d", tc.in, got, tc.want)
		}
	}
}
