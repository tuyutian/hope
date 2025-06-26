package bcrypt

import (
	"fmt"
	"testing"
)

func TestCustomNumber(t *testing.T) {
	c := NewCustomNumber("qwertyuioplkjhgfdsazxcvbnmqazwsx", "stlDEFABCNOPyzghijQRSTUwxkVWXYZabcdefIJK67nopqr89LMmGH012345uv")

	plaintext := "20250427"
	ciphertext, err := c.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := c.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("plaintext: %s\nciphertext: %s\ndecrypted: %s\n%t\n",
		plaintext, ciphertext, decrypted, plaintext == decrypted)

	// 错误密文解密
	if len(ciphertext) > 1 {
		fmt.Println(c.Decrypt(ciphertext[:(len(ciphertext)-1)>>1]))
	}
	fmt.Println(c.Decrypt("abcdefghijklmnopabcdefghijklmnop"))
	fmt.Println(c.Decrypt(ciphertext + "abcdefghijklmnopabcdefghijklmnop"))
	fmt.Println(c.Decrypt("abcdefghijklmnopabcdefghijklmnop" + ciphertext))
}
