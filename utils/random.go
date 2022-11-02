package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alpha = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alpha)

	for i := 0; i < n; i++ {
		c := alpha[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomCurrency() string {
	curs := []string{"USD", "EUR", "CAD"}
	return curs[rand.Intn(len(curs))]
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(10, 1000)
}
