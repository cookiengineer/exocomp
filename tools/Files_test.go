package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestFiles_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewFiles("coder", sandbox)

	if tool != nil {

		result1, err1 := tool.List(".")
		_, err2 := tool.Write("./First.txt", "This is the first file content!")
		result3, err3 := tool.List(".")
		_, err4 := tool.Write("./2nd.txt", "This is the second file content!")
		result5, err5 := tool.List(".")

		if result1 != "" || err1 == nil {
			t.Errorf("Expected folder %s to not exist", sandbox)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		lines3 := strings.Split(result3, "\n")

		if len(lines3) == 2 {

			if lines3[0] !=  "files.List: . contains 1 entries." {
				t.Errorf("Expected %d folder entries", 1)
			}

			if strings.Contains(lines3[1], "Name: First.txt") == false {
				t.Errorf("Expected entry %s to be Name: First.txt", lines3[1])
			}

			if strings.Contains(lines3[1], "Type: file") == false {
				t.Errorf("Expected entry %s to be Type: file", lines3[1])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines3), 2)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		lines5 := strings.Split(result5, "\n")

		if len(lines5) == 3 {

			if lines5[0] !=  "files.List: . contains 2 entries." {
				t.Errorf("Expected %d folder entries", 2)
			}

			if strings.Contains(lines5[1], "Name: 2nd.txt") == false {
				t.Errorf("Expected entry %s to be Name: 2nd.txt", lines5[1])
			}

			if strings.Contains(lines5[1], "Type: file") == false {
				t.Errorf("Expected entry %s to be Type: file", lines5[1])
			}

			if strings.Contains(lines5[2], "Name: First.txt") == false {
				t.Errorf("Expected entry %s to be Name: First.txt", lines5[2])
			}

			if strings.Contains(lines5[2], "Type: file") == false {
				t.Errorf("Expected entry %s to be Type: file", lines5[2])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines5), 3)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err4)
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

func TestFiles_Read(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewFiles("coder", sandbox)

	if tool != nil {

		_, err01 := tool.Write("./file.txt", "This is the file content!")
		_, err02 := tool.Write("./file.txt", "This is the file content!")
		result1, err1 := tool.Read("./file.txt")
		result2, err2 := tool.Read("./../../../file.txt")
		result3, err3 := tool.Read("./..\\..\\../file.txt")
		result4, err4 := tool.Read("/etc/passwd")
		result5, err5 := tool.Read("../../../../../../../../etc/passwd")

		if err01 != nil {
			t.Errorf("Expected %v to be nil", err01)
		}

		if err02 != nil {
			t.Errorf("Expected %v to be nil", err02)
		}

		if strings.Contains(result1, "This is the file content!") == false {
			t.Errorf("Expected %s to be allowed", "./file.txt")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result2 != "" {
			t.Errorf("Expected %s to be empty", result2)
		}

		if err2 != nil {

			if strings.Contains(err2.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err2)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err2)
		}

		if result3 != "" {
			t.Errorf("Expected %s to be empty", result3)
		}

		if err3 != nil {

			if strings.Contains(err3.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err3)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err3)
		}

		if result4 != "" {
			t.Errorf("Expected %s to be empty", result4)
		}

		if err4 != nil {

			if strings.Contains(err4.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err4)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err4)
		}

		if result5 != "" {
			t.Errorf("Expected %s to be empty", result5)
		}

		if err5 != nil {

			if strings.Contains(err5.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err5)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err5)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestFiles_Stat(t *testing.T) {
	// TODO
	t.Errorf("TODO: Implement Files.Stat unit test")
}

func TestFiles_Write(t *testing.T) {
	// TODO
	t.Errorf("TODO: Implement Files.Write unit test")
}
