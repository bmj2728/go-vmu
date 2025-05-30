package utils

import "fmt"

func QuoteString(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
