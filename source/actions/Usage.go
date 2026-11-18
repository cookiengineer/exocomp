package actions

import "fmt"
import "os"
import "path/filepath"
import "slices"
import "strings"

func Usage(options_mode []string) {

	program := filepath.Base(os.Args[0])

	if program == "" {
		program = "exocomp"
	}

	example_mode := ""

	for _, mode := range []string{"web", "terminal"} {

		if slices.Contains(options_mode, mode) == true {
			example_mode = mode
			break
		}

	}

	if example_mode == "" && len(options_mode) > 0 {
		example_mode = options_mode[0]
	}

	default_mode := "unset"

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "Usage: %s <mode> [flags]\n", program)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "Arguments:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "  <mode> string          UI type\n")
	fmt.Fprint(os.Stdout, "                         (options: "+strings.Join(options_mode, ", ")+")\n")
	fmt.Fprint(os.Stdout, "                         (default: "+default_mode+")\n")
	fmt.Fprint(os.Stdout, "\n")

	printUsageFlags()

	if program != "" && example_mode != "" {

		fmt.Fprint(os.Stdout, "Examples:\n")
		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "  # single-agent mode\n")
		fmt.Fprintf(os.Stdout, "  %s %s --role=architect;\n", program, example_mode)
		fmt.Fprintf(os.Stdout, "  %s %s --role=architect --model=\"huihui_ai/Qwen3.8-abliterated:27b\" --temperature=\"0.5\";\n", program, example_mode)
		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "  # multi-agent mode\n")
		fmt.Fprintf(os.Stdout, "  %s %s --role=planner;\n", program, example_mode)
		fmt.Fprintf(os.Stdout, "  %s %s --role=planner --model=\"deepseek-v4-pro:cloud\" --temperature=\"0.7\";\n", program, example_mode)
		fmt.Fprint(os.Stdout, "\n")

	}

	os.Exit(1)

}
