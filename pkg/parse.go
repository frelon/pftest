package pkg

import (
	"strings"

	"github.com/pkg/errors"
)

type ParseFunc func(rule *Rule, tokens []string) ([]string, error)

func ParseAction(rule *Rule, tokens []string) ([]string, error) {
	switch tokens[0] {
	case "match":
		rule.Action = Match
	case "block":
		rule.Action = Block
	case "pass":
		rule.Action = Pass
	default:
		return tokens, errors.Errorf("failed to parse %v", tokens[0])
	}

	return tokens[1:], nil
}

func ParseDirection(rule *Rule, tokens []string) ([]string, error) {
	switch tokens[0] {
	case "in":
		rule.Direction = In
	case "out":
		rule.Direction = Out
	default:
		rule.Direction = Any
		return tokens, nil
	}

	return tokens[1:], nil
}

func ParseLog(rule *Rule, tokens []string) ([]string, error) {
	switch tokens[0] {
	case "log":
		return ParseLogParams(rule, tokens[1:])
	default:
		return tokens, nil
	}
}

func ParseLogParams(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == "(all)" {
		return tokens[1:], nil
	}

	return tokens, nil
}

func ParseInterface(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == On {
		rule.Interface = tokens[1]
		return tokens[2:], nil
	}

	return tokens, nil
}

func ParseFromTo(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == All {
		rule.From = Any
		rule.To = Any
		return tokens[1:], nil
	}

	if tokens[0] == "from" {

	}

	if tokens[0] == "to" {

	}

	return tokens, nil
}

func ParseScrub(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == "scrub" {
		return Take(tokens[1:], func(t string) bool {
			return strings.Contains(t, ")")
		}), nil
	}

	return tokens, nil
}

func ParseNAT(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == "nat-to" {

	}

	return tokens, nil
}

func ParseAF(rule *Rule, tokens []string) ([]string, error) {
	switch tokens[0] {
	case "inet":
		rule.AddressFamily = "inet"
	case "inet6":
		rule.AddressFamily = "inet6"
	default:
		return tokens[:], nil
	}

	return tokens[1:], nil
}

func ParseProto(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == "proto" {
		rule.Protocol = tokens[1]
		return tokens[2:], nil
	}

	return tokens, nil
}

func ParseQuick(rule *Rule, tokens []string) ([]string, error) {
	if tokens[0] == "quick" {
		rule.Quick = true
		return tokens[1:], nil
	}

	return tokens, nil
}

func ParseRule(line string) (Rule, error) {
	parsers := []ParseFunc{
		ParseAction,
		ParseDirection,
		ParseLog,
		ParseQuick,
		ParseInterface,
		ParseAF,
		ParseProto,
		ParseFromTo,
		ParseScrub,
		ParseNAT,
	}

	tokens := Tokenize(line)

	r := Rule{}
	for _, p := range parsers {
		if len(tokens) == 0 {
			break
		}

		tokensLeft, err := p(&r, tokens)
		if err != nil {
			return r, errors.Wrap(err, "failed to parse")
		}

		tokens = tokensLeft
	}

	if len(tokens) != 0 {
		return r, errors.Errorf("unparsed tokens: %v", tokens)
	}

	return r, nil
}

func Tokenize(line string) []string {
	return strings.Split(line, " ")
}
