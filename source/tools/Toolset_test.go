package tools

import "exocomp/schemas"
import "encoding/json"
import "strings"
import "testing"

func TestToolSchemas_DoNotMentionInternalState(t *testing.T) {

	schema_groups := [][]schemas.Tool{
		AgentsSchema,
		AssetsSchema,
		BugsSchema,
		ChangelogSchema,
		ContainersSchema,
		FilesSchema,
		FindingsSchema,
		HumansSchema,
		ProgramsSchema,
		RequirementsSchema,
		ShellSchema,
		SkillsSchema,
		WebsitesSchema,
	}

	for _, entries := range schema_groups {

		for _, entry := range entries {

			buffer, err := json.Marshal(entry)

			if err != nil {
				t.Fatalf("Expected %v to be nil", err)
			}

			if strings.Contains(string(buffer), ".exocomp") == true {
				t.Errorf("Expected tool schema %T not to mention internal state: %s", entry, string(buffer))
			}

		}

	}

}

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
