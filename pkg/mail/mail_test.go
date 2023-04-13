package mail

import (
	"reflect"
	"testing"
)

func TestBuildSender(t *testing.T) {
	type test struct {
		input Party
		want  string
	}

	tests := []test{
		{
			input: Party{
				Name:  "Michael Knight",
				Email: "mn@thefoundation.local",
			},
			want: "Michael Knight<mn@thefoundation.local>",
		},
		{
			input: Party{
				Email: "mn@thefoundation.local",
			},
			want: "mn@thefoundation.local",
		},
	}

	for _, tc := range tests {
		got := buildSender(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestBuildReceivers(t *testing.T) {
	type test struct {
		input []Party
		want  []string
	}

	tests := []test{
		{
			input: []Party{
				{
					Name:  "Michael Knight",
					Email: "mn@thefoundation.local",
				},
				{
					Email: "dm@thefoundation.local",
				},
			},
			want: []string{
				"Michael Knight<mn@thefoundation.local>",
				"dm@thefoundation.local",
			},
		},
	}

	for _, tc := range tests {
		got := buildReceivers(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
