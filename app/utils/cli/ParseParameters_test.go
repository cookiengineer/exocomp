package cli

import "testing"

func TestParseParameters_String(t *testing.T) {

	result := ParseParameters("name=\"hello world\"")

	if result["name"] != "hello world" {
		t.Errorf("Expected %q, got %v", "hello world", result["name"])
	}

}

func TestParseParameters_Number(t *testing.T) {

	result := ParseParameters("step=42")

	if result["step"] != float64(42) {
		t.Errorf("Expected %v, got %v", float64(42), result["step"])
	}

}

func TestParseParameters_Boolean(t *testing.T) {

	result := ParseParameters("force=true")

	if result["force"] != true {
		t.Errorf("Expected %v, got %v", true, result["force"])
	}

}

func TestParseParameters_Object(t *testing.T) {

	result := ParseParameters("payload={\"key\":\"value\"}")

	object, ok := result["payload"].(map[string]any)

	if ok == false {
		t.Fatalf("Expected payload to be an object, got %v", result["payload"])
	}

	if object["key"] != "value" {
		t.Errorf("Expected %q, got %v", "value", object["key"])
	}

}

func TestParseParameters_Multiple(t *testing.T) {

	result := ParseParameters("name=\"foo\" step=1 force=false")

	if result["name"] != "foo" {
		t.Errorf("Expected %q, got %v", "foo", result["name"])
	}

	if result["step"] != float64(1) {
		t.Errorf("Expected %v, got %v", float64(1), result["step"])
	}

	if result["force"] != false {
		t.Errorf("Expected %v, got %v", false, result["force"])
	}

}
