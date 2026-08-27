package yaml

import "fmt"

type ParseError struct {
	Line    int
	Message string
}

func (err ParseError) Error() string {

	return fmt.Sprintf(
		"yaml parse error on line %d: %s",
		err.Line,
		err.Message,
	)

}
