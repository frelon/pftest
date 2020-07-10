package main

import (
	"flag"

	"frelon.se/pftest/pkg"
	log "github.com/sirupsen/logrus"
)

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

	pass := pkg.Action(pkg.Pass)
	RunTests(pass, rules, tests[pass])
	if err != nil {
		log.Errorf("error running passing tests: %v", err)
		return
	}

	block := pkg.Action(pkg.Block)
	RunTests(block, rules, tests[block])
	if err != nil {
		log.Errorf("error running block tests: %v", err)
		return
	}

	log.Infof("done, exiting")
}

func RunTests(expected pkg.Action, rules pkg.RuleSet, packets []pkg.Packet) error {
	var (
		success int
		fail    int
	)

	for _, p := range packets {
		lastRule, _, err := rules.Evaluate(p)
		if err != nil {
			return err
		}

		if lastRule == nil {
			log.Warning("no matching rule")
			fail++
			continue
		}

		if lastRule.Action == expected {
			success++
		} else {
			fail++
		}
	}

	if fail == 0 {
		log.Infof("All tests successful")
	} else {
		log.Warningf("%v/%v tests failed", fail, fail+success)
	}

	return nil
}
