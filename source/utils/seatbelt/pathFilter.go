package seatbelt

import "fmt"

func pathFilter(resolved_path resolvedPath) string {

	if resolved_path.exists == true {

		matcher := "subpath"

		if resolved_path.isDir == false {
			matcher = "literal"
		}

		if resolved_path.original == resolved_path.resolved {
			return fmt.Sprintf("(%s %s)", matcher, quoteSBPL(resolved_path.original))
		} else {
			return fmt.Sprintf(
				"(require-any (%s %s) (%s %s))",
				matcher, quoteSBPL(resolved_path.original),
				matcher, quoteSBPL(resolved_path.resolved),
			)
		}

	} else {

		if resolved_path.original == resolved_path.resolved {
			return fmt.Sprintf(
				"(require-any (literal %s) (subpath %s))",
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.original),
			)
		} else {
			return fmt.Sprintf(
				"(require-any (literal %s) (subpath %s) (literal %s) (subpath %s))",
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.resolved),
				quoteSBPL(resolved_path.resolved),
			)
		}

	}

}
