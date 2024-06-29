package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"io"

	"github.com/Stuhub-io/core/ports"
	"golang.org/x/crypto/scrypt"
)

type scryptImpl struct {
	saltLength uint32
	n          int
	r          int
	p          int
	keyLength  int
	secretKey  []byte
}

func newScryptDefault(secretKey []byte) *scryptImpl {
	return &scryptImpl{
		saltLength: 8,
		n:          1 << 15,
		r:          8,
		p:          1,
		keyLength:  32,
		secretKey:  secretKey,
	}
}

// NewScrypt creates a new scrypt password helper.
func NewScrypt(secretKey []byte) ports.Hasher {
	return newScryptDefault(secretKey)
}

func (h scryptImpl) GenerateSalt() string {
	var salt = make([]byte, h.saltLength)

	_, err := io.ReadFull(rand.Reader, salt)

	if err != nil {
		panic(err)
	}

	return base64.RawStdEncoding.EncodeToString(salt)
}

func (h scryptImpl) Hash(password, salt string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	passwordWithKey := append(passwordBytes, h.secretKey...)

	hash, err := scrypt.Key(passwordWithKey, saltBytes, h.n, h.r, h.p, h.keyLength)
	if err != nil {
		return "", err
	}

	hashedStr := base64.RawStdEncoding.EncodeToString(hash)

	return hashedStr, nil
}

func (h scryptImpl) Compare(password, hashedPassword, salt string) bool {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	passwordWithKey := append([]byte(password), h.secretKey...)

	otherHash, err := scrypt.Key(passwordWithKey, saltBytes, h.n, h.r, h.p, h.keyLength)
	if err != nil {
		return false
	}

	hashedPasswordBytes, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	if subtle.ConstantTimeCompare(hashedPasswordBytes, otherHash) == 1 {
		return true
	}

	return false
}
