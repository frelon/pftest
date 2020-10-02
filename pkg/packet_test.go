package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Match(t *testing.T) {
	cases := []struct {
		name      string
		rule      Rule
		packet    Packet
		wantMatch bool
	}{
		{
			name:      "same interface matches",
			wantMatch: true,
			packet:    Packet{Interface: "em0"},
			rule:      Rule{Interface: "em0"},
		},
		{
			name:      "different interface does not match",
			wantMatch: false,
			packet:    Packet{Interface: "em0"},
			rule:      Rule{Interface: "em1"},
		},
		{
			name:      "matches from ipv4-address",
			wantMatch: true,
			packet:    Packet{Source: "192.168.0.2"},
			rule:      Rule{Interface: Any, From: "192.168.0.2"},
		},
		{
			name:      "does not match different from ipv4-address",
			wantMatch: false,
			packet:    Packet{Interface: "em0", Source: "192.168.0.3"},
			rule:      Rule{Interface: "em0", From: "192.168.0.2"},
		},
		{
			name:      "matches ipv4 net",
			wantMatch: true,
			packet:    Packet{Source: "192.168.0.3"},
			rule:      Rule{Interface: Any, From: "192.168.0.0/24"},
		},
		{
			name:      "does not match ipv4 address outside of net",
			wantMatch: false,
			packet:    Packet{Source: "172.16.0.3"},
			rule:      Rule{Interface: Any, From: "192.168.0.0/24"},
		},
		{
			name:      "matches ipv4 net",
			wantMatch: true,
			packet:    Packet{Source: "192.168.0.3"},
			rule:      Rule{Interface: Any, From: "192.168.0.0/24"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matches := tc.rule.Matches(tc.packet)
			assert.Equal(t, tc.wantMatch, matches)
		})
	}
}
