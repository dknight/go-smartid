package smartid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// Supported hashing algorithms.
const (
	// SHA256 algorithm for encryption.
	SHA256 = "SHA256"

	// SHA384 algorithm for encryption.
	SHA384 = "SHA384"

	// SHA512 algorithm for encryption.
	SHA512 = "SHA512"
)

var randByteGenerator = generateRandomBytes

// nbytes how many random bytes to generate.
const nbytes int = 64

// AuthHash contains the hash sum for authentication and signing requests
// and the calculation of the verification code.
type AuthHash []byte

// GenerateAuthHash generates a new random hashe.
func GenerateAuthHash(algo string) AuthHash {
	bs := randByteGenerator(nbytes)
	var sum []byte
	var sumSha256 [sha256.Size]byte
	var sumSha384 [sha512.Size384]byte
	var sumSha512 [sha512.Size]byte
	switch algo {
	case SHA256:
		sumSha256 = sha256.Sum256(bs)
		sum = sumSha256[:]
	case SHA384:
		sumSha384 = sha512.Sum384(bs)
		sum = sumSha384[:]
	case SHA512:
		sumSha512 = sha512.Sum512(bs)
		sum = sumSha512[:]
	}
	return AuthHash(sum[:])
}

// EncodeBase64 encodes hash sum to base64. Normally you should not convert
// hash to base64 manually. Library does it automatically on the request.
func (h AuthHash) EncodeBase64() []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(h)))
	base64.StdEncoding.Encode(dst, h)
	return dst
}

// ToBase64String coverts AuthHash to base64 string.
func (h AuthHash) ToBase64String() string {
	return string(h.EncodeBase64())
}

// CalculateVerificationCode computes the verification 4-digit verification
// that you can show it to use.
// Examine how VerificationCode is computed.
// https://github.com/SK-EID/smart-id-documentation#23122-computing-the-verification-code
func (h AuthHash) CalculateVerificationCode() string {
	var x uint16
	h256 := sha256.New()
	h256.Write(h)
	sum := h256.Sum(nil)
	buf := bytes.NewReader(sum[len(sum)-2:])
	binary.Read(buf, binary.BigEndian, &x)
	return fmt.Sprintf("%04d", x%10000)
}

// generateRandomBytes create N number bytes.
func generateRandomBytes(n int) []byte {
	bs := make([]byte, n)
	rand.Read(bs)
	return bs
}
