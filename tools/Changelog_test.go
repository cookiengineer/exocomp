package tools

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
		t.Errorf("Expected tool to be not nil")
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
		t.Errorf("Expected tool to be not nil")
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
		t.Errorf("Expected tool to be not nil")
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
		t.Errorf("Expected tool to be not nil")
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

		_, err1 := tool.Add("./path/to/Cache.go", "Store", "Added new cache implementation")
		time.Sleep(2 * time.Second)
		_, err2 := tool.Change("./path/to/Cache.go", "Store", "Changed new cache implementation")
		time.Sleep(2 * time.Second)
		_, err3 := tool.Fix("./path/to/Cache.go", "Store", "Fixed new cache implementation")
		time.Sleep(2 * time.Second)
		_, err4 := tool.Deprecate("./path/to/Cache.go", "Store", "Deprecated new cache implementation")
		time.Sleep(2 * time.Second)
		_, err5 := tool.Remove("./path/to/Cache.go", "Store", "Removed new cache implementation")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		result, err6 := tool.List()
		lines := strings.Split(result, "\n")

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
		t.Errorf("Expected tool to be not nil")
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
		t.Errorf("Expected tool to be not nil")
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

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {

		// This should not be listed
		tool.Sandbox = playground
		tool.Add("./sub/Another.go", "Hidden", "Added secret functionality")
		tool.Sandbox = sandbox

		_, err1 := tool.Add("./path/to/Cache.go", "Store", "Added new cache implementation")
		time.Sleep(2 * time.Second)
		_, err2 := tool.Add("./path/to/Cache.go", "StoreItem", "Added new cache implementation")
		time.Sleep(2 * time.Second)
		_, err3 := tool.Change("./path/to/Cache.go", "Store", "Changed new cache implementation")
		time.Sleep(2 * time.Second)
		_, err4 := tool.Fix("./path/to/Cache.go", "StoreItem", "Fixed new cache implementation")
		time.Sleep(2 * time.Second)
		_, err5 := tool.Deprecate("./path/to/Cache.go", "Store", "Deprecated new cache implementation")
		time.Sleep(2 * time.Second)
		_, err6 := tool.Remove("./path/to/Cache.go", "Store", "Removed new cache implementation")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

		result1, err7 := tool.Search("./path/to/Cache.go", "")
		result2, err8 := tool.Search("./path/to/Cache.go", "Store")
		result3, err9 := tool.Search("./path/to/Cache.go", "StoreItem")

		lines1 := strings.Split(result1, "\n")
		lines2 := strings.Split(result2, "\n")
		lines3 := strings.Split(result3, "\n")

		if len(lines1) == 7 {

			if lines1[0] != "changelog.Search: ./path/to/Cache.go contains 6 changelog entries." {
				t.Errorf("Expected %d changelog entries", 6)
			}

			if strings.Contains(lines1[1], "Type: Add,") == false {
				t.Errorf("Expected %s to be Type: Add", lines1[1])
			}

			if strings.Contains(lines1[1], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines1[1])
			}

			if strings.Contains(lines1[2], "Type: Change,") == false {
				t.Errorf("Expected %s to be Type: Change", lines1[2])
			}

			if strings.Contains(lines1[2], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines1[2])
			}

			if strings.Contains(lines1[3], "Type: Deprecate,") == false {
				t.Errorf("Expected %s to be Type: Deprecate", lines1[3])
			}

			if strings.Contains(lines1[3], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines1[3])
			}

			if strings.Contains(lines1[4], "Type: Remove,") == false {
				t.Errorf("Expected %s to be Type: Remove", lines1[4])
			}

			if strings.Contains(lines1[4], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines1[4])
			}

			if strings.Contains(lines1[5], "Type: Add,") == false {
				t.Errorf("Expected %s to be Type: Add", lines1[4])
			}

			if strings.Contains(lines1[5], "Symbol: StoreItem,") == false {
				t.Errorf("Expected %s to be Symbol: StoreItem", lines1[5])
			}

			if strings.Contains(lines1[6], "Type: Fix,") == false {
				t.Errorf("Expected %s to be Type: Fix", lines1[6])
			}

			if strings.Contains(lines1[6], "Symbol: StoreItem,") == false {
				t.Errorf("Expected %s to be Symbol: StoreItem", lines1[6])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines1), 7)
		}

		if err7 != nil {
			t.Errorf("Expected %v to be nil", err7)
		}

		if len(lines2) == 5 {

			if lines2[0] != "changelog.Search: ./path/to/Cache.go#Store contains 4 changelog entries." {
				t.Errorf("Expected %d changelog entries", 4)
			}

			if strings.Contains(lines2[1], "Type: Add,") == false {
				t.Errorf("Expected %s to be Type: Add", lines2[1])
			}

			if strings.Contains(lines2[1], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines2[1])
			}

			if strings.Contains(lines2[2], "Type: Change,") == false {
				t.Errorf("Expected %s to be Type: Change", lines2[2])
			}

			if strings.Contains(lines2[2], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines2[2])
			}

			if strings.Contains(lines2[3], "Type: Deprecate,") == false {
				t.Errorf("Expected %s to be Type: Deprecate", lines2[3])
			}

			if strings.Contains(lines2[3], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines2[3])
			}

			if strings.Contains(lines2[4], "Type: Remove,") == false {
				t.Errorf("Expected %s to be Type: Remove", lines2[4])
			}

			if strings.Contains(lines2[4], "Symbol: Store,") == false {
				t.Errorf("Expected %s to be Symbol: Store", lines2[4])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines2), 5)
		}

		if err8 != nil {
			t.Errorf("Expected %v to be nil", err8)
		}

		if len(lines3) == 3 {

			if lines3[0] != "changelog.Search: ./path/to/Cache.go#StoreItem contains 2 changelog entries." {
				t.Errorf("Expected %d changelog entries", 2)
			}

			if strings.Contains(lines3[1], "Type: Add,") == false {
				t.Errorf("Expected %s to be Type: Add", lines1[1])
			}

			if strings.Contains(lines3[1], "Symbol: StoreItem,") == false {
				t.Errorf("Expected %s to be Symbol: StoreItem", lines3[1])
			}

			if strings.Contains(lines3[2], "Type: Fix,") == false {
				t.Errorf("Expected %s to be Type: Fix", lines1[2])
			}

			if strings.Contains(lines3[2], "Symbol: StoreItem,") == false {
				t.Errorf("Expected %s to be Symbol: StoreItem", lines3[2])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines3), 3)
		}

		if err9 != nil {
			t.Errorf("Expected %v to be nil", err9)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
