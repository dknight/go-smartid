package smartid

import (
	"crypto"
	"io/ioutil"
	"testing"
)

func TestSignature_resolveSignatureAlgo(t *testing.T) {
	sig1 := Signature{
		Algorithm: "sha256WithRSAEncryption",
	}
	if sig1.resolveSignatureAlgo() != crypto.SHA256 {
		t.Error("expected", crypto.SHA256, "got", sig1.resolveSignatureAlgo())
	}

	sig2 := Signature{
		Algorithm: "sha384WithRSAEncryption",
	}
	if sig2.resolveSignatureAlgo() != crypto.SHA384 {
		t.Error("expected", crypto.SHA384, "got", sig2.resolveSignatureAlgo())
	}

	sig3 := Signature{
		Algorithm: "sha512WithRSAEncryption",
	}
	if sig3.resolveSignatureAlgo() != crypto.SHA512 {
		t.Error("expected", crypto.SHA512, "got", sig3.resolveSignatureAlgo())
	}

	sig4 := Signature{
		Algorithm: "sha224WithRSAEncryption",
	}
	if sig4.resolveSignatureAlgo() != 0 {
		t.Error("expected", 0, "got", sig4.resolveSignatureAlgo())
	}
}

func TestSignature_IsValid(t *testing.T) {
	certValue, _ := ioutil.ReadFile("./files/test.crt")
	hash := GenerateAuthHash(SHA512)
	cert := Cert{
		CertificateLevel: CertLevelQualified,
		Value:            string(certValue),
	}
	cert.createX509CertIfNeeded()
	sig := Signature{
		Algorithm: "sha512WithRSAEncryption",
		Value:     "foobar",
	}
	if sig.IsValid(cert, hash) {
		t.Error("Test fail signature should be invalid")
	}

	sig = Signature{
		Algorithm: "sha512WithRSAEncryption",
		Value:     "Zm9vYmFy",
	}
	if sig.IsValid(cert, hash) {
		t.Error("Test fail signature should be invalid")
	}
}
