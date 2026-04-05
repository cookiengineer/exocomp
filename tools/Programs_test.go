package tools

import "os"
import "strings"
import "testing"

func TestPrograms(t *testing.T) {

	t.Run("TestProgram(ls -la)", func(t *testing.T) {

		cwd, err1 := os.Getwd()

		if err1 == nil {

			programs := NewPrograms("tester", cwd, []string{"ls"})

			result, err2 := programs.Call("Execute", map[string]interface{}{
				"program":   "ls",
				"arguments": []interface{}{
					"-la",
				},
			})

			if err2 == nil {

				if strings.Contains(result, "-rw-r--r--") {
					// Do Nothing
				} else {
					t.Errorf("Expected %s to contain -rw-r--r--", result)
				}

			} else {
				t.Errorf("Expected %s to be nil", err2.Error())
			}

		} else {
			t.Errorf("Expected %s to be nil", err1.Error())
		}

	})

}

