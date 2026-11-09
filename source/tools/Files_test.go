package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestFiles_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

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

			if lines3[0] != "files.List: \".\" contains 1 entries." {
				t.Errorf("Expected %d folder entries", 1)
			}

			if strings.Contains(lines3[1], "Path: \"./First.txt\"") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines3[1], "Path: \"./First.txt\"")
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

			if lines5[0] != "files.List: \".\" contains 2 entries." {
				t.Errorf("Expected %d folder entries", 2)
			}

			if strings.Contains(lines5[1], "Path: \"./2nd.txt\"") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[1], "Path: \"./2nd.txt\"")
			}

			if strings.Contains(lines5[1], "Type: file") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[1], "Type: file")
			}

			if strings.Contains(lines5[2], "Path: \"./First.txt\"") == false {
				t.Errorf("Expected entry \"%s\" to contain \"%s\"", lines5[2], "Path: \"./First.txt\"")
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
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

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
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

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

			if strings.Contains(err1.Error(), "File \"./does-not-exist.txt\" does not exist.") == false {
				t.Errorf("Expected %v to be file doesn't exist error", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

		lines2 := strings.Split(result2, "\n")

		if len(lines2) == 6 {

			if lines2[0] != "files.Stat: \"./file.txt\" is a file." {
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
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

	if tool != nil {

		result1, err1 := tool.Write("./file.txt", "This is the file content!")
		result2, err2 := tool.Write("./../../../file.txt", "This is the file content!")
		result3, err3 := tool.Write("/etc/passwd", "This is the file content!")

		if result1 != "files.Write: File \"./file.txt\" with 26 B written." {
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

func TestFiles_Copy(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

	if tool != nil {

		_, err1 := tool.Write("./source.txt", "This is the source file content!")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		result2, err2 := tool.Copy("./source.txt", "./copy/source.txt")
		result3, err3 := tool.Read("./copy/source.txt")
		result4, err4 := tool.Copy("./../../../source.txt", "./escape.txt")

		if strings.Contains(result2, "File \"./source.txt\" copied to \"./copy/source.txt\"") == false {
			t.Errorf("Expected file to be copied, got %s", result2)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "This is the source file content!") == false {
			t.Errorf("Expected copied file to contain source content, got %s", result3)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
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

func TestFiles_ReadSymbol(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "ReadSymbol", "Stat", "Write", "WriteSymbol"}, playground, sandbox)

	if tool != nil {

		_, err1 := tool.Write("./core.go", "package core\n\nfunc FirstFunction(current int64) (string, error) {\n\treturn \"\", nil\n}\n\ntype MyByte uint8\n")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		result2, err2 := tool.ReadSymbol("./core.go", "FirstFunction")
		result3, err3 := tool.ReadSymbol("./core.go", "MyByte")
		result4, err4 := tool.ReadSymbol("./core.go", "Missing")

		if strings.Contains(result2, "Symbol \"FirstFunction\" with Type \"func\"") == false {
			t.Errorf("Expected function symbol to be read, got %s", result2)
		}

		if strings.Contains(result2, "func FirstFunction(current int64) (string, error)") == false {
			t.Errorf("Expected function declaration in result, got %s", result2)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "Symbol \"MyByte\" with Type \"type\"") == false {
			t.Errorf("Expected type symbol to be read, got %s", result3)
		}

		if strings.Contains(result3, "type MyByte uint8") == false {
			t.Errorf("Expected type declaration in result, got %s", result3)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if result4 != "" {
			t.Errorf("Expected %s to be empty", result4)
		}

		if err4 == nil {
			t.Errorf("Expected %v to be not nil", err4)
		} else if strings.Contains(err4.Error(), "has no Symbol") == false {
			t.Errorf("Expected missing symbol error, got %v", err4)
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

func TestFiles_Search(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "ReadSymbol", "Search", "Stat", "Write", "WriteSymbol"}, playground, sandbox)

	if tool != nil {

		_, err01 := tool.Write("./core.go", "package core\n\nfunc ParseInput() {\n\treturn\n}\n\ntype MyByte uint8\n")
		_, err02 := tool.Write("./utils/Parse.go", "package utils\n\nfunc ParseOutput() string {\n\treturn \"\"\n}\n")
		_, err03 := tool.Write("./sub/Deep.go", "package sub\n\nfunc DeepSearch() {\n\treturn\n}\n")

		if err01 != nil {
			t.Errorf("Expected %v to be nil", err01)
		}

		if err02 != nil {
			t.Errorf("Expected %v to be nil", err02)
		}

		if err03 != nil {
			t.Errorf("Expected %v to be nil", err03)
		}

		result1, err1 := tool.Search(".", "ParseOutput")
		result2, err2 := tool.Search(".", "parseoutput")
		result3, err3 := tool.Search(".", "MyByte")
		result4, err4 := tool.Search(".", "DeepSearch")
		result5, err5 := tool.Search(".", "missingquery")
		result6, err6 := tool.Search("./../../../etc", "ParseOutput")

		if strings.Contains(result1, "files.Search: \".\" for \"ParseOutput\" contains 1 matches.") == false {
			t.Errorf("Expected header with 1 match, got %s", result1)
		}

		if strings.Contains(result1, "File: \"./utils/Parse.go\", Symbol: \"ParseOutput\", Type: \"func\"") == false {
			t.Errorf("Expected match in utils/Parse.go, got %s", result1)
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result2, "File: \"./utils/Parse.go\", Symbol: \"ParseOutput\", Type: \"func\"") == false {
			t.Errorf("Expected case-insensitive match in utils/Parse.go, got %s", result2)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "File: \"./core.go\", Symbol: \"MyByte\", Type: \"type\"") == false {
			t.Errorf("Expected type match in core.go, got %s", result3)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if strings.Contains(result4, "File: \"./sub/Deep.go\", Symbol: \"DeepSearch\", Type: \"func\"") == false {
			t.Errorf("Expected recursive match in sub/Deep.go, got %s", result4)
		}

		if err4 != nil {
			t.Errorf("Expected %v to be nil", err4)
		}

		if strings.Contains(result5, "contains 0 matches.") == false {
			t.Errorf("Expected 0 matches, got %s", result5)
		}

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if result6 != "" {
			t.Errorf("Expected %s to be empty", result6)
		}

		if err6 != nil {

			if strings.Contains(err6.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err6)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err6)
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

func TestFiles_WriteSymbol(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "ReadSymbol", "Stat", "Write", "WriteSymbol"}, playground, sandbox)

	if tool != nil {

		_, err1 := tool.Write("./core.go", "package core\n\nfunc FirstFunction() int {\n\treturn 1\n}\n\ntype MyByte uint8\n")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		result2, err2 := tool.WriteSymbol("./core.go", "FirstFunction", "func FirstFunction() int {\n\treturn 2\n}")
		result3, err3 := tool.WriteSymbol("./core.go", "MyByte", "type MyByte string")
		result4, err4 := tool.WriteSymbol("./core.go", "Missing", "func Missing() {}")

		if strings.Contains(result2, "Symbol \"FirstFunction\" written") == false {
			t.Errorf("Expected function symbol to be written, got %s", result2)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result3, "Symbol \"MyByte\" written") == false {
			t.Errorf("Expected type symbol to be written, got %s", result3)
		}

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if result4 != "" {
			t.Errorf("Expected %s to be empty", result4)
		}

		if err4 == nil {
			t.Errorf("Expected %v to be not nil", err4)
		} else if strings.Contains(err4.Error(), "has no Symbol") == false {
			t.Errorf("Expected missing symbol error, got %v", err4)
		}

		content, err5 := tool.Read("./core.go")

		if err5 != nil {
			t.Errorf("Expected %v to be nil", err5)
		}

		if strings.Contains(content, "return 2") == false {
			t.Errorf("Expected function body to be overwritten, got %s", content)
		}

		if strings.Contains(content, "return 1") == true {
			t.Errorf("Expected old function body to be gone, got %s", content)
		}

		if strings.Contains(content, "type MyByte string") == false {
			t.Errorf("Expected type to be overwritten, got %s", content)
		}

		if strings.Contains(content, "type MyByte uint8") == true {
			t.Errorf("Expected old type to be gone, got %s", content)
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

func TestFiles_Symlinks(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-symlink-*")
	sandbox := filepath.Join(playground, "files")
	outside := filepath.Join(playground, "outside")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	err1 := os.MkdirAll(outside, 0755)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	err2 := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("TOP SECRET"), 0644)

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	err3 := os.WriteFile(filepath.Join(outside, "secret.go"), []byte("package secret\n\nfunc LeakedSymbol() {\n\treturn\n}\n"), 0644)

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	err4 := os.WriteFile(filepath.Join(sandbox, "normal.txt"), []byte("normal content"), 0644)

	if err4 != nil {
		t.Errorf("Expected %v to be nil", err4)
	}

	err5 := os.WriteFile(filepath.Join(sandbox, "internal.go"), []byte("package files\n\nfunc InternalSymbol() {\n\treturn\n}\n"), 0644)

	if err5 != nil {
		t.Errorf("Expected %v to be nil", err5)
	}

	err6 := os.MkdirAll(filepath.Join(sandbox, "tree"), 0755)

	if err6 != nil {
		t.Errorf("Expected %v to be nil", err6)
	}

	err7 := os.WriteFile(filepath.Join(sandbox, "tree", "normal.txt"), []byte("tree normal"), 0644)

	if err7 != nil {
		t.Errorf("Expected %v to be nil", err7)
	}

	createTestSymlink(t, filepath.Join(outside, "secret.txt"), filepath.Join(sandbox, "leak.txt"))
	createTestSymlink(t, filepath.Join(outside, "secret.go"), filepath.Join(sandbox, "leak.go"))
	createTestSymlink(t, outside, filepath.Join(sandbox, "leakdir"))
	createTestSymlink(t, "normal.txt", filepath.Join(sandbox, "internal_link.txt"))
	createTestSymlink(t, filepath.Join(outside, "secret.txt"), filepath.Join(sandbox, "tree", "link_outside.txt"))
	createTestSymlink(t, outside, filepath.Join(sandbox, "tree", "link_outside_dir"))
	createTestSymlink(t, filepath.Join(sandbox, "normal.txt"), filepath.Join(sandbox, "tree", "link_internal.txt"))

	tool := NewFiles([]string{"Copy", "List", "Read", "ReadSymbol", "Search", "Stat", "Write", "WriteSymbol"}, playground, sandbox)

	if tool == nil {
		t.Fatalf("Expected tool to be not nil")
	}

	result1, err8 := tool.Read("./leak.txt")
	result2, err9 := tool.Read("./leakdir/secret.txt")
	result3, err10 := tool.Read("./internal_link.txt")
	result4, err11 := tool.ReadSymbol("./leak.go", "LeakedSymbol")
	result5, err12 := tool.Stat("./leak.txt")

	if result1 != "" || err8 == nil {
		t.Errorf("Expected outside file symlink to be rejected")
	} else if strings.Contains(err8.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err8)
	} else if strings.Contains(err8.Error(), outside) {
		t.Errorf("Expected error not to leak outside path, got %v", err8)
	}

	if result2 != "" || err9 == nil {
		t.Errorf("Expected outside dir symlink child to be rejected")
	} else if strings.Contains(err9.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err9)
	}

	if result3 != "" || err10 == nil {
		t.Errorf("Expected internal symlink to be rejected")
	} else if strings.Contains(err10.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err10)
	}

	if result4 != "" || err11 == nil {
		t.Errorf("Expected ReadSymbol on symlink to be rejected")
	} else if strings.Contains(err11.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err11)
	}

	if result5 != "" || err12 == nil {
		t.Errorf("Expected Stat on symlink to be rejected")
	} else if strings.Contains(err12.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err12)
	}

	result6, err13 := tool.Write("./leak.txt", "pwned")

	if result6 != "" || err13 == nil {
		t.Errorf("Expected Write through outside file symlink to be rejected")
	}

	content1, err14 := os.ReadFile(filepath.Join(outside, "secret.txt"))

	if err14 != nil {
		t.Errorf("Expected %v to be nil", err14)
	} else if string(content1) != "TOP SECRET" {
		t.Errorf("Expected outside file to be unchanged, got %q", string(content1))
	}

	result7, err15 := tool.Write("./leakdir/new.txt", "pwned")

	if result7 != "" || err15 == nil {
		t.Errorf("Expected Write through outside dir symlink to be rejected")
	}

	_, err16 := os.Stat(filepath.Join(outside, "new.txt"))

	if os.IsNotExist(err16) == false {
		t.Errorf("Expected no file to be created outside the sandbox")
	}

	result8, err17 := tool.Write("./internal_link.txt", "pwned")

	if result8 != "" || err17 == nil {
		t.Errorf("Expected Write through internal symlink to be rejected")
	}

	content2, err18 := os.ReadFile(filepath.Join(sandbox, "normal.txt"))

	if err18 != nil {
		t.Errorf("Expected %v to be nil", err18)
	} else if string(content2) != "normal content" {
		t.Errorf("Expected internal file to be unchanged, got %q", string(content2))
	}

	result9, err19 := tool.WriteSymbol("./leak.go", "LeakedSymbol", "func LeakedSymbol() {\n\treturn\n}")

	if result9 != "" || err19 == nil {
		t.Errorf("Expected WriteSymbol on symlink to be rejected")
	}

	content3, err20 := os.ReadFile(filepath.Join(outside, "secret.go"))

	if err20 != nil {
		t.Errorf("Expected %v to be nil", err20)
	} else if strings.Contains(string(content3), "pwned") {
		t.Errorf("Expected outside Go file to be unchanged")
	}

	result10, err21 := tool.List(".")

	if err21 != nil {
		t.Errorf("Expected %v to be nil", err21)
	} else if strings.Contains(result10, "leak.txt") || strings.Contains(result10, "leak.go") || strings.Contains(result10, "leakdir") || strings.Contains(result10, "internal_link.txt") {
		t.Errorf("Expected List to skip symlinks, got %s", result10)
	} else if strings.Contains(result10, "normal.txt") == false {
		t.Errorf("Expected List to contain normal.txt, got %s", result10)
	}

	result11, err22 := tool.Search(".", "LeakedSymbol")

	if err22 != nil {
		t.Errorf("Expected %v to be nil", err22)
	} else if strings.Contains(result11, "contains 0 matches.") == false {
		t.Errorf("Expected Search to skip outside symlink, got %s", result11)
	}

	result12, err23 := tool.Copy("./leak.txt", "./copied.txt")

	if result12 != "" || err23 == nil {
		t.Errorf("Expected Copy of outside file symlink to be rejected")
	}

	_, err24 := os.Stat(filepath.Join(sandbox, "copied.txt"))

	if os.IsNotExist(err24) == false {
		t.Errorf("Expected no copy of outside file symlink")
	}

	result13, err25 := tool.Copy("./leakdir", "./copieddir")

	if result13 != "" || err25 == nil {
		t.Errorf("Expected Copy of outside dir symlink to be rejected")
	}

	_, err26 := os.Stat(filepath.Join(sandbox, "copieddir"))

	if os.IsNotExist(err26) == false {
		t.Errorf("Expected no copy of outside dir symlink")
	}

	result14, err27 := tool.Copy("./tree", "./tree_copy")

	if err27 != nil {
		t.Errorf("Expected %v to be nil, got %v", nil, err27)
	} else if strings.Contains(result14, "copied") == false {
		t.Errorf("Expected tree to be copied, got %s", result14)
	}

	_, err28 := os.Stat(filepath.Join(sandbox, "tree_copy", "normal.txt"))

	if err28 != nil {
		t.Errorf("Expected copied tree to contain normal.txt, got %v", err28)
	}

	_, err29 := os.Stat(filepath.Join(sandbox, "tree_copy", "link_outside.txt"))

	if os.IsNotExist(err29) == false {
		t.Errorf("Expected folder copy to skip symlink entries")
	}

	_, err30 := os.Stat(filepath.Join(sandbox, "tree_copy", "link_outside_dir"))

	if os.IsNotExist(err30) == false {
		t.Errorf("Expected folder copy to skip symlink directories")
	}

	_, err31 := os.Stat(filepath.Join(sandbox, "tree_copy", "link_internal.txt"))

	if os.IsNotExist(err31) == false {
		t.Errorf("Expected folder copy to skip internal symlinks")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
