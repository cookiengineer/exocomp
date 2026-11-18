package capabilities

import "net"
import "os"
import "path/filepath"
import "testing"

func TestCapabilityNormalizeFiles(t *testing.T) {

	folder := t.TempDir()
	file := filepath.Join(folder, "file.txt")

	err0 := os.WriteFile(file, []byte("test"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create test file: %v", err0)
	}

	normalized, err1 := Normalize(Capability{Type: CapabilityFilesRead, Scope: file})
	resolved, _ := filepath.EvalSymlinks(file)

	if err1 == nil {
		if normalized.Scope != resolved {
			t.Errorf("Expected %q, got %q", resolved, normalized.Scope)
		}
	} else {
		t.Errorf("Expected %v to be nil", err1)
	}

	_, err2 := Normalize(Capability{Type: CapabilityFilesRead, Scope: "/"})

	if err2 == nil {
		t.Errorf("Expected / to be rejected")
	}

	_, err3 := Normalize(Capability{Type: CapabilityFilesWrite, Scope: "relative/path"})

	if err3 == nil {
		t.Errorf("Expected relative path to be rejected")
	}

}

func TestCapabilityNormalizeNetwork(t *testing.T) {

	normalized, err0 := Normalize(Capability{Type: CapabilityNetworkConnect, Scope: "Example.COM"})

	if err0 == nil {
		if normalized.Scope != "example.com" {
			t.Errorf("Expected %q, got %q", "example.com", normalized.Scope)
		}
	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

	normalized_ip, err1 := Normalize(Capability{Type: CapabilityNetworkConnect, Scope: "10.0.0.1"})

	if err1 == nil {
		if normalized_ip.Scope != "10.0.0.1" {
			t.Errorf("Expected %q, got %q", "10.0.0.1", normalized_ip.Scope)
		}
	} else {
		t.Errorf("Expected %v to be nil", err1)
	}

	normalized_cidr, err2 := Normalize(Capability{Type: CapabilityNetworkConnect, Scope: "10.0.0.1/8"})

	if err2 == nil {
		if normalized_cidr.Scope != "10.0.0.0/8" {
			t.Errorf("Expected %q, got %q", "10.0.0.0/8", normalized_cidr.Scope)
		}
	} else {
		t.Errorf("Expected %v to be nil", err2)
	}

	for _, scope := range []string{"example.com:443", "https://example.com", "*.example.com"} {

		_, err3 := Normalize(Capability{Type: CapabilityNetworkConnect, Scope: scope})

		if err3 == nil {
			t.Errorf("Expected scope %q to be rejected", scope)
		}

	}

}

func TestSet(t *testing.T) {

	folder := t.TempDir()
	file := filepath.Join(folder, "file.txt")

	err0 := os.WriteFile(file, []byte("test"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create test file: %v", err0)
	}

	set := NewSet()
	err1 := set.Grant([]Capability{
		{Type: CapabilityFilesRead, Scope: folder},
		{Type: CapabilityFilesWrite, Scope: folder},
		{Type: CapabilityNetworkConnect, Scope: "example.com"},
		{Type: CapabilityNetworkConnect, Scope: "10.0.0.0/8"},
	})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if set.CanRead(file) == false {
		t.Errorf("Expected read access")
	}

	if set.CanWrite(file) == false {
		t.Errorf("Expected write access")
	}

	if set.CanConnect("example.com") == false {
		t.Errorf("Expected example.com access")
	}

	if set.CanConnect("api.example.com") == true {
		t.Errorf("Expected exact domain matching only")
	}

	if set.CanConnect("10.1.2.3") == false {
		t.Errorf("Expected CIDR access")
	}

	if set.CanConnect("192.168.1.1") == true {
		t.Errorf("Expected unrelated IP to be rejected")
	}

	read_paths, write_paths := set.SandboxPaths()
	resolved_folder, _ := filepath.EvalSymlinks(folder)

	if len(read_paths) != 1 || read_paths[0] != resolved_folder {
		t.Errorf("Expected read path %q, got %v", resolved_folder, read_paths)
	}

	if len(write_paths) != 1 || write_paths[0] != resolved_folder {
		t.Errorf("Expected write path %q, got %v", resolved_folder, write_paths)
	}

	missing, err2 := set.Missing([]Capability{
		{Type: CapabilityFilesRead, Scope: file},
		{Type: CapabilityNetworkConnect, Scope: "other.example.com"},
	})

	if err2 == nil {
		if len(missing) != 1 || missing[0].Scope != "other.example.com" {
			t.Errorf("Expected one missing network capability, got %v", missing)
		}
	} else {
		t.Errorf("Expected %v to be nil", err2)
	}

	set.Clear()

	if len(set.List()) != 0 {
		t.Errorf("Expected capabilities to be cleared")
	}

}

func TestCapabilityPathHelpers(t *testing.T) {

	folder := t.TempDir()
	file := filepath.Join(folder, "file.txt")
	outside := filepath.Join(t.TempDir(), "outside.txt")

	if IsPathWithin(folder, file) == false {
		t.Errorf("Expected file to be inside folder")
	}

	if IsPathWithin(folder, outside) == true {
		t.Errorf("Expected outside file to be outside folder")
	}

	if IsPathDenied("/proc/self") == false {
		t.Errorf("Expected /proc to be denied")
	}

	if IsPathDenied("/home") == true {
		t.Errorf("Expected /home not to be denied")
	}

	canonical_folder, _ := filepath.EvalSymlinks(folder)

	if IsCanonicalPath(canonical_folder) == false {
		t.Errorf("Expected temp folder to be canonical")
	}

	if IsCanonicalPath(folder+string(filepath.Separator)+".") == true {
		t.Errorf("Expected non-clean path not to be canonical")
	}

	if IsValidDomain("example.com") == false {
		t.Errorf("Expected example.com to be a valid domain")
	}

	if IsValidDomain("*.example.com") == true {
		t.Errorf("Expected wildcard domain to be rejected")
	}

	if NetworkContains("10.0.0.0/8", net.ParseIP("10.1.2.3")) == false {
		t.Errorf("Expected IP to be inside CIDR")
	}

	if NetworkContains("10.0.0.0/8", net.ParseIP("192.168.1.1")) == true {
		t.Errorf("Expected unrelated IP to be outside CIDR")
	}

}
