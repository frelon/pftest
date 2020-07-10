package pkg

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type TestCases map[Action][]Packet

type Packet struct {
	Source      string
	Destination string
	Interface   string
}

func LoadTestsFile(path string) (TestCases, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return LoadTests(f)
}

func LoadTests(reader io.Reader) (TestCases, error) {
	tests := TestCases{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")

		packet := Packet{
			Source:      tokens[1],
			Destination: tokens[2],
			Interface:   tokens[3],
		}

		action := Action(tokens[0])

		tests[action] = append(tests[action], packet)
	}
	return tests, nil
}
