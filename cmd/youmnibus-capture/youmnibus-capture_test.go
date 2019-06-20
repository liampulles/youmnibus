package main

import "os"
import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("YOUTUBE_API_KEY", "AIzaSyCCCxwvYLd23yqbrF-61sscgueIXpAlyJ4")
			main()
		})
	}
}
