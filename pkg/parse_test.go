package pkg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseActionRule(t *testing.T) {
	cases := []struct {
		name       string
		wantErr    string
		rule       string
		wantRule   Rule
		wantTokens []string
	}{
		{
			name: "parses match rule",
			rule: "match",
			wantRule: Rule{
				Action: Match,
			},
			wantTokens: []string{},
		},
		{
			name: "parses block rule with direction",
			rule: "block in all",
			wantRule: Rule{
				Action: Block,
			},
			wantTokens: []string{In, All},
		},
		{
			name: "parses pass rule with direction",
			rule: "pass out all",
			wantRule: Rule{
				Action: Pass,
			},
			wantTokens: []string{Out, All},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := strings.Split(tc.rule, " ")
			rule := Rule{}
			tokensLeft, err := ParseAction(&rule, tokens)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantRule, rule)
			assert.Equal(t, tc.wantTokens, tokensLeft)
		})
	}
}

func Test_Take(t *testing.T) {
	cases := []struct {
		name     string
		tokens   []string
		takeFunc TakeFunc
		wantLeft []string
	}{
		{
			name:     "take everything",
			tokens:   []string{"match", "all", "test"},
			takeFunc: func(token string) bool { return false },
			wantLeft: []string{},
		},
		{
			name:     "take until open param",
			tokens:   []string{"match", "in", "all", "scrub", "(no-df", "random-id", "max-mss", "1440)"},
			takeFunc: func(token string) bool { return strings.Contains(token, "(") },
			wantLeft: []string{"random-id", "max-mss", "1440)"},
		},
		{
			name:     "take until close param",
			tokens:   []string{"match", "in", "all", "scrub", "(no-df", "random-id", "max-mss", "1440)"},
			takeFunc: func(token string) bool { return strings.Contains(token, ")") },
			wantLeft: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			actual := Take(tc.tokens, tc.takeFunc)

			assert.Equal(t, tc.wantLeft, actual)
		})
	}
}

func Test_ParseFromTo(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		wantRule Rule
		wantErr  string
	}{
		{
			name: "from any to 64:ff9b::/96",
			wantRule: Rule{
				From: Any,
				To:   "64:ff9b::/96",
			},
		},
		{
			name: "from self",
			wantRule: Rule{
				From: "self",
				To:   Any,
			},
		},
		{
			name: "to port 25",
			wantRule: Rule{
				From: Any,
				To:   "port 25",
			},
		},
		{
			name: "from any to 1.2.3.4 port > 123",
			wantRule: Rule{
				From: Any,
				To:   "1.2.3.4 port > 123",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := strings.Split(tc.name, " ")
			rule := Rule{}
			tokensLeft, err := ParseFromTo(&rule, tokens)

			require.Equal(t, []string{}, tokensLeft)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

