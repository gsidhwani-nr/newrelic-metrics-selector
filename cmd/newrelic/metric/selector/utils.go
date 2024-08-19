package main

import "strings"

func maskAPIKey(apiKey string) string {
	if len(apiKey) > 8 {
		return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
	}
	return strings.Repeat("*", len(apiKey))
}
