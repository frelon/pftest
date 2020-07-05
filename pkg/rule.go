package pkg

type Action string

type Direction string

const (
	Block     = "block"
	Pass      = "pass"
	Match     = "match"
	Antispoof = "antispoof"

	All = "all"
	In  = "in"
	Out = "out"
	Any = "any"
	On  = "on"
)

var (
	BlockAll = Rule{
		Action:    Block,
		From:      Any,
		To:        Any,
		Interface: Any,
	}

	PassAll = Rule{
		Action:    Pass,
		From:      Any,
		To:        Any,
		Interface: Any,
	}
)

type Rule struct {
	Action        Action
	Direction     Direction
	Interface     string
	From          string
	To            string
	AddressFamily string
	Protocol      string
	Quick         bool
}

type RuleSet []Rule

// Evaluate takes in a packet and runs it through the RuleSet, returning the
// last matching rule, and an array of all matching rules.
func (r RuleSet) Evaluate(packet Packet) (Rule, []Rule, error) {
	matches := []Rule{}

	for _, rule := range r {
		if rule.Matches(packet) {
			matches = append(matches, rule)
		}
	}

	return matches[len(matches)-1], matches, nil
}

func (r Rule) Matches(packet Packet) bool {
	return true
}

type VariableSet map[string]string

