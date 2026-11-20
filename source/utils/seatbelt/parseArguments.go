package seatbelt

import "fmt"
import "strings"

func parseArguments(arguments []string) (string, []string, []string, bool, string, []string, error) {

	sandbox := ""
	target := ""
	network := "deny"
	read_paths := make([]string, 0)
	write_paths := make([]string, 0)
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
		} else if strings.HasPrefix(argument, "--write=") {
			write_paths = append(write_paths, strings.TrimSpace(argument[len("--write="):]))
		} else if strings.HasPrefix(argument, "--network=") {
			network = strings.TrimSpace(argument[len("--network="):])
		} else {
			return "", nil, nil, false, "", nil, fmt.Errorf("seatbelt.parseArguments: Invalid argument.")
		}

	}

	if network == "allow" || network == "deny" {

		if sandbox != "" && target != "" {
			return sandbox, read_paths, write_paths, network == "allow", target, target_arguments, nil
		} else {
			return "", nil, nil, false, "", nil, fmt.Errorf("seatbelt.parseArguments: Invalid argument.")
		}

	} else {
		return "", nil, nil, false, "", nil, fmt.Errorf("seatbelt.parseArguments: Invalid argument.")
	}
}
