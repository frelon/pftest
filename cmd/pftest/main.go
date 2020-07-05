package main

import (
	"flag"

	"frelon.se/pftest/pkg"
	log "github.com/sirupsen/logrus"
)

func main() {
	var rulesPath = flag.String("f", "/etc/pf.conf", "path to rule set")

	flag.Parse()

	rules, err := pkg.LoadRuleSetFile(*rulesPath)
	if err != nil {
		log.Error("failed loading rules", err)
		return
	}

	log.Infof("loaded %v rules", len(rules))
}
