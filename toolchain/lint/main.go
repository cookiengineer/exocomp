package main

import "encoding/json"
import "flag"
import "fmt"
import "os"
import "sort"
import "strings"

func main() {

	root := flag.String("root", "../source", "source tree to lint")
	is_including_tests := flag.Bool("tests", false, "also lint _test.go files")
	rule_filter := flag.String("rules", "", "comma separated rule IDs to run")
	list := flag.Bool("list", false, "list available rules")
	json_output := flag.Bool("json", false, "emit JSON diagnostics")

	flag.Parse()

	if *list == true {

		printRules()
		return

	}

	context, err0 := Load(*root, *is_including_tests)

	if err0 != nil {

		fmt.Fprintf(os.Stderr, "lint: %s\n", err0.Error())
		os.Exit(2)

	}

	for _, parse_error := range context.ParseErrors {
		fmt.Fprintf(os.Stderr, "lint: parse error: %s\n", parse_error)
	}

	selected := parseRuleFilter(*rule_filter)
	diagnostics := context.Run(selected)

	if *json_output == true {
		printJSON(diagnostics)
	} else {
		printText(diagnostics)
	}

	if len(diagnostics) > 0 {
		os.Exit(1)
	}

}

func parseRuleFilter(raw string) map[string]bool {

	if strings.TrimSpace(raw) == "" {
		return nil
	} else {

		result := make(map[string]bool)

		for _, entry := range strings.Split(raw, ",") {

			trimmed := strings.TrimSpace(entry)

			if trimmed != "" {
				result[trimmed] = true
			}

		}

		return result

	}

}

func printRules() {

	type entry struct {
		ID      string
		Summary string
	}

	entries := make([]entry, 0)

	for _, rule := range fileRules {
		entries = append(entries, entry{ID: rule.ID, Summary: rule.Summary})
	}

	for _, rule := range packageRules {
		entries = append(entries, entry{ID: rule.ID, Summary: rule.Summary})
	}

	sort.Slice(entries, func(a int, b int) bool {
		return entries[a].ID < entries[b].ID
	})

	for _, item := range entries {
		fmt.Printf("%-28s %s\n", item.ID, item.Summary)
	}

}

func printText(diagnostics []Diagnostic) {

	for _, diagnostic := range diagnostics {
		fmt.Println(diagnostic.String())
	}

	counts := make(map[string]int)

	for _, diagnostic := range diagnostics {
		counts[diagnostic.Rule]++
	}

	keys := make([]string, 0, len(counts))

	for key := range counts {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(diagnostics) > 0 {
		fmt.Fprintln(os.Stdout, "")
	}

	for _, key := range keys {
		fmt.Printf("%5d  %s\n", counts[key], key)
	}

	fmt.Printf("\n%d errors\n", len(diagnostics))

}

func printJSON(diagnostics []Diagnostic) {

	output := make([]jsonDiagnostic, 0, len(diagnostics))

	for _, diagnostic := range diagnostics {
		output = append(output, jsonDiagnostic{
			File:    diagnostic.Position.Filename,
			Line:    diagnostic.Position.Line,
			Column:  diagnostic.Position.Column,
			Rule:    diagnostic.Rule,
			Message: diagnostic.Message,
		})
	}

	encoded, err := json.MarshalIndent(output, "", "\t")

	if err != nil {

		fmt.Fprintf(os.Stderr, "lint: %s\n", err.Error())
		return

	}

	fmt.Println(string(encoded))

}
