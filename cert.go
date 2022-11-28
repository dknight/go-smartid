package smartid

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"time"
)

// Certificate levels. QUALIFIED is the more modern way to use it. Some
// accounts or services might support only non-qualified certificates,
// which are known as basic accounts (ADVANCED).
const (
	// CertLevelQualified QUALIFIED level of certificate.
	CertLevelQualified = "QUALIFIED" // recommended

	// CertLevelAdvanced ADVANCED level of certificate.
	CertLevelAdvanced = "ADVANCED"
)

const (
	certBegin = "-----BEGIN CERTIFICATE-----\n"
	certEnd   = "\n-----END CERTIFICATE-----"
)

// Cert represents certificate from session response.
type Cert struct {
	// Value is the base64 encoded string of certificate.
	Value string `json:"value"`

	// CertificateLevel is the level of the certificate:
	//	QUALIFIED
	//	ADVANCED
	CertificateLevel string `json:"certificateLevel"`

	// x509Cert is the X509 certificate.
	x509Cert *x509.Certificate
}

// IsExpired checks that certificate has expired.
func (c *Cert) IsExpired() bool {
	return time.Now().Before(c.x509Cert.NotBefore)
}

// IsNotActive checks that certificate is not yet active.
func (c *Cert) IsNotActive() bool {
	return time.Now().After(c.x509Cert.NotAfter)
}

// IsSameLevel checks that certificate is the same level as argument.
func (c *Cert) IsSameLevel(lvl string) bool {
	return c.CertificateLevel == lvl
}

// Verify certificate by file system paths.
// TODO: check more carefully.
func (c *Cert) Verify(paths []string) (bool, error) {
	roots := x509.NewCertPool()

	for _, path := range paths {
		cert, err := createCertFromPath(path)
		if err != nil {
			return false, err
		}
		roots.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := c.x509Cert.Verify(opts); err != nil {
		return false, err
	}
	return true, nil
}

// GetX509Cert returns X509 certificate from response.
func (c *Cert) GetX509Cert() *x509.Certificate {
	return c.x509Cert
}

// GetSubject get subject from certificate in PKIX format.
func (c *Cert) GetSubject() *pkix.Name {
	return &c.GetX509Cert().Subject
}

// GetIssuer get issuer from certificate in PKIX format.
func (c *Cert) GetIssuer() *pkix.Name {
	return &c.GetX509Cert().Issuer
}

// createCertFromPath certificate from given file system path.
func createCertFromPath(path string) (*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cert, err := createCertFromString(string(bs))
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// createCertFromString creates certificate from string.
func createCertFromString(s string) (*x509.Certificate, error) {
	block := decodePEMBlock([]byte(s))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// decodePEMBlock decodes PEM block.
func decodePEMBlock(pemData []byte) *pem.Block {
	block, _ := pem.Decode(pemData)
	return block
}

// createX509CertIfNeeded creates X509 certificate from response if not yet
// certificate exists.
func (c *Cert) createX509CertIfNeeded() {
	if c.GetX509Cert() == nil {
		cert, _ := createCertFromString(certBegin + c.Value + certEnd)
		c.x509Cert = cert
	}
}
