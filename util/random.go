package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijk"

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomInt generates random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates string with given length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner geneates random owner
func RandomOwner() string {
	return RandomString(6)
}

// RandomBalance generates random balance
func RandomBalance() int64 {
	return RandomInt(1, 1000)
}

// RandomCurrency generates random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "MMK"}
	return currencies[rand.Intn(3)]
}

// RandomTransactionType generates random transaction
func RandomTransactionType() string {
	txnTypes := []string{"Credit", "Debit"}
	return txnTypes[rand.Intn(2)]
}
