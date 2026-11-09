package seatbelt

import "fmt"
import "strings"

func validatePath(path string) error {

	if strings.ContainsAny(path, "\x00\r\n") == true {
		return fmt.Errorf("seatbelt.validatePath: Invalid path.")
	} else {
		return nil
	}

}
