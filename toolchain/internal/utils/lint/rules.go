package lint

var fileRules []FileRule
var packageRules []PackageRule

func init() {

	fileRules = []FileRule{
		{"dep/tool-isolation", "tool packages must not import each other", checkToolIsolation},
		{"dep/tools-allowed", "tool packages may only import types, schemas, agents and utils", checkToolsAllowed},
		{"dep/types-no-capabilities", "types must not import utils/capabilities", checkTypesNoCapabilities},
		{"dep/utils-no-tools", "utils must not import tools", checkUtilsNoTools},
		{"branch/guard", "terminal guard clauses become positive-first branches", checkGuard},
		{"branch/guard-chain", "consecutive guards become one else-if chain", checkGuardChain},
		{"cond/bare", "conditions are compared explicitly", checkBareCondition},
		{"cond/not", "no ! negation in conditions", checkNot},
		{"errors/lowercase", "error messages start with a capital letter", checkErrorLowercase},
		{"errors/prefix", "error messages carry a component prefix", checkErrorPrefix},
		{"import/group", "project imports first, then stdlib, then third-party", checkImportGroup},
		{"import/order", "imports are sorted alphabetically within their group", checkImportOrder},
		{"import/paren", "one import per line, no parenthesized block", checkImportParen},
		{"import/testing-last", "testing is the last stdlib import in tests", checkTestingLast},
		{"layout/function-file", "a lone free function lives in a file named after it", checkFunctionFile},
		{"layout/init", "at most one init() per file", checkInit},
		{"layout/one-function", "every free function lives in its own file", checkOneFunction},
		{"layout/one-type", "one primary type per file", checkOneType},
		{"layout/tools-root", "source/tools contains only Toolset.go", checkToolsRoot},
		{"layout/type-file", "a lone unexported type lives in a file named after it", checkTypeFile},
		{"layout/utils-mirror", "utils/<name> must not mirror a top-level package", checkUtilsMirror},
		{"naming/boolean", "boolean values read as predicates", checkBooleanNaming},
		{"naming/file", "Go file names use the identifier spelling, no underscores", checkFileName},
		{"naming/from", "converters are named <Type>From<Source>", checkFromNaming},
		{"naming/func", "unexported functions use camelCase, never snake_case", checkFuncNaming},
		{"naming/type-initialism", "type names are spelled out, not collapsed to initials", checkTypeInitialism},
		{"style/any", "use any instead of interface{}", checkAny},
		{"style/blank-block", "blank lines around blocks and above switch cases", checkBlankBlock},
		{"style/gofmt", "file is formatted", checkGofmt},
		{"style/json-indent", "JSON indentation uses tabs", checkJSONIndent},
		{"style/local-case", "local variables use snake_case", checkLocalCase},
		{"test/file", "test files are named after a source file", checkTestFile},
		{"test/func", "test functions are named Test<Type>_<Method>", checkTestFunc},
	}

	packageRules = []PackageRule{
		{"layout/constructor-file", "New<Type> and Parse<Type> live in the type file", checkConstructorFile},
		{"layout/method-file", "methods live in the file of their receiver type", checkMethodFile},
		{"naming/receiver", "one receiver name per type", checkReceiverNaming},
	}

}
