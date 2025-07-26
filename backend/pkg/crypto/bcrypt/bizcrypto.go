package bcrypt

// BCrypto crypto for biz
type BCrypto interface {
	Encrypt(s string) (string, error)
	Decrypt(payload string) (string, error)
}
