package tools

import "testing"

func TestToolset_RegistersAssetsAndFindingsImport(t *testing.T) {

	playground := t.TempDir()

	registered := Toolset(
		playground,
		playground,
		"test-model",
		nil,
		false,
		nil,
		nil,
		[]string{"assets.Add", "assets.Import", "findings.Import"},
		"tester",
		nil,
	)

	names := make(map[string]bool)

	for _, tool := range registered {

		names[tool.Name()] = true

		schemas := tool.Schemas()

		if len(schemas) == 0 {
			t.Errorf("Expected schemas for tool %q", tool.Name())
		}

	}

	if names["assets"] == false {
		t.Errorf("Expected assets tool to be registered")
	}

	if names["findings"] == false {
		t.Errorf("Expected findings tool to be registered")
	}

}
