package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string of length equal n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a random owner name
func RandomOwner() string {
	return RandomString(15)
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate a random currency code
func RandomCurrency() string {
	currency := []string{USD, RUB, EUR}
	k := len(currency)

	return currency[rand.Intn(k)]
}

// RandomEmail generate a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@%s.%s", RandomString(10), RandomString(5), RandomString(2))
}
