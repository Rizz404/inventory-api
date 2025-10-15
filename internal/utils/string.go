package utils

import "regexp"

// * Helper function to replace all occurrences of a string (case-insensitive)
func ReplaceAllCaseInsensitive(s, old, new string) string {
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(old))
	return re.ReplaceAllString(s, new)
}
