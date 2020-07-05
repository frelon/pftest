package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Evaluate(t *testing.T) {
	cases := []struct {
		name              string
		ruleSet           RuleSet
		packet            Packet
		wantErr           string
		wantMatchingRules []Rule
	}{
		{
			name: "block all blocks all",
			ruleSet: RuleSet{
				{
					Action:    Block,
					From:      Any,
					To:        Any,
					Interface: Any,
				},
			},
			packet: Packet{},
			wantMatchingRules: []Rule{
				{Action: Block, From: Any, To: Any, Interface: Any},
			},
		},
		{
			name: "pass all passes all",
			ruleSet: RuleSet{
				BlockAll,
			},
			packet: Packet{},
			wantMatchingRules: []Rule{
				{Action: Pass, From: Any, To: Any, Interface: Any},
			},
		},
		{
			name: "passes in on correct interface and matches all rules",
			ruleSet: RuleSet{
				BlockAll,
				{Action: Pass, From: Any, To: Any, Interface: "em0"},
			},
			packet: Packet{
				Source:      "10.0.0.1",
				Destination: "10.0.0.2",
				Interface:   "em0",
			},
			wantMatchingRules: []Rule{},
		},
		{
			name: "blocks on incorrect interface",
			ruleSet: RuleSet{
				BlockAll,
				{Action: Pass, From: Any, To: Any, Interface: "em0"},
			},
			packet: Packet{
				Source:      "10.0.0.1",
				Destination: "10.0.0.2",
				Interface:   "em1",
			},
			wantMatchingRules: []Rule{},
		},
		{
			name: "Block all quick blocks all quick",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lastRule, matchingRules, err := tc.ruleSet.Evaluate(tc.packet)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantMatchingRules, matchingRules)
			assert.Equal(t, tc.wantMatchingRules[len(tc.wantMatchingRules)-1], lastRule)
		})
	}
}
