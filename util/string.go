package util

import (
	cryptoRand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"math/big"
	"math/rand"
	"regexp"
	"strings"
)

var sanitizeRegex = regexp.MustCompile(`[^a-zA-Z0-9-_.]+`)

func SanitizeName(n string, allowUnderscore, allowDot bool) string {
	if !allowUnderscore {
		n = strings.ReplaceAll(n, "_", "-")
	}

	if !allowDot {
		n = strings.ReplaceAll(n, ".", "-")
	}

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

var ansiControlRegex = regexp.MustCompile("[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]")

func StripAnsiControl(in string) string {
	return ansiControlRegex.ReplaceAllString(in, "")
}

const numChars = "0123456789"
const lowerChars = "abcdefghijklmnopqrstuvwxyz"
const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const specialChars = "!@#$%&*()-_=+[]{}<>:?"

func RandomStringCrypto(length int) string {
	return randomString(true, true, true, true, length, true)
}

func RandomString(length int) string {
	return randomString(true, true, true, true, length, false)
}

func RandomStringCryptoCustom(lower, upper, numeric, special bool, length int) string {
	return randomString(lower, upper, numeric, special, length, true)
}

func RandomStringCustom(lower, upper, numeric, special bool, length int) string {
	return randomString(lower, upper, numeric, special, length, false)
}

func charsPool(lower, upper, numeric, special bool) string {
	var chars string

	if lower {
		chars += lowerChars
	}

	if upper {
		chars += upperChars
	}

	if numeric {
		chars += numChars
	}

	if special {
		chars += specialChars
	}

	return chars
}

func randomString(lower, upper, numeric, special bool, length int, crypto bool) string {
	chars := charsPool(lower, upper, numeric, special)
	bytes := make([]byte, length)

	if crypto {
		err := randomCharsCrypto(bytes, chars)
		if err != nil {
			randomChars(bytes, chars)
		}
	} else {
		randomChars(bytes, chars)
	}

	return string(bytes)
}

func randomCharsCrypto(bytes []byte, chars string) error {
	setLen := big.NewInt(int64(len(chars)))

	for i := range bytes {
		idx, err := cryptoRand.Int(cryptoRand.Reader, setLen)
		if err != nil {
			return err
		}

		bytes[i] = chars[idx.Int64()]
	}

	return nil
}

func randomChars(bytes []byte, chars string) {
	setLen := len(chars)

	for i := range bytes {
		bytes[i] = chars[rand.Intn(setLen)]
	}
}
