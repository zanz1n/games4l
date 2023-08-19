package cmd

import (
	"fmt"
	"strings"
)

func ParseArgs(args []string) (map[string]string, []string, error) {
	parsed := make(map[string]string)
	cmd := []string{}

	// Loop state
	var (
		initialState = true
		skip         = false
		ll           = len(args)
	)

	for pos, arg := range args {
		if !skip && strings.HasPrefix(arg, "--") {
			initialState = false

			statement := arg[2:]
			if statement == "" {
				return nil, nil, fmt.Errorf("empty arg provided")
			}

			// lookup
			if ll-1 > pos {
				lookup := args[pos+1]
				if !strings.HasPrefix(lookup, "--") {
					skip = true
					parsed[statement] = args[pos+1]
				} else {
					parsed[statement] = "true"
				}
			} else {
				parsed[statement] = "true"
			}

			continue
		}

		if skip {
			skip = false
		}

		if initialState {
			cmd = append(cmd, arg)
		}
	}

	return parsed, cmd, nil
}
