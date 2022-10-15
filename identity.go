package smartid

import (
	"crypto/x509/pkix"
	"strings"
)

// Identity represents simpler format of PKIX Subject.
type Identity struct {
	Country            string
	Organization       string
	OrganizationalUnit string
	Locality           string
	Province           string
	StreetAddress      string
	PostalCode         string
	SerialNumber       string
	CommonName         string
}

// newIdentity makes identity from PKIX Subject retrieved
// from certificate.
//
// This is very primitive PKIX Subject parsing, but should be enough
// for Smart-ID.
//
// TODO: refactor
func newIdentity(n *pkix.Name) *Identity {
	return &Identity{
		Country:            strings.Join(n.Country, ""),
		Organization:       strings.Join(n.Organization, ""),
		OrganizationalUnit: strings.Join(n.OrganizationalUnit, ""),
		Locality:           strings.Join(n.Locality, ""),
		Province:           strings.Join(n.Province, ""),
		StreetAddress:      strings.Join(n.StreetAddress, ""),
		PostalCode:         strings.Join(n.PostalCode, ""),
		SerialNumber:       n.SerialNumber,
		CommonName:         n.CommonName,
	}
}
