package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadRuleset(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   string
		ruleset   string
		wantRules RuleSet
	}{
		{
			name:      "empty input returns empty ruleset",
			wantRules: RuleSet{},
		},
		{
			name:    "single rule returns 1 rule",
			ruleset: "block all",
			wantRules: RuleSet{
				{Action: Block, Direction: Any, From: Any, To: Any, Interface: Any},
			},
		},
		{
			name: "two rules returns 2 rules",
			ruleset: `block all
pass out all`,
			wantRules: RuleSet{
				{Action: Block, Direction: Any, From: Any, To: Any, Interface: Any},
				{Action: Pass, Direction: Out, From: Any, To: Any, Interface: Any},
			},
		},
		{
			name:    "junk input returns error",
			ruleset: "wjnckjwahenwajngvkwah",
			wantErr: "failed to parse",
		},
		{
			name:      "does not parse variable as a rule",
			ruleset:   `var = "test"`,
			wantRules: RuleSet{},
		},
		{
			name: "replaces variable inside other variable",
			ruleset: `var1 = test
var2 = $var1
pass on $var2 all`,
			wantRules: RuleSet{
				{
					Action:    Pass,
					Direction: Any,
					Interface: "test",
					From:      Any,
					To:        Any,
				},
			},
		},
		{
			name: "replaces variables in following rules",
			ruleset: `var = "test"
block on $var all`,
			wantRules: RuleSet{
				{Action: Block, Direction: Any, From: Any, To: Any, Interface: "test"},
			},
		},
		{
			name: "parses multiline rule",
			ruleset: `pass \
  out \
  all`,
			wantRules: RuleSet{
				{Action: Pass, Direction: Out, From: Any, To: Any, Interface: Any},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tc.ruleset)
			actual, err := LoadRuleSet(reader)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantRules, actual)
		})
	}
}

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
				PassAll,
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
			wantMatchingRules: []Rule{
				BlockAll,
				{Action: Pass, From: Any, To: Any, Interface: "em0"},
			},
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
			wantMatchingRules: []Rule{
				BlockAll,
			},
		},
		{
			name:   "Block all quick blocks all quick",
			packet: Packet{},
			ruleSet: RuleSet{
				BlockAllQuick,
				PassAll,
			},
			wantMatchingRules: []Rule{BlockAllQuick},
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
			assert.Equal(t, tc.wantMatchingRules[len(tc.wantMatchingRules)-1], *lastRule)
		})
	}
}
