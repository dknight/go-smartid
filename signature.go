package smartid

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
)

// Signature represents signature from session response.
type Signature struct {
	// Value is the base64 encoded string of signature.
	Value string `json:"value"`

	// Algorithm represents the algorithm used to encrypt the signature.
	Algorithm string `json:"algorithm"`
}

// IsValid checks the validity of signature.
func (sig Signature) IsValid(c Cert, h AuthHash) bool {
	decodedSig, err := base64.StdEncoding.DecodeString(sig.Value)
	if err != nil {
		return false
	}

	c.createX509CertIfNeeded()
	pubkey := c.x509Cert.PublicKey.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(pubkey, sig.resolveSignatureAlgo(), h, decodedSig)
	if err != nil {
		return false
	}
	return true
}

// Currently not in use, SK SmartID uses SHA512.
// resolveSignatureAlgo resolves encryption algorithm based on service response.
func (sig Signature) resolveSignatureAlgo() crypto.Hash {
	switch sig.Algorithm {
	case "sha256WithRSAEncryption":
		return crypto.SHA256
	case "sha384WithRSAEncryption":
		return crypto.SHA384
	case "sha512WithRSAEncryption":
		return crypto.SHA512
	default:
		return 0
	}
}
