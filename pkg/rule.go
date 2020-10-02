package pkg

import (
	"fmt"
	"net"
	"strings"
)

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

	BlockAllQuick = Rule{
		Action:    Block,
		From:      Any,
		To:        Any,
		Interface: Any,
		Quick:     true,
	}

	PassAll = Rule{
		Action:    Pass,
		From:      Any,
		To:        Any,
		Interface: Any,
	}
)

type VariableSet map[string]string

type Rule struct {
	Action        Action
	BlockPolicy   string
	Direction     Direction
	Interface     string
	From          string
	FromPort      string
	To            string
	ToPort        string
	AddressFamily string
	Protocol      string
	NAT           string
	RedirectTo    string
	Quick         bool
}

type RuleSet []Rule

// Evaluate takes in a packet and runs it through the RuleSet, returning the
// last matching rule, and an array of all matching rules.
func (r RuleSet) Evaluate(packet Packet) (*Rule, []Rule, error) {
	matches := []Rule{}

	for _, rule := range r {
		if rule.Matches(packet) {
			matches = append(matches, rule)

			if rule.Quick {
				break
			}
		}
	}

	if len(matches) < 1 {
		return nil, []Rule{}, nil
	}

	return &matches[len(matches)-1], matches, nil
}

func (r Rule) Matches(packet Packet) bool {
	if r.Interface != Any && r.Interface != packet.Interface {
		return false
	}

	if r.From != Any && !matchesAddress(r.From, packet.Source) {
		return false
	}

	if r.To != Any && !matchesAddress(r.To, packet.Destination) {
		return false
	}

	return true
}

func matchesAddress(rawNet, rawAddress string) bool {
	if strings.Contains(rawNet, "/") {
		address := net.ParseIP(rawAddress)
		_, net, err := net.ParseCIDR(rawNet)

		if err != nil {
			fmt.Printf("Err parsing net: %v\n", err)
			return false
		}

		return net.Contains(address)
	}

	return rawAddress == rawNet
}
