package commands

import "exocomp/schemas"
import "testing"

func testCatalog() *Catalog {

	properties := map[string]schemas.ToolFunctionParameterProperty{
		"path": {
			Type:        "string",
			Description: "Relative path",
		},
		"count": {
			Type:        "integer",
			Description: "Number of items",
		},
		"recursive": {
			Type:        "boolean",
			Description: "Recurse",
		},
		"tags": {
			Type:        "array",
			Description: "Tags",
		},
	}

	return NewCatalog([]schemas.Tool{
		{
			Type: "function",
			Function: schemas.ToolFunction{
				Name:        "files.Read",
				Description: "Read a file",
				Parameters: schemas.ToolFunctionParameters{
					Type:       "object",
					Properties: properties,
					Required:   []string{"path"},
				},
			},
		},
	})

}

func TestResolveTypedArguments(t *testing.T) {

	command, err := Resolve(testCatalog(), `/files.Read path="./foo.txt" count=3 recursive=true tags=["a","b"]`)

	if err == nil {

		if command.Name != "files.Read" {
			t.Errorf("Expected files.Read, got %s", command.Name)
		}

		if command.Method != "Read" {
			t.Errorf("Expected Read, got %s", command.Method)
		}

		if command.Arguments["path"] != "./foo.txt" {
			t.Errorf("Expected path ./foo.txt, got %v", command.Arguments["path"])
		}

		if command.Arguments["count"] != int64(3) {
			t.Errorf("Expected int64(3), got %T %v", command.Arguments["count"], command.Arguments["count"])
		}

		if command.Arguments["recursive"] != true {
			t.Errorf("Expected true, got %v", command.Arguments["recursive"])
		}

	} else {
		t.Errorf("Expected nil, got %v", err)
	}

}

func TestResolveMissingRequired(t *testing.T) {

	_, err := Resolve(testCatalog(), `/files.Read count=3`)

	if err != nil {
		if err.Error() != "commands.Resolve: Missing required parameter \"path\" for \"files.Read\"." {
			t.Errorf("Expected missing path error, got %v", err)
		}
	} else {
		t.Errorf("Expected an error for missing required parameter")
	}

}

func TestResolveUnknownParameter(t *testing.T) {

	_, err := Resolve(testCatalog(), `/files.Read path="x" missing=1`)

	if err != nil {
		if err.Error() != "commands.Resolve: Unknown parameter \"missing\" for \"files.Read\"." {
			t.Errorf("Expected unknown parameter error, got %v", err)
		}
	} else {
		t.Errorf("Expected an error for unknown parameter")
	}

}

func TestCompleteParameter(t *testing.T) {

	suggestions := testCatalog().Complete("/files.Read ", 8)

	if len(suggestions) == 4 {
		if suggestions[0].Kind != "parameter" {
			t.Errorf("Expected parameter suggestion, got %s", suggestions[0].Kind)
		}
	} else {
		t.Errorf("Expected 4 parameter suggestions, got %d", len(suggestions))
	}

}
