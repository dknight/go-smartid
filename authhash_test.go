package smartid

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"reflect"
	"testing"
)

var randHashTest = []byte{65, 65, 65, 65, 65, 65, 65, 65, 65, 65}

func TestGenerateAuthHash(t *testing.T) {
	randByteGenerator = func(n int) []byte {
		return randHashTest
	}
	var hash AuthHash
	hash = GenerateAuthHash(SHA256)
	if len(hash) != sha256.Size {
		t.Error("expected", sha256.Size, "got", len(hash))
	}
	hash = GenerateAuthHash(SHA384)
	if len(hash) != sha512.Size384 {
		t.Error("expected", sha512.Size384, "got", len(hash))
	}
	hash = GenerateAuthHash(SHA512)
	if len(hash) != sha512.Size {
		t.Error("expected", sha512.Size, "got", len(hash))
	}
}

func TestHashSum_EncodeBase64(t *testing.T) {
	hash := GenerateAuthHash(SHA512)
	hashEncoded := hash.EncodeBase64()
	result := make([]byte, base64.StdEncoding.EncodedLen(len(hash)))
	base64.StdEncoding.Encode(result, hash)
	if !reflect.DeepEqual(hashEncoded, result) {
		t.Error("expected", result, "got", hashEncoded)
	}
}

func TestHashSum_ToBase64String(t *testing.T) {
	hash := GenerateAuthHash(SHA512)
	hashEncoded := hash.ToBase64String()
	result := make([]byte, base64.StdEncoding.EncodedLen(len(hash)))
	base64.StdEncoding.Encode(result, hash)
	if string(hashEncoded) != string(result) {
		t.Error("expected", string(result), "got", string(hashEncoded))
	}
}

func TestAuthHash_CalculateVerificationCode(t *testing.T) {
	vc := GenerateAuthHash(SHA512).CalculateVerificationCode()
	result := "3174"
	if vc != result {
		t.Error("expected", result, "got", vc)
	}
}
