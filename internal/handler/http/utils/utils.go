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

func ParseDiscordError(err error) string {
	switch {
	case strings.Contains(err.Error(), "10013"):
		return "User not found"
	case strings.Contains(err.Error(), "50035"):
		return "Incorrect user ID. Check ID for correctness"
	}

	return "Something went wrong"
}
