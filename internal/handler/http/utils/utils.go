package utils

import (
	"strconv"
	"strings"
)

func ParseDelay(value string) (int, error) {
	value = strings.TrimSpace(value)
	switch value {
	case "":
		return 3, nil
	default:
		return strconv.Atoi(value)
	}
}

func ParseError(err error) (string, string) {
	switch {
	case strings.Contains(err.Error(), "10013"):
		return "Not Found", "User not found"
	case strings.Contains(err.Error(), "50035"):
		return "Incorrect ID", "Check your ID for correctness"
	}

	return "Unknown Error", "Something went wrong"
}
