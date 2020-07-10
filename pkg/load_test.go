package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExplodeLists(t *testing.T) {
	cases := []struct {
		name      string
		line      string
		wantLines []string
	}{
		{
			name:      "line with no lists is returned as is",
			line:      "pass in",
			wantLines: []string{"pass in"},
		},
		{
			name:      "line with list with two items is expanded to two lines",
			line:      "pass in on { em0 em1 }",
			wantLines: []string{"pass in on em0", "pass in on em1"},
		},
		{
			name: "line with two lists with two items is expanded to four lines",
			line: "pass in on { em0 em1 } proto { tcp udp }",
			wantLines: []string{
				"pass in on em0 proto tcp",
				"pass in on em0 proto udp",
				"pass in on em1 proto tcp",
				"pass in on em1 proto udp",
			},
		},
		{
			name: "line with two lists with two items is expanded to four lines",
			line: "pass in on { em0 em1 } proto { tcp udp }",
			wantLines: []string{
				"pass in on em0 proto tcp",
				"pass in on em0 proto udp",
				"pass in on em1 proto tcp",
				"pass in on em1 proto udp",
			},
		},
		{
			name: "line with three lists with two items is expanded to eight lines",
			line: "pass in on { em0 em1 } proto { tcp udp } from { 10.0.0.0/24 10.0.1.0/24 }",
			wantLines: []string{
				"pass in on em0 proto tcp from 10.0.0.0/24",
				"pass in on em0 proto tcp from 10.0.1.0/24",
				"pass in on em0 proto udp from 10.0.0.0/24",
				"pass in on em0 proto udp from 10.0.1.0/24",
				"pass in on em1 proto tcp from 10.0.0.0/24",
				"pass in on em1 proto tcp from 10.0.1.0/24",
				"pass in on em1 proto udp from 10.0.0.0/24",
				"pass in on em1 proto udp from 10.0.1.0/24",
			},
		},
		{
			name: "line with two lists with three items is expanded to nine lines",
			line: "pass in on { em0 em1 em2 } proto { tcp udp smtp }",
			wantLines: []string{
				"pass in on em0 proto tcp",
				"pass in on em0 proto udp",
				"pass in on em0 proto smtp",
				"pass in on em1 proto tcp",
				"pass in on em1 proto udp",
				"pass in on em1 proto smtp",
				"pass in on em2 proto tcp",
				"pass in on em2 proto udp",
				"pass in on em2 proto smtp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExplodeList(tc.line)

			assert.Equal(t, tc.wantLines, actual)
		})
	}
}
