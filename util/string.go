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

var (
	validEnvVarRegex = regexp.MustCompile(`[^A-Za-z0-9_]`)
)

func SanitizeEnvVar(in string) string {
	return validEnvVarRegex.ReplaceAllString(in, "_")
}

var ansiRegex = regexp.MustCompile("[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]")

func StripAnsi(in string) string {
	return ansiRegex.ReplaceAllString(in, "")
}
