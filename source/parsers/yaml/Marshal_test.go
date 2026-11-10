package yaml_test

import "exocomp/parsers/yaml"
import "testing"

func TestMarshalSortsObjectKeys(t *testing.T) {

	value := struct {
		Name  string `yaml:"name"`
		Model string `yaml:"model"`
		Debug bool   `yaml:"debug"`
	}{
		Name:  "test",
		Model: "model",
		Debug: true,
	}

	out, err := yaml.Marshal(value)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := "debug: true\nmodel: model\nname: test\n"

	if string(out) != expected {
		t.Fatalf("expected sorted output %q, got: %q", expected, string(out))
	}

}
