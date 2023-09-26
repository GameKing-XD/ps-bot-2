package randstr

import (
	"math/rand"
	"strings"
)

const Specialchars = "()[]{}-=+/,.?"
const Numbers = "0123456789"
const Lowercase = "abcdefghijklmnopqrstuvwxyz"
const Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Letters = Uppercase + Lowercase
const Validchars = Uppercase + Lowercase + Numbers + Specialchars

func Randstr(set string, length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		random := rand.Intn(len(set))
		sb.WriteByte(set[random])
	}

	return sb.String()
}

func Concat(strs ...string) string {
	return strings.Join(strs, "")
}
