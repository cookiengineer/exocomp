package tools

import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import "encoding/json"
import "os"
import "path/filepath"
import "strings"
import "testing"
import "time"

func TestHumans_AskPermission(t *testing.T) {

	folder := t.TempDir()
	file := filepath.Join(folder, "file.txt")

	err0 := os.WriteFile(file, []byte("test"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create test file: %v", err0)
	}

	capabilities := utils_capabilities.NewSet()
	tool := NewHumansWithCapabilities([]string{"AskPermission"}, ".", ".", "Cody Codera", capabilities)

	results := make(chan string, 1)
	errors := make(chan error, 1)

	go func() {

		result, err1 := tool.AskPermission("The target file is outside the sandbox.", []types.Capability{
			{Type: types.CapabilityFilesRead, Scope: file},
		})

		if err1 == nil {
			results <- result
		} else {
			errors <- err1
		}

	}()

	question := waitForPermissionQuestion(t, tool)

	if question.Agent != "Cody Codera" {
		t.Errorf("Expected agent %q, got %q", "Cody Codera", question.Agent)
	}

	if question.Type != "Permission" {
		t.Errorf("Expected type %q, got %q", "Permission", question.Type)
	}

	if len(question.Capabilities) != 1 || question.Capabilities[0].Type != types.CapabilityFilesRead {
		t.Errorf("Expected one files.Read capability, got %v", question.Capabilities)
	}

	answer, err2 := json.Marshal(question.Capabilities)

	if err2 != nil {
		t.Fatalf("Cannot encode permission answer: %v", err2)
	}

	_, err3 := tool.Answer(question.Question, string(answer))

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

	select {
	case result := <-results:

		if strings.Contains(result, "Granted 1 capability") == false {
			t.Errorf("Expected grant report, got %s", result)
		}

	case err4 := <-errors:
		t.Fatalf("Expected %v to be nil", err4)

	case <-time.After(2 * time.Second):
		t.Fatalf("Timed out waiting for permission result")
	}

	if capabilities.CanRead(file) == false {
		t.Errorf("Expected file to be readable after grant")
	}

}

func TestHumans_AskPermissionDenyAll(t *testing.T) {

	folder := t.TempDir()
	capabilities := utils_capabilities.NewSet()
	tool := NewHumansWithCapabilities([]string{"AskPermission"}, ".", ".", "Cody Codera", capabilities)

	results := make(chan string, 1)

	go func() {
		result, _ := tool.AskPermission("Need folder access.", []types.Capability{
			{Type: types.CapabilityFilesRead, Scope: folder},
		})
		results <- result
	}()

	question := waitForPermissionQuestion(t, tool)
	_, err0 := tool.Answer(question.Question, "[]")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	select {
	case result := <-results:

		if strings.Contains(result, "No capabilities granted") == false {
			t.Errorf("Expected deny report, got %s", result)
		}

	case <-time.After(2 * time.Second):
		t.Fatalf("Timed out waiting for permission result")
	}

	if capabilities.CanRead(folder) == true {
		t.Errorf("Expected folder not to be readable after deny")
	}

}

func TestHumans_AskPermissionAlreadyGranted(t *testing.T) {

	folder := t.TempDir()
	capabilities := utils_capabilities.NewSet()

	err0 := capabilities.Grant([]types.Capability{
		{Type: types.CapabilityFilesRead, Scope: folder},
	})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	tool := NewHumansWithCapabilities([]string{"AskPermission"}, ".", ".", "Cody Codera", capabilities)

	result, err1 := tool.AskPermission("Need folder access.", []types.Capability{
		{Type: types.CapabilityFilesRead, Scope: folder},
	})

	if err1 == nil {
		if strings.Contains(result, "already granted") == false {
			t.Errorf("Expected already-granted report, got %s", result)
		}
	} else {
		t.Errorf("Expected %v to be nil", err1)
	}

}

func waitForPermissionQuestion(t *testing.T, tool *Humans) *types.Question {

	t.Helper()

	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {

		pending := tool.Snapshot()

		for _, question := range pending {

			if question.Type == "Permission" {
				return question
			}

		}

		time.Sleep(10 * time.Millisecond)

	}

	t.Fatalf("Timed out waiting for permission question")

	return nil

}

func TestHumansSchema_AskPermission(t *testing.T) {

	for _, schema := range HumansSchema {

		if schema.Function.Name == "humans.AskPermission" {

			capabilities := schema.Function.Parameters.Properties["capabilities"]

			if capabilities.Items == nil {
				t.Fatalf("Expected capabilities item schema")
			}

			if len(capabilities.Items.Properties) != 2 {
				t.Errorf("Expected type and scope properties, got %v", capabilities.Items.Properties)
			}

			return

		}

	}

	t.Fatalf("Expected humans.AskPermission schema")

}
