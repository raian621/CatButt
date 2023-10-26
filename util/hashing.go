package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func DefaultGenerateHash(text string) (string, error) {
	return GenerateHash(text, &ArgonParams{
		memory:      64 * 1024,
		iterations:  1,
		parallelism: 4,
		saltLength:  16,
		keyLength:   32,
	})
}

func GenerateHash(text string, p *ArgonParams) (string, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(text), salt, p.iterations, p.memory,
		p.parallelism, p.keyLength)

	return encodeArgonString(hash, salt, p), nil
}

func MatchHash(text, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(text), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (*ArgonParams, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("Invalid hash provided")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.New("Invalid hash provided: Error in field 2")
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("Incompatible Argon version")
	}

	p := &ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations,
		&p.parallelism)
	if err != nil {
		return nil, nil, nil, errors.New("Invalid hash provided: Error in field 3")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encodeArgonString(hash, salt []byte, p *ArgonParams) string {
	b64hash := base64.RawStdEncoding.Strict().EncodeToString(hash)
	b64salt := base64.RawStdEncoding.Strict().EncodeToString(salt)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory,
		p.iterations, p.parallelism, b64salt, b64hash)
}
