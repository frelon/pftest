package main

import (
	"flag"
	"os"

	"frelon.se/pftest/pkg"
	log "github.com/sirupsen/logrus"
)

type TestResult struct {
	Successful pkg.TestCases
	Failed     pkg.TestCases
}

func main() {
	var rulesPath = flag.String("f", "/etc/pf.conf", "path to rule file")
	var testsPath = flag.String("t", "/etc/pftest.conf", "path to test file")

	flag.Parse()

	rules, err := pkg.LoadRuleSetFile(*rulesPath)
	if err != nil {
		log.Errorf("failed to load rules: %v", err)
		return
	}

	log.Infof("loaded %v rules", len(rules))

	tests, err := pkg.LoadTestsFile(*testsPath)
	if err != nil {
		log.Errorf("failed to load tests: %v", err)
		return
	}

	log.Infof("loaded %v tests", len(tests))

	result, err := RunTests(rules, tests)
	if err != nil {
		log.Errorf("error running passing tests: %v", err)
		return
	}

	if len(result.Failed) > 0 {
		log.Errorf("failed, %v/%v tests failed", len(result.Failed), len(tests))
		os.Exit(1)
		return
	}

	log.Infof("successful, %v tests passed", len(result.Successful))
}

func RunTests(rules pkg.RuleSet, packets pkg.TestCases) (TestResult, error) {
	fail := pkg.TestCases{}
	success := pkg.TestCases{}

	for _, p := range packets {
		lastRule, _, err := rules.Evaluate(p.Packet)
		if err != nil {
			return TestResult{}, err
		}

		if lastRule == nil {
			log.Warningf("no matching rule: %v", p)
			fail = append(fail, p)
			continue
		}

		if lastRule.Action == p.ExpectedAction {
			success = append(success, p)
		} else {
			log.Errorf("test failed: %v", p)
			fail = append(fail, p)
		}
	}

	result := TestResult{
		Successful: success,
		Failed:     fail,
	}

	return result, nil
}
