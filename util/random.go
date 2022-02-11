package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwyxz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt generate a random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // ritorna numero tra min e max
}

//RandomString
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//RandomOwner
func RandomOwner() string {
	return RandomString(6)
}

//RandomMoney
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency ritorna una currency casuale tra quelle definite
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, PND}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
