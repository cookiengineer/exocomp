package findings

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"
import "path/filepath"

func writeFindingsReport(tool *Findings, report string) error {

	resolved, err0 := utils_sandbox.ResolvePath(tool.Playground, filepath.Join(".exocomp", "findings", "report.md"))

	if err0 == nil {

		if err1 := os.WriteFile(resolved, []byte(report), 0666); err1 == nil {
			return nil
		} else {
			return fmt.Errorf("Cannot write findings report.")
		}

	} else {
		return fmt.Errorf("Cannot resolve findings report.")
	}

}
