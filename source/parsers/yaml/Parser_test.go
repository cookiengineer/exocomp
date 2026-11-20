package yaml_test

import "testing"

import "exocomp/parsers/yaml"

type mockProvider struct {
	URL   string `yaml:"url"`
	Alias string `yaml:"alias"`
	Token string `yaml:"token"`
}

type mockConfig struct {
	Name      string                  `yaml:"name"`
	Role      string                  `yaml:"role"`
	Model     string                  `yaml:"model"`
	Providers map[string]mockProvider `yaml:"providers"`
}

const configYAML = `name: HyRell
role: planner
model: "deepseek-v4-pro:cloud"
providers:
  deepseek-v4-pro:cloud:
    url: "https://api.deepseek.com"
    alias: "deepseek-v4-pro"
    token: "sk-abc123"
`

func TestUnmarshalConfigWithColonProviderKey(t *testing.T) {

	config := mockConfig{}

	err := yaml.Unmarshal([]byte(configYAML), &config)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.Name != "HyRell" {
		t.Fatalf("expected Name %q, got %q", "HyRell", config.Name)
	}

	if config.Role != "planner" {
		t.Fatalf("expected Role %q, got %q", "planner", config.Role)
	}

	if config.Model != "deepseek-v4-pro:cloud" {
		t.Fatalf("expected Model %q, got %q", "deepseek-v4-pro:cloud", config.Model)
	}

	provider, exists := config.Providers["deepseek-v4-pro:cloud"]

	if exists == false {
		t.Fatalf("expected provider %q to exist, got providers: %#v", "deepseek-v4-pro:cloud", config.Providers)
	}

	if provider.URL != "https://api.deepseek.com" {
		t.Fatalf("expected provider URL %q, got %q", "https://api.deepseek.com", provider.URL)
	}

	if provider.Alias != "deepseek-v4-pro" {
		t.Fatalf("expected provider Alias %q, got %q", "deepseek-v4-pro", provider.Alias)
	}

	if provider.Token != "sk-abc123" {
		t.Fatalf("expected provider Token %q, got %q", "sk-abc123", provider.Token)
	}

}

func TestParserColonProviderKeyTree(t *testing.T) {

	parser := yaml.NewParser([]byte(configYAML))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	providers, exists := root.ObjectChildren["providers"]

	if exists == false {
		t.Fatalf("expected providers node to exist, got: %#v", root.ObjectChildren)
	}

	if providers.Kind != yaml.ObjectNode {
		t.Fatalf("expected providers to be an object, got kind: %v", providers.Kind)
	}

	provider, exists := providers.ObjectChildren["deepseek-v4-pro:cloud"]

	if exists == false {
		t.Fatalf("expected provider key %q to exist, got: %#v", "deepseek-v4-pro:cloud", providers.ObjectChildren)
	}

	if provider.Kind != yaml.ObjectNode {
		t.Fatalf("expected provider to be an object, got kind: %v", provider.Kind)
	}

	url_node, exists := provider.ObjectChildren["url"]

	if exists == false {
		t.Fatalf("expected url child to exist, got: %#v", provider.ObjectChildren)
	}

	if url_node.Kind != yaml.ScalarNode {
		t.Fatalf("expected url child to be scalar, got kind: %v", url_node.Kind)
	}

	if url_node.Value != "https://api.deepseek.com" {
		t.Fatalf("expected url value %q, got %q", "https://api.deepseek.com", url_node.Value)
	}

}

func TestParserArrayOfScalars(t *testing.T) {

	data := "name: Whatever\nallowed:\n  - go\n  - gofmt\n  - gopls\n"

	parser := yaml.NewParser([]byte(data))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	allowed, exists := root.ObjectChildren["allowed"]

	if exists == false {
		t.Fatalf("expected allowed node to exist, got: %#v", root.ObjectChildren)
	}

	if allowed.Kind != yaml.ArrayNode {
		t.Fatalf("expected allowed to be an array, got kind: %v", allowed.Kind)
	}

	expected := []string{"go", "gofmt", "gopls"}

	if len(allowed.ArrayChildren) != len(expected) {
		t.Fatalf("expected %d array children, got %d", len(expected), len(allowed.ArrayChildren))
	}

	for index, child := range allowed.ArrayChildren {

		if child.Kind != yaml.ScalarNode {
			t.Fatalf("expected array child %d to be scalar, got kind: %v", index, child.Kind)
		}

		if child.Value != expected[index] {
			t.Fatalf("expected array child %d value %q, got %q", index, expected[index], child.Value)
		}

	}

}

func TestParserArrayOfObjects(t *testing.T) {

	data := "servers:\n  - name: foo\n    url: https://foo.example\n  - name: baz\n    url: https://baz.example\n"

	parser := yaml.NewParser([]byte(data))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	servers, exists := root.ObjectChildren["servers"]

	if exists == false {
		t.Fatalf("expected servers node to exist, got: %#v", root.ObjectChildren)
	}

	if servers.Kind != yaml.ArrayNode {
		t.Fatalf("expected servers to be an array, got kind: %v", servers.Kind)
	}

	if len(servers.ArrayChildren) != 2 {
		t.Fatalf("expected 2 server entries, got %d", len(servers.ArrayChildren))
	}

	first := servers.ArrayChildren[0]

	if first.Kind != yaml.ObjectNode {
		t.Fatalf("expected first server to be an object, got kind: %v", first.Kind)
	}

	name_node, exists := first.ObjectChildren["name"]

	if exists == false || name_node.Value != "foo" {
		t.Fatalf("expected first server name %q, got %#v", "foo", name_node)
	}

	url_node, exists := first.ObjectChildren["url"]

	if exists == false || url_node.Value != "https://foo.example" {
		t.Fatalf("expected first server url %q, got %#v", "https://foo.example", url_node)
	}

}

func TestParserArrayOfObjectsWithNestedObject(t *testing.T) {

	data := "servers:\n  - name: foo\n    config:\n      port: 8080\n      tls: true\n  - name: baz\n    config:\n      port: 9090\n"

	parser := yaml.NewParser([]byte(data))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	servers := root.ObjectChildren["servers"]

	if servers == nil || len(servers.ArrayChildren) != 2 {
		t.Fatalf("expected 2 servers, got: %#v", servers)
	}

	first := servers.ArrayChildren[0]
	config := first.ObjectChildren["config"]

	if config == nil || config.Kind != yaml.ObjectNode {
		t.Fatalf("expected first server config object, got: %#v", config)
	}

	port_node := config.ObjectChildren["port"]

	if port_node == nil || port_node.Value != "8080" {
		t.Fatalf("expected first server config port %q, got: %#v", "8080", port_node)
	}

}

func TestParserNestedArrays(t *testing.T) {

	data := "matrix:\n  - - 1\n    - 2\n  - - 3\n    - 4\n"

	parser := yaml.NewParser([]byte(data))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	matrix := root.ObjectChildren["matrix"]

	if matrix == nil || matrix.Kind != yaml.ArrayNode {
		t.Fatalf("expected matrix array, got: %#v", matrix)
	}

	if len(matrix.ArrayChildren) != 2 {
		t.Fatalf("expected 2 matrix rows, got %d", len(matrix.ArrayChildren))
	}

	first_row := matrix.ArrayChildren[0]

	if first_row.Kind != yaml.ArrayNode || len(first_row.ArrayChildren) != 2 {
		t.Fatalf("expected first row to be array of 2, got: %#v", first_row)
	}

	if first_row.ArrayChildren[0].Value != "1" {
		t.Fatalf("expected first row first value %q, got %q", "1", first_row.ArrayChildren[0].Value)
	}

	if first_row.ArrayChildren[1].Value != "2" {
		t.Fatalf("expected first row second value %q, got %q", "2", first_row.ArrayChildren[1].Value)
	}

}

func TestParserDeeplyNestedObjects(t *testing.T) {

	data := "a:\n  b:\n    c:\n      d: value\n"

	parser := yaml.NewParser([]byte(data))
	root, err := parser.Root()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	a := root.ObjectChildren["a"]
	b := a.ObjectChildren["b"]
	c := b.ObjectChildren["c"]
	d := c.ObjectChildren["d"]

	if a == nil || b == nil || c == nil || d == nil {
		t.Fatalf("expected a.b.c.d chain, got: %#v", root)
	}

	if d.Value != "value" {
		t.Fatalf("expected d value %q, got %q", "value", d.Value)
	}

}

func TestRootSurfacesParseError(t *testing.T) {

	data := "providers:\n  deepseek-v4-pro:cloud:\n      url: x\n    alias: y\n"

	parser := yaml.NewParser([]byte(data))
	_, err := parser.Root()

	if err == nil {
		t.Fatalf("expected a parse error, got nil")
	}

	if err.Error() == "missing document root node" {
		t.Fatalf("expected the real parse error, got generic missing root error")
	}

}
