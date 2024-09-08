package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Generates random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

// Generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for a := 0; a < n; a++ {
		idx := rand.Intn(k)
		sb.WriteByte(alphabet[idx])
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}