package main

import (
	"testing"
)

func Test_readChannels(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "No Input",
			input: "",
			want:  []string{},
		},
		{
			name:  "Empty String as Input",
			input: ",,,",
			want:  []string{},
		},
		{
			name:  "Single Channel",
			input: "channel_1",
			want:  []string{"channel_1"},
		},
		{
			name:  "Multiple Channels",
			input: "channel_1,channel_2,channel_3",
			want:  []string{"channel_1", "channel_2", "channel_3"},
		},
		{
			name:  "Multiple Channels with Empty Entries",
			input: "channel_1,,channel_3,,",
			want:  []string{"channel_1", "channel_3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := readChannels(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("readChannels() = %v, want %v", got, tt.want)
				return
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("readChannels() = %v, want %v", got, tt.want)
					break
				}
			}
		})
	}
}
