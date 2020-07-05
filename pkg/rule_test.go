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
				{Action: Block, Direction: Any, From: Any, To: Any},
			},
		},
		{
			name: "two rules returns 2 rules",
			ruleset: `block all
pass out all`,
			wantRules: RuleSet{
				{Action: Block, Direction: Any, From: Any, To: Any},
				{Action: Pass, Direction: Out, From: Any, To: Any},
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
				{Action: Pass, Direction: Out, From: Any, To: Any},
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

