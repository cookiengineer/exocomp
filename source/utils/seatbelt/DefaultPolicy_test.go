package seatbelt

import "path/filepath"
import "testing"

func TestDefaultPolicy(t *testing.T) {

	sandbox := "/tmp/exocomp-sandbox"
	target := "/usr/bin/cat"
	policy := DefaultPolicy(sandbox, target, []string{"/tmp/scripts"}, []string{"/tmp/output"})

	found_sandbox := false
	found_target := false
	found_extra := false
	found_write := false
	found_device := false
	found_deny := false

	for _, path := range policy.Deny {

		if path == filepath.Join(sandbox, ".exocomp") {
			found_deny = true
		}

	}

	for _, path := range policy.ReadWrite {

		if path == sandbox {
			found_sandbox = true
		} else if path == "/tmp/output" {
			found_write = true
		}

	}

	for _, path := range policy.ReadOnly {

		if path == filepath.Dir(target) {
			found_target = true
		} else if path == "/tmp/scripts" {
			found_extra = true
		}

	}

	for _, path := range policy.Devices {

		if path == "/dev/null" {
			found_device = true
		}

	}

	if found_sandbox == false {
		t.Errorf("Expected sandbox to be read-write")
	}

	if found_target == false {
		t.Errorf("Expected target folder to be read-only")
	}

	if found_extra == false {
		t.Errorf("Expected extra read path to be read-only")
	}

	if found_write == false {
		t.Errorf("Expected extra write path to be read-write")
	}

	if found_device == false {
		t.Errorf("Expected device node to be allowed")
	}

	if found_deny == false {
		t.Errorf("Expected .exocomp folder to be denied")
	}

	if policy.Network == true {
		t.Errorf("Expected network access to be denied by default")
	}

}
