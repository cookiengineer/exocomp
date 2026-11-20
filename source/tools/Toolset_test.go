package tools

import "exocomp/schemas"
import "exocomp/tools/bugs"
import "exocomp/tools/changelog"
import "exocomp/tools/containers"
import "exocomp/tools/files"
import "exocomp/tools/humans"
import "exocomp/tools/programs"
import "exocomp/tools/requirements"
import "exocomp/tools/skills"
import "exocomp/tools/websites"
import "encoding/json"
import "strings"
import "testing"

func TestToolSchemas_DoNotMentionInternalState(t *testing.T) {

	schema_groups := [][]schemas.Tool{
		AgentsSchema,
		AssetsSchema,
		bugs.BugsSchema,
		changelog.ChangelogSchema,
		containers.ContainersSchema,
		files.FilesSchema,
		FindingsSchema,
		humans.HumansSchema,
		programs.ProgramsSchema,
		requirements.RequirementsSchema,
		ShellSchema,
		skills.SkillsSchema,
		websites.WebsitesSchema,
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
