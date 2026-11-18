package actions

import "fmt"
import "os"
import "path/filepath"

func UsageAgent() {

	program := filepath.Base(os.Args[0])

	if program == "" {
		program = "exocomp-agent"
	}

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "Usage: %s [flags]\n", program)
	fmt.Fprint(os.Stdout, "\n")

	printUsageFlags()

	fmt.Fprint(os.Stdout, "Examples:\n")
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "  %s --role=coder --prompt=\"Implement the CalculateFibonacci method.\";\n", program)
	fmt.Fprintf(os.Stdout, "  %s --role=tester --prompt=\"Write the unit tests for CalculateFibonacci.\";\n", program)
	fmt.Fprint(os.Stdout, "\n")

	os.Exit(1)

}
