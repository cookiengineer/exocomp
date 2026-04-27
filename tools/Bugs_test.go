package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestBugs_Add(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-bugs-*")
	sandbox       := filepath.Join(playground, "bugs")
	tool          := NewBugs(playground, sandbox)

	if tool != nil {

		result1, err1 := tool.Add("./path/to/File.go", "Whatever", "A new feature broke the Whatever method.")
		result2, err2 := tool.Add("./path/to/File.go", "Whatever", "A new feature broke the Whatever method.")

		if result1 != "bugs.Add: Bug report with 40 B written." {
			t.Errorf("Expected Bug report to be written")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result2 != "bugs.Add: Bug report with 40 B updated." {
			t.Errorf("Expected Bug report to be written")
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestBugs_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-bugs-*")
	sandbox       := filepath.Join(playground, "bugs")
	tool          := NewBugs(playground, sandbox)

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

		if strings.HasPrefix(result, "bugs.List: 2 unfixed bug reports.") == false {
			t.Errorf("Expected 2 unfixed bug reports:\n%s", result)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestBugs_Fix(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-bugs-*")
	sandbox       := filepath.Join(playground, "bugs")
	tool          := NewBugs(playground, sandbox)

	if tool != nil {

		_, err1 := tool.Add("./path/to/First.go",  "Foo", "A new feature broke the Foo method.")
		_, err2 := tool.Add("./path/to/Second.go", "Bar", "A new feature broke the Bar method.")
		_, err3 := tool.Add("./path/to/Third.go",  "Qux", "A new feature broke the Qux method.")
		_, err4 := tool.Fix("./path/to/Third.go",  "Qux")

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

		if strings.HasPrefix(result, "bugs.List: 2 unfixed bug reports.") == false {
			t.Errorf("Expected 2 unfixed bug reports:\n%s", result)
		}

		if strings.Contains(result, "path/to/Third.go") == true {
			t.Errorf("Expected %s to be fixed.", "./path/to/Third.go")
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestBugs_Search(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-bugs-*")
	sandbox       := filepath.Join(playground, "bugs")
	tool          := NewBugs(playground, sandbox)

	if tool != nil {

		_, err1 := tool.Add("./path/to/First.go",  "Foo", "A new feature broke the Foo method.")
		_, err2 := tool.Add("./path/to/Second.go", "Bar", "A new feature broke the Bar method.")
		_, err3 := tool.Add("./path/to/Third.go",  "Qux", "A new feature broke the Qux method.")
		_, err4 := tool.Add("./path/to/Third.go",  "Doo", "A new feature broke the Doo method.")

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

		result1, err5 := tool.Search("./path/to/Third.go", "Qux")
		result2, err6 := tool.Search("./path/to/Third.go", "")

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if err6 != nil {
			t.Errorf("Expected %v to be nil", err6)
		}

		if strings.HasPrefix(result1, "bugs.Search: ./path/to/Third.go#Qux contains 1 bug reports.") == false {
			t.Errorf("Expected 1 unfixed bug report:\n%s", result1)
		}

		if strings.HasPrefix(result2, "bugs.Search: ./path/to/Third.go contains 2 bug reports.") == false {
			t.Errorf("Expected 2 unfixed bug reports:\n%s", result2)
		}

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})
}
