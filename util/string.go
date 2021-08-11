package util

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
)

var sanitizeRegex = regexp.MustCompile("[^a-zA-Z0-9-]+")

func SanitizeName(n string) string {
	n = strings.ReplaceAll(n, "_", "-")
	n = strings.ReplaceAll(n, ".", "-")

	return sanitizeRegex.ReplaceAllString(n, "")
}

func SHAString(n string) string {
	h := sha1.New()
	h.Write([]byte(n))
	sha := h.Sum(nil)

	return hex.EncodeToString(sha)
}

func LimitString(n string, lim int) string {
	if len(n) > lim {
		return n[:lim]
	}

	return n
}

func StringSliceContains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}
