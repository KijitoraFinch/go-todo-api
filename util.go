package main

import (
	"golang.org/x/exp/slices"
	"net/url"
)

func HasCommonElement(a, b []string) bool {
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}
func ShouldSkip(query url.Values, key string, taskValue string) bool {
	if len(query[key]) > 0 && query[key][0] != "false" {
		for _, value := range query[key] {
			if value == taskValue {
				return false
			}
		}
		return true
	}
	return false
}
