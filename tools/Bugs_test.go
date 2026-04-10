package tools

import "fmt"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestBugs_Add(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewBugs("tester", sandbox, playground)

	if tool != nil {

		result1, err1 := tool.Add("./path/to/File.go", "Whatever", "A new feature broke the Whatever method.")

		if result1 != "bugs.Add: Bug report with 40 B written." {
			t.Errorf("Expected Bug report to be written")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

}

func TestBugs_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewBugs("tester", sandbox, playground)

	if tool != nil {

		_, err1 := tool.Add("./path/to/First.go",  "Foo", "A new feature broke the Foo method.")
		_, err2 := tool.Add("./path/to/Second.go", "Bar", "A new feature broke the Bar method.")
		_, err3 := tool.Add("./path/to/Third.go",  "Qux", "A new feature broke the Qux method.")
		_, err4 := tool.Fix("./path/to/Second.go", "Bar")

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

		result, err5 := tool.List()

		if strings.HasPrefix(result, "bugs.List: 2 unfixed bug reports.") != true {
			t.Errorf("Expected 2 unfixed bug reports:\n%s", result)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

}

func TestBugs_Fix(t *testing.T) {

	// TODO: Test Fix
	fmt.Println("TODO: Bugs.Fix()")

}

func TestBugs_Search(t *testing.T) {

	// TODO: Test Search
	fmt.Println("TODO: Bugs.Search()")

}
