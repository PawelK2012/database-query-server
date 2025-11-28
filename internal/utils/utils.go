package utils

import (
	"fmt"
	"strings"
)

// checkFirstWord checks if the first word in a query is "SELECT"
func CheckFirstWord(input string) error {
	words := strings.Fields(input)
	if len(words) > 0 && words[0] != "SELECT" {
		// for security reasons we don’t include too many details in the error message
		return fmt.Errorf("query is not valid")
	}
	return nil
}
