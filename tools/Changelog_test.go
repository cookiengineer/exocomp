package tools

import "fmt"
import "os"
import "path/filepath"
import "strings"
import "testing"
import "time"

func TestChangelog_Add(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Add("./path/to/File.go", "Whatever", "New cache implementation")
		result2, err2 := tool.Add("./path/to/File.go", "Whatever", "New updated cache implementation")
		result3, err3 := tool.Add("./path/to/File.go", "Whatever", "New cache implementation")

		if strings.Contains(result1, "changelog.Add: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result2, "changelog.Add: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "changelog.Add: Log entry already exists for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to already exist")
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestChangelog_Change(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Change("./path/to/File.go", "Whatever", "New cache implementation")
		result2, err2 := tool.List()

		if strings.Contains(result1, "changelog.Change: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.HasPrefix(result2, "changelog.List: 1 changelog entries.") == false {
			t.Errorf("Expected 1 changelog entries")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestChangelog_Deprecate(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Deprecate("./path/to/File.go", "Whatever", "New cache implementation")
		result2, err2 := tool.List()

		if strings.Contains(result1, "changelog.Deprecate: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.HasPrefix(result2, "changelog.List: 1 changelog entries.") == false {
			t.Errorf("Expected 1 changelog entries")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestChangelog_Fix(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Fix("./path/to/File.go", "Whatever", "New cache implementation")
		result2, err2 := tool.List()

		if strings.Contains(result1, "changelog.Fix: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.HasPrefix(result2, "changelog.List: 1 changelog entries.") == false {
			t.Errorf("Expected 1 changelog entries")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestChangelog_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		// This should not be listed
		tool.Sandbox = playground
		tool.Add("./sub/Another.go", "Hidden", "Added secret functionality")
		tool.Sandbox = sandbox

		result1, err1 := tool.Add("./path/to/Cache.go", "Store", "Added new cache implementation")
		time.Sleep(2 * time.Second)
		result2, err2 := tool.Change("./path/to/Cache.go", "Store", "Changed new cache implementation")
		time.Sleep(2 * time.Second)
		result3, err3 := tool.Fix("./path/to/Cache.go", "Store", "Fixed new cache implementation")
		time.Sleep(2 * time.Second)
		result4, err4 := tool.Deprecate("./path/to/Cache.go", "Store", "Deprecated new cache implementation")
		time.Sleep(2 * time.Second)
		result5, err5 := tool.Remove("./path/to/Cache.go", "Store", "Removed new cache implementation")

		if strings.Contains(result1, "changelog.Add: Log entry created for ./path/to/Cache.go#Store at") == false {
			t.Errorf("Expected \"Add\" changelog entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result2, "changelog.Change: Log entry created for ./path/to/Cache.go#Store at") == false {
			t.Errorf("Expected \"Change\" changelog entry to be created")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "changelog.Fix: Log entry created for ./path/to/Cache.go#Store at") == false {
			t.Errorf("Expected \"Fix\" changelog entry to be created")
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if strings.Contains(result4, "changelog.Deprecate: Log entry created for ./path/to/Cache.go#Store at") == false {
			t.Errorf("Expected \"Deprecate\" changelog entry to be created")
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if strings.Contains(result5, "changelog.Remove: Log entry created for ./path/to/Cache.go#Store at") == false {
			t.Errorf("Expected \"Remove\" changelog entry to be created")
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		result6, err6 := tool.List()

		lines := strings.Split(result6, "\n")

		if len(lines) == 6 {

			if lines[0] != "changelog.List: 5 changelog entries." {
				t.Errorf("Expected %d changelog entries", 5)
			}

			if strings.Contains(lines[1], "Type: Add") == false {
				t.Errorf("Expected %s to be Type: Add", lines[1])
			}

			if strings.Contains(lines[2], "Type: Change") == false {
				t.Errorf("Expected %s to be Type: Change", lines[2])
			}

			if strings.Contains(lines[3], "Type: Fix") == false {
				t.Errorf("Expected %s to be Type: Fix", lines[3])
			}

			if strings.Contains(lines[4], "Type: Deprecate") == false {
				t.Errorf("Expected %s to be Type: Deprecate", lines[4])
			}

			if strings.Contains(lines[5], "Type: Remove") == false {
				t.Errorf("Expected %s to be Type: Remove", lines[5])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines), 6)
		}

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

}

func TestChangelog_Remove(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Remove("./path/to/File.go", "Whatever", "New cache implementation")
		result2, err2 := tool.List()

		if strings.Contains(result1, "changelog.Remove: Log entry created for ./path/to/File.go#Whatever at") == false {
			t.Errorf("Expected Log entry to be created")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.HasPrefix(result2, "changelog.List: 1 changelog entries.") == false {
			t.Errorf("Expected 1 changelog entries")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

	} else {
		t.Errorf("Expected %s to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestChangelog_Search(t *testing.T) {
	// TODO
	fmt.Println("TODO")
	t.Errorf("Changelog.Search unit test is not implemented yet")
}
