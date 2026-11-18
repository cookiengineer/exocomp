package findings

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "path/filepath"
import "strings"

func EvidenceArtifactPath(playground string, finding_id string, evidence_id string, filename string) (string, string, error) {

	filename = filepath.Base(strings.TrimSpace(filename))

	if filename == "" || filename == "." || filename == string(filepath.Separator) {
		filename = "evidence.txt"
	}

	filename = strings.Map(func(character rune) rune {

		if character >= 'a' && character <= 'z' {
			return character
		} else if character >= 'A' && character <= 'Z' {
			return character
		} else if character >= '0' && character <= '9' {
			return character
		} else if character == '.' || character == '-' || character == '_' {
			return character
		} else {
			return '-'
		}

	}, filename)

	relative := filepath.Join(".exocomp", "evidence", finding_id, evidence_id+"-"+filename)
	absolute, err0 := utils_sandbox.ResolvePath(playground, relative)

	if err0 == nil {

		display, err1 := utils_sandbox.SanitizePath(playground, absolute)

		if err1 == nil {
			return absolute, display, nil
		} else {
			return "", "", fmt.Errorf("Cannot resolve evidence artifact path.")
		}

	} else {
		return "", "", fmt.Errorf("Cannot resolve evidence artifact path.")
	}

}
