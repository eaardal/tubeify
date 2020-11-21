package app

import (
	"fmt"
)

func logf(prefix string, msg string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s: %s\n", prefix, msg), args...)
}
