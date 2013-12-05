package curl

import (
	"fmt"
)

type curler_error struct {
	code   int
	reason string
}

func (this curler_error) Error() string {
	return fmt.Sprintf("%d: %v", this.code, this.reason)
}
