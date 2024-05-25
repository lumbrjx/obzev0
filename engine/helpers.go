package main

import (
	"strings"
)

func IsHTTPMethod(method string) bool {
	validMethods := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
		"HEAD",
		"OPTIONS",
		"CONNECT",
		"TRACE",
	}
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}

func ExtractURL(request string) (string, string, bool) {
	lines := strings.Split(request, "\n")
	if len(lines) == 0 {
		return "", "", false
	}

	components := strings.Split(lines[0], " ")
	if len(components) < 3 {
		return "", "", false
	}

	method := components[0]
	if !IsHTTPMethod(method) {
		return "", "", false
	}
	url := components[1]

	return method, url, true
}
