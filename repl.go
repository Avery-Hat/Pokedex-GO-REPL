// repl.go
package main

import "strings"

// cleanInput lowercases, trims, and splits on any whitespace.
func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{}
	}
	lower := strings.ToLower(trimmed)
	// Fields splits on any run of whitespace (spaces, tabs, newlines).
	return strings.Fields(lower)
}
