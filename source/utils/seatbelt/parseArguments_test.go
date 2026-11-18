package seatbelt

import "testing"

func TestParseArguments(t *testing.T) {

	sandbox, read_paths, write_paths, allow_network, target, target_arguments, err0 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--read=/scripts",
		"--write=/output",
		"--network=allow",
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

		if len(write_paths) != 1 || write_paths[0] != "/output" {
			t.Errorf("Expected write paths %v, got %v", []string{"/output"}, write_paths)
		}

		if allow_network == false {
			t.Errorf("Expected network access to be allowed")
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

	_, _, _, allow_default, _, _, err1 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--",
		"/usr/bin/cat",
	})

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if allow_default == true {
		t.Errorf("Expected network access to be denied by default")
	}

	_, _, _, allow_deny, _, _, err2 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--network=deny",
		"--",
		"/usr/bin/cat",
	})

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	if allow_deny == true {
		t.Errorf("Expected network access to be denied")
	}

	_, _, _, _, _, _, err3 := parseArguments([]string{"--sandbox=/sandbox"})

	if err3 != nil {
		// expected: a target is required
	} else {
		t.Errorf("Expected invalid arguments to fail")
	}

	_, _, _, _, _, _, err4 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--network=sometimes",
		"--",
		"/usr/bin/cat",
	})

	if err4 == nil {
		t.Errorf("Expected invalid network mode to fail")
	}

}
