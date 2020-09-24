package random

import(
	"bytes"
	"math/rand"
)

const (
	RANDOM_STRING = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	RANDOM_ROLE = [...]string{"user", "admin", "anonymous", "supplier", "operation", "system", "service"}
)

func RandomString(number int) string {
	var buffer bytes.Buffer

    for i := 0; i < number; i++ {
        buffer.WriteByte(RANDOM_STRING[rand.Intn(len(RANDOM_STRING))])
	}
	
	return buffer.String()
}

func RandomRole() string {
	return RANDOM_ROLE[rand.Intn(len(RANDOM_ROLE))]
}