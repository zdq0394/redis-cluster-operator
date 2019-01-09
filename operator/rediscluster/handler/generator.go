package handler

import (
	"fmt"
)

func generateName(prefix, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}
