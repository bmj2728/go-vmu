package utils

import "al.essio.dev/pkg/shellescape"

func QuoteString(s string) string {
	return shellescape.Quote(s)
}
