package utils

import (
	"strings"
)

func NormalizeName(fullName string) (firstName string, lastName string) {

	fullName = strings.TrimSpace(fullName)
	items := strings.Split(fullName, " ")

	firstName = items[0]
	if len(items) > 1 {
		lastName = strings.Join(items[1:], " ")
		lastName = strings.TrimSpace(lastName)
	}

	return
}
