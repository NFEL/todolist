package api_error

import (
	"fmt"
)

func UsernameExists(s string) error {
	return fmt.Errorf("username %s already exists.", s)
}
