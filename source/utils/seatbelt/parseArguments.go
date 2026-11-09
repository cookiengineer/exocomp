package seatbelt

import "fmt"
import "strings"

func parseArguments(arguments []string) (string, []string, string, []string, error) {

	sandbox := ""
	target := ""
	read_paths := make([]string, 0)
	target_arguments := make([]string, 0)
	after_separator := false

	for _, argument := range arguments {

		if after_separator == true {

			if target == "" {
				target = argument
			} else {
				target_arguments = append(target_arguments, argument)
			}

		} else if argument == "--" {
			after_separator = true
		} else if strings.HasPrefix(argument, "--sandbox=") {
			sandbox = strings.TrimSpace(argument[len("--sandbox="):])
		} else if strings.HasPrefix(argument, "--read=") {
			read_paths = append(read_paths, strings.TrimSpace(argument[len("--read="):]))
		} else {
			return "", nil, "", nil, fmt.Errorf("seatbelt.parseArguments: Invalid argument.")
		}

	}

	if sandbox != "" && target != "" {
		return sandbox, read_paths, target, target_arguments, nil
	} else {
		return "", nil, "", nil, fmt.Errorf("seatbelt.parseArguments: Invalid argument.")
	}

}
