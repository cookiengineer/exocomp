package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestFiles_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "files")
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
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines3[1], "Name: First.txt")
			}

			if strings.Contains(lines3[1], "Type: file") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines3[1], "Type: file")
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
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[1], "Name: 2nd.txt")
			}

			if strings.Contains(lines5[1], "Type: file") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[1], "Type: file")
			}

			if strings.Contains(lines5[2], "Name: First.txt") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[2], "Name: First.txt")
			}

			if strings.Contains(lines5[2], "Type: file") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[2], "Type: file")
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
	sandbox       := filepath.Join(playground, "files")
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

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestFiles_Stat(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "files")
	tool          := NewFiles("coder", sandbox)

	if tool != nil {

		_, err0 := tool.Write("./file.txt", "This is the file content!")

		if err0 != nil {
			t.Errorf("Expected %v to be nil", err0)
		}

		result1, err1 := tool.Stat("./does-not-exist.txt")
		result2, err2 := tool.Stat("./file.txt")

		if result1 != "" {
			t.Errorf("Expected %s to be empty", result1)
		}

		if err1 != nil {

			if strings.Contains(err1.Error(), "File doesn't exist") == false {
				t.Errorf("Expected %v to be file doesn't exist error", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

		lines2 := strings.Split(result2, "\n")

		if len(lines2) == 6 {

			if lines2[0] != "files.Stat: ./file.txt is a file." {
				t.Errorf("Expected \"%s\" to be a file", lines2[0])
			}

			if strings.Contains(lines2[1], "Name: file.txt") == false {
				t.Errorf("Expected \"%s\" to be \"%s\"", lines2[1], "Name: file.txt")
			}

			if strings.Contains(lines2[2], "Type: file") == false {
				t.Errorf("Expected \"%s\" to be \"%s\"", lines2[2], "Type: file")
			}

			if strings.Contains(lines2[3], "Size: 26 B") == false {
				t.Errorf("Expected \"%s\" to be \"%s\"", lines2[3], "Size: 26 B")
			}

			if strings.Contains(lines2[4], "Mode: file (readable, writable)") == false {
				t.Errorf("Expected \"%s\" to be \"%s\"", lines2[4], "Mode: file (readable, writable)")
			}

			if strings.Contains(lines2[5], "Modified: ") == false {
				t.Errorf("Expected \"%s\" to contain \"%s\"", lines2[5], "Modified: YYYY-MM-DD HH:ii:ss")
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines2), 6)
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

func TestFiles_Write(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "files")
	tool          := NewFiles("coder", sandbox)

	if tool != nil {

		result1, err1 := tool.Write("./file.txt", "This is the file content!")
		result2, err2 := tool.Write("./../../../file.txt", "This is the file content!")
		result3, err3 := tool.Write("/etc/passwd", "This is the file content!")

		if result1 != "files.Write: ./file.txt with 26 B written." {
			t.Errorf("Expected file to be written")
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
