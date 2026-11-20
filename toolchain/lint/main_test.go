package main

import "os"
import "path/filepath"
import "testing"

func writeTree(t *testing.T, files map[string]string) string {

	t.Helper()

	root := t.TempDir()

	for name, content := range files {

		full := filepath.Join(root, name)

		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("Cannot create folder for %s: %v", name, err)
		}

		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("Cannot write %s: %v", name, err)
		}

	}

	return root

}

func lintRule(t *testing.T, rule string, files map[string]string) []Diagnostic {

	t.Helper()

	files["go.mod"] = "module exocomp\n\ngo 1.27\n"

	root := writeTree(t, files)

	context, err := Load(root, true)

	if err != nil {
		t.Fatalf("Cannot load tree: %v", err)
	}

	result := make([]Diagnostic, 0)

	for _, diagnostic := range context.Run(nil) {

		if diagnostic.Rule == rule {
			result = append(result, diagnostic)
		}

	}

	return result

}

type lintCase struct {
	Name  string
	Rule  string
	Files map[string]string
	Want  int
}

func TestRules(t *testing.T) {

	cases := []lintCase{
		{
			Name: "method-file/missing",
			Rule: "layout/method-file",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\ntype A struct{}\n",
				"pkg/b.go": "package pkg\n\nfunc (a *A) Foo() {}\n",
			},
			Want: 1,
		},
		{
			Name: "method-file/ok",
			Rule: "layout/method-file",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\ntype A struct{}\n\nfunc (a *A) Foo() {}\n",
			},
			Want: 0,
		},
		{
			Name: "one-type/missing",
			Rule: "layout/one-type",
			Files: map[string]string{
				"pkg/types.go": "package pkg\n\ntype A struct{}\n\ntype B struct{}\n",
			},
			Want: 1,
		},
		{
			Name: "one-type/ok",
			Rule: "layout/one-type",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\ntype A struct{}\n",
			},
			Want: 0,
		},
		{
			Name: "type-file/missing",
			Rule: "layout/type-file",
			Files: map[string]string{
				"pkg/foo.go": "package pkg\n\ntype bar struct{}\n",
			},
			Want: 1,
		},
		{
			Name: "type-file/ok",
			Rule: "layout/type-file",
			Files: map[string]string{
				"pkg/bar.go": "package pkg\n\ntype bar struct{}\n",
			},
			Want: 0,
		},
		{
			Name: "one-function/missing",
			Rule: "layout/one-function",
			Files: map[string]string{
				"pkg/Multi.go": "package pkg\n\nfunc Multi() {}\n\nfunc Other() {}\n",
			},
			Want: 1,
		},
		{
			Name: "one-function/ok",
			Rule: "layout/one-function",
			Files: map[string]string{
				"pkg/foo.go": "package pkg\n\nfunc Foo() {}\n",
			},
			Want: 0,
		},
		{
			Name: "function-file/missing",
			Rule: "layout/function-file",
			Files: map[string]string{
				"pkg/foo.go": "package pkg\n\nfunc Bar() {}\n",
			},
			Want: 1,
		},
		{
			Name: "function-file/ok",
			Rule: "layout/function-file",
			Files: map[string]string{
				"pkg/Bar.go": "package pkg\n\nfunc Bar() {}\n",
			},
			Want: 0,
		},
		{
			Name: "constructor-file/missing",
			Rule: "layout/constructor-file",
			Files: map[string]string{
				"pkg/Thing.go": "package pkg\n\ntype Thing struct{}\n",
				"pkg/other.go": "package pkg\n\nfunc NewThing() *Thing { return nil }\n",
			},
			Want: 1,
		},
		{
			Name: "constructor-file/ok",
			Rule: "layout/constructor-file",
			Files: map[string]string{
				"pkg/Thing.go": "package pkg\n\ntype Thing struct{}\n\nfunc NewThing() *Thing { return nil }\n",
			},
			Want: 0,
		},
		{
			Name: "tools-root/missing",
			Rule: "layout/tools-root",
			Files: map[string]string{
				"tools/other.go": "package tools\n",
			},
			Want: 1,
		},
		{
			Name: "tools-root/ok",
			Rule: "layout/tools-root",
			Files: map[string]string{
				"tools/Toolset.go": "package tools\n",
			},
			Want: 0,
		},
		{
			Name: "init/missing",
			Rule: "layout/init",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc init() {}\n\nfunc init() {}\n",
			},
			Want: 1,
		},
		{
			Name: "utils-mirror/missing",
			Rule: "layout/utils-mirror",
			Files: map[string]string{
				"types/a.go":       "package types\n",
				"utils/types/b.go": "package types\n",
			},
			Want: 1,
		},
		{
			Name: "tool-isolation/missing",
			Rule: "dep/tool-isolation",
			Files: map[string]string{
				"tools/a/a.go": "package a\n\nimport \"exocomp/tools/b\"\n",
				"tools/b/b.go": "package b\n",
			},
			Want: 1,
		},
		{
			Name: "tools-allowed/missing",
			Rule: "dep/tools-allowed",
			Files: map[string]string{
				"tools/a/a.go": "package a\n\nimport \"exocomp/engine\"\n",
			},
			Want: 1,
		},
		{
			Name: "tools-allowed/ok",
			Rule: "dep/tools-allowed",
			Files: map[string]string{
				"tools/a/a.go": "package a\n\nimport \"exocomp/types\"\n",
			},
			Want: 0,
		},
		{
			Name: "types-no-capabilities/missing",
			Rule: "dep/types-no-capabilities",
			Files: map[string]string{
				"types/a.go": "package types\n\nimport \"exocomp/utils/capabilities\"\n",
			},
			Want: 1,
		},
		{
			Name: "utils-no-tools/missing",
			Rule: "dep/utils-no-tools",
			Files: map[string]string{
				"utils/a/a.go": "package a\n\nimport \"exocomp/tools\"\n",
			},
			Want: 1,
		},
		{
			Name: "import-paren/missing",
			Rule: "import/paren",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport (\n\t\"fmt\"\n)\n",
			},
			Want: 1,
		},
		{
			Name: "import-group/blank-between-project-and-stdlib",
			Rule: "import/group",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"exocomp/types\"\n\nimport \"fmt\"\n",
			},
			Want: 1,
		},
		{
			Name: "import-group/ok",
			Rule: "import/group",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"exocomp/types\"\nimport \"fmt\"\n\nimport \"golang.org/x/net/html\"\n",
			},
			Want: 0,
		},
		{
			Name: "import-order/missing",
			Rule: "import/order",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"net/url\"\nimport \"fmt\"\n",
			},
			Want: 1,
		},
		{
			Name: "testing-last/missing",
			Rule: "import/testing-last",
			Files: map[string]string{
				"pkg/a_test.go": "package pkg\n\nimport \"testing\"\nimport \"strings\"\n\nfunc TestThing_Do(t *testing.T) {}\n",
				"pkg/a.go":      "package pkg\n",
			},
			Want: 1,
		},
		{
			Name: "naming-file/missing",
			Rule: "naming/file",
			Files: map[string]string{
				"pkg/my_file.go": "package pkg\n",
			},
			Want: 1,
		},
		{
			Name: "naming-file/ok-platform",
			Rule: "naming/file",
			Files: map[string]string{
				"pkg/check_linux.go": "package pkg\n",
			},
			Want: 0,
		},
		{
			Name: "naming-func/missing",
			Rule: "naming/func",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc do_thing() {}\n",
			},
			Want: 1,
		},
		{
			Name: "naming-from/missing",
			Rule: "naming/from",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc FromInt(value int) int { return value }\n",
			},
			Want: 1,
		},
		{
			Name: "naming-type-initialism/missing",
			Rule: "naming/type-initialism",
			Files: map[string]string{
				"pkg/rgb.go": "package pkg\n\ntype RGB struct{}\n",
			},
			Want: 1,
		},
		{
			Name: "naming-boolean/missing",
			Rule: "naming/boolean",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(debug bool) {}\n",
			},
			Want: 1,
		},
		{
			Name: "naming-boolean/ok",
			Rule: "naming/boolean",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(is_debug bool) {}\n",
			},
			Want: 0,
		},
		{
			Name: "naming-receiver/missing",
			Rule: "naming/receiver",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\ntype A struct{}\n\nfunc (a *A) One() {}\n\nfunc (other *A) Two() {}\n",
			},
			Want: 1,
		},
		{
			Name: "branch-guard/missing",
			Rule: "branch/guard",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(err error) error {\n\n\tif err != nil {\n\t\treturn err\n\t}\n\n\treturn nil\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "branch-guard/ok-nil-receiver",
			Rule: "branch/guard",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\ntype A struct{}\n\nfunc (a *A) run() int {\n\n\tif a == nil {\n\t\treturn 0\n\t}\n\n\treturn 1\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "branch-guard-chain/missing",
			Rule: "branch/guard-chain",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(a error) error {\n\n\tif a != nil {\n\t\treturn a\n\t}\n\n\tif a == nil {\n\t\treturn nil\n\t}\n\n\treturn nil\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "cond-not/missing",
			Rule: "cond/not",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(ok bool) {\n\n\tif !ok {\n\t}\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "cond-bare/missing",
			Rule: "cond/bare",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run(ok bool) {\n\n\tif ok {\n\t}\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "errors-prefix/missing",
			Rule: "errors/prefix",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"fmt\"\n\nfunc run() error {\n\n\treturn fmt.Errorf(\"nope\")\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "errors-prefix/tool-helper-inherits",
			Rule: "errors/prefix",
			Files: map[string]string{
				"tools/x/helper.go": "package x\n\nimport \"fmt\"\n\nfunc helper() error {\n\n\treturn fmt.Errorf(\"Nope.\")\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "errors-prefix/tool-exported-needs-prefix",
			Rule: "errors/prefix",
			Files: map[string]string{
				"tools/x/x.go": "package x\n\nimport \"fmt\"\n\nfunc Run() error {\n\n\treturn fmt.Errorf(\"Nope.\")\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "errors-prefix/placeholders",
			Rule: "errors/prefix",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"fmt\"\n\nfunc run(method string) error {\n\n\treturn fmt.Errorf(\"pkg.%s: Nope.\", method)\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "errors-prefix/ok",
			Rule: "errors/prefix",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"fmt\"\n\nfunc run() error {\n\n\treturn fmt.Errorf(\"pkg.Run: Nope.\")\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "errors-lowercase/missing",
			Rule: "errors/lowercase",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"fmt\"\n\nfunc run() error {\n\n\treturn fmt.Errorf(\"pkg.Run: nope\")\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "style-any/missing",
			Rule: "style/any",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nvar value interface{}\n",
			},
			Want: 1,
		},
		{
			Name: "style-any/ok",
			Rule: "style/any",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nvar value any\n",
			},
			Want: 0,
		},
		{
			Name: "style-blank-block/missing",
			Rule: "style/blank-block",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run() {\n\tone := 1\n\ttwo := 2\n\n\t_ = one\n\t_ = two\n}\n",
			},
			Want: 2,
		},
		{
			Name: "style-blank-block/ok",
			Rule: "style/blank-block",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run() {\n\n\tone := 1\n\ttwo := 2\n\n\t_ = one\n\t_ = two\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "style-gofmt/missing",
			Rule: "style/gofmt",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run() {\n  value := 1\n  _ = value\n}\n",
			},
			Want: 1,
		},
		{
			Name: "style-json-indent/missing",
			Rule: "style/json-indent",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nimport \"encoding/json\"\n\nfunc run(value any) ([]byte, error) {\n\n\treturn json.MarshalIndent(value, \"\", \"  \")\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "style-local-case/missing",
			Rule: "style/local-case",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run() {\n\n\tonSwitch := 1\n\t_ = onSwitch\n\n}\n",
			},
			Want: 1,
		},
		{
			Name: "style-local-case/ok-acronym",
			Rule: "style/local-case",
			Files: map[string]string{
				"pkg/a.go": "package pkg\n\nfunc run() {\n\n\tbaseURL := \"x\"\n\t_ = baseURL\n\n}\n",
			},
			Want: 0,
		},
		{
			Name: "test-file/missing",
			Rule: "test/file",
			Files: map[string]string{
				"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n",
			},
			Want: 1,
		},
		{
			Name: "test-file/ok",
			Rule: "test/file",
			Files: map[string]string{
				"pkg/foo.go":      "package pkg\n",
				"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n",
			},
			Want: 0,
		},
		{
			Name: "test-func/missing",
			Rule: "test/func",
			Files: map[string]string{
				"pkg/foo.go":      "package pkg\n",
				"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc Test_bad(t *testing.T) {}\n",
			},
			Want: 1,
		},
	}

	for _, test_case := range cases {

		t.Run(test_case.Name, func(t *testing.T) {

			diagnostics := lintRule(t, test_case.Rule, test_case.Files)

			if len(diagnostics) != test_case.Want {
				t.Errorf("Expected %d %s diagnostics, got %d", test_case.Want, test_case.Rule, len(diagnostics))

				for _, diagnostic := range diagnostics {
					t.Logf("  %s", diagnostic.String())
				}
			}

		})

	}

}

func TestSuppression(t *testing.T) {

	diagnostics := lintRule(t, "style/any", map[string]string{
		"pkg/a.go": "package pkg\n\nvar value interface{} //nolint:style/any\n",
	})

	if len(diagnostics) != 0 {
		t.Errorf("Expected the nolint comment to suppress the diagnostic, got %d", len(diagnostics))
	}

}
