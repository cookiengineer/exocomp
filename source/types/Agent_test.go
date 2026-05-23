package types

import "exocomp/encoding/yaml"
import "strings"
import "testing"

func TestUnmarshal_Agent(t *testing.T) {

	file := `
name:   "Peanut Hamper"
role:   pentester
model:  qwen3-coder:30b
prompt: |
  You are a C and C++ Pentest Expert
  running within an Agentic Environment
  and you help to develop exploits for
  Linux. Your preferred languages are
  C, CGo and Go.
temperature: 0.75
allowed-programs:
  - go
  - gofmt
  - gopls
allowed-tools:
  - files.Copy
  - files.Read
  - files.Stat
  - files.Write
`

	agent := Agent{}
	err   := yaml.Unmarshal([]byte(file), &agent)

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if agent.Name != "Peanut Hamper" {
		t.Errorf("Expected agent.Name \"%s\" to be \"%s\"", agent.Name, "Peanut Hamper")
	}

	if agent.Role != "pentester" {
		t.Errorf("Expected agent.Role \"%s\" to be \"%s\"", agent.Role, "pentester")
	}

	if agent.Model != "qwen3-coder:30b" {
		t.Errorf("Expected agent.Model \"%s\" to be \"%s\"", agent.Model, "qwen3-coder:30b")
	}

	if !strings.HasPrefix(agent.Prompt, "You are a C and C++ Pentest Expert") || !strings.HasSuffix(agent.Prompt, "C, CGo and Go.") {
		t.Errorf("Expected agent.Prompt \"%s\"", agent.Prompt)
	}

	if agent.Temperature != 0.75 {
		t.Errorf("Expected agent.Temperature \"%.2f\" to be \"%.2f\"", agent.Temperature, 0.75)
	}

	if len(agent.AllowedPrograms) == 3 {

		if agent.AllowedPrograms[0] != "go" {
			t.Errorf("Expected agent.AllowedPrograms[0] \"%s\" to be \"%s\"", agent.AllowedPrograms[0], "go")
		}

		if agent.AllowedPrograms[1] != "gofmt" {
			t.Errorf("Expected agent.AllowedPrograms[1] \"%s\" to be \"%s\"", agent.AllowedPrograms[1], "gofmt")
		}

		if agent.AllowedPrograms[2] != "gopls" {
			t.Errorf("Expected agent.AllowedPrograms[2] \"%s\" to be \"%s\"", agent.AllowedPrograms[2], "gopls")
		}

	} else {
		t.Errorf("Expected %d agent.AllowedPrograms to be %d", len(agent.AllowedPrograms), 3)
	}

	if len(agent.AllowedTools) == 4 {

		if agent.AllowedTools[0] != "files.Copy" {
			t.Errorf("Expected agent.AllowedTools[0] \"%s\" to be \"%s\"", agent.AllowedTools[0], "files.Copy")
		}

		if agent.AllowedTools[1] != "files.Read" {
			t.Errorf("Expected agent.AllowedTools[1] \"%s\" to be \"%s\"", agent.AllowedTools[1], "files.Read")
		}

		if agent.AllowedTools[2] != "files.Stat" {
			t.Errorf("Expected agent.AllowedTools[0] \"%s\" to be \"%s\"", agent.AllowedTools[2], "files.Stat")
		}

		if agent.AllowedTools[3] != "files.Write" {
			t.Errorf("Expected agent.AllowedTools[0] \"%s\" to be \"%s\"", agent.AllowedTools[3], "files.Write")
		}

	} else {
		t.Errorf("Expected %d agent.AllowedTools to be %d", len(agent.AllowedTools), 4)
	}

}

