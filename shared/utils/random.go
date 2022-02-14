package utils

import (
	cryptoRand "crypto/rand"
	"io"
	"math/rand"

	"time"
)

// PickRandomString picks a random string
func PickRandomString(str ...string) string {

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(str)
	return str[n]
}

// PickRandomInt picks a random string
func PickRandomInt(len int) int {

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len
	return n
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// GenerateOTPCode generates random number for otp code
func GenerateOTPCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(cryptoRand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// GenerateRandomToken does generate a random string
func GenerateRandomToken(max int) string {
	b := make([]byte, max)
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	rand.Seed(time.Now().Unix())
	for i := 0; i < max; i++ {
		b[i] = table[rand.Int()%10]
	}

	return string(b)
}
