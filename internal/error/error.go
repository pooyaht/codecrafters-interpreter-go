package cerror

import "fmt"

type RuntimeError struct {
	Message string
	Line    int
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s \n[line %d]\n", e.Message, e.Line)
}
