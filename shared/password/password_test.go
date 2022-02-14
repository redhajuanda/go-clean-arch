package password

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	hashedPwd, _ := HashAndSalt([]byte("password1234"))
	fmt.Println(hashedPwd)
	fmt.Println(ComparePasswords(hashedPwd, []byte("password1234")))
}
