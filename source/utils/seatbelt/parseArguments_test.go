package seatbelt

import "testing"

func TestParseArguments(t *testing.T) {

	sandbox, read_paths, target, target_arguments, err0 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--read=/scripts",
		"--",
		"/usr/bin/cat",
		"file.txt",
	})

	if err0 == nil {

		if sandbox != "/sandbox" {
			t.Errorf("Expected %q, got %q", "/sandbox", sandbox)
		}

		if len(read_paths) != 1 || read_paths[0] != "/scripts" {
			t.Errorf("Expected read paths %v, got %v", []string{"/scripts"}, read_paths)
		}

		if target != "/usr/bin/cat" {
			t.Errorf("Expected %q, got %q", "/usr/bin/cat", target)
		}

		if len(target_arguments) != 1 || target_arguments[0] != "file.txt" {
			t.Errorf("Expected target arguments %v, got %v", []string{"file.txt"}, target_arguments)
		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

	_, _, _, _, err1 := parseArguments([]string{"--sandbox=/sandbox"})

	if err1 != nil {
		// expected: a target is required
	} else {
		t.Errorf("Expected invalid arguments to fail")
	}

}
