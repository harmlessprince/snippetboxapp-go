package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2023 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2023 at 10:00",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.tm)
			if hd != test.want {
				t.Errorf("want %q; got %q", test.want, hd)
			}
		})
	}
}
