package pkg

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	listRegex = regexp.MustCompile(`\{.+?\}`)
)

func LoadRuleSetFile(filename string) (RuleSet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %v", filename)
	}

	return LoadRuleSet(f)
}

func LoadRuleSet(reader io.Reader) (RuleSet, error) {
	rules := RuleSet{}
	vars := VariableSet{}

	multiline := ""
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := CleanString(scanner.Text())

		if strings.HasSuffix(line, `\`) {
			if multiline == "" {
				multiline = CleanString(strings.TrimSuffix(line, `\`))
				continue
			}

			multiline = multiline + " " + CleanString(strings.TrimSuffix(line, `\`))
			continue
		}

		if multiline != "" {
			line = multiline + " " + line
			multiline = ""
		}

		if !IsRuleLine(line) {
			continue
		}

		replaced := ReplaceVars(vars, line)
		if IsVariableDeclaration(replaced) {
			vars = AddVar(vars, replaced)
			continue
		}

		lines := ExplodeList(replaced)
		for _, l := range lines {
			rule, err := ParseRule(l)
			if err != nil {
				return rules, err
			}

			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func CleanString(line string) string {
	return strings.ToLower(strings.TrimSpace(line))
}

func ReplaceVars(vars VariableSet, line string) string {
	l := line

	for s, v := range vars {
		l = strings.ReplaceAll(l, "$"+s, strings.Trim(v, `"`))
	}

	return l
}

func IsRuleLine(line string) bool {
	if len(line) == 0 {
		return false
	}

	return !strings.HasPrefix(line, "#") &&
		!strings.HasPrefix(line, "set") &&
		!strings.HasPrefix(line, "table") &&
		!strings.HasPrefix(line, "antispoof")
}

func IsVariableDeclaration(line string) bool {
	return strings.Contains(line, "=")
}

func AddVar(vars VariableSet, line string) VariableSet {
	tokens := strings.Split(line, "=")
	name := strings.TrimSpace(tokens[0])
	value := strings.TrimSpace(tokens[1])
	vars[name] = value
	return vars
}

func ExplodeList(line string) []string {
	lines := []string{}
	m := listRegex.FindString(line)
	if m == "" {
		return []string{line}
	}

	cleaned := strings.ReplaceAll(m, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "{", "")
	cleaned = strings.ReplaceAll(cleaned, "}", "")
	cleaned = strings.TrimSpace(cleaned)
	split := strings.Split(cleaned, " ")

	for _, t := range split {
		replaced := strings.Replace(line, m, t, 1)
		lines = append(lines, replaced)
	}

	ret := []string{}
	for _, l := range lines {
		exploded := ExplodeList(l)
		ret = append(ret, exploded...)
	}

	return ret
}
