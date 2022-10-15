package smartid

import "fmt"

// Supported countries by Smart-ID.
const (
	CountryEE = "EE" // Estonia
	CountryLV = "LV" // Latvia
	CountryLT = "LT" // Lithuania
	CountryKZ = "KZ" // Kazakhstan
)

// People can be identified by their ETSI Natural Person Semantics Identifier
// specified in ETSI319412-1. Other way it might be passport number
// id card number, this depends on country or company internal politics.
const (
	// IdentifierTypePAS for identification based on passport number.
	IdentifierTypePAS = "PAS"

	// IdentifierTypeIDC for identification based on national identity
	// card number.
	IdentifierTypeIDC = "IDC"

	// IdentifierPNO for identification based on (national) personal
	// number (national civic registration number).
	IdentifierTypePNO = "PNO"
)

// NewSemanticIdentifier creates new semantic identifier as string.
func NewSemanticIdentifier(typ, country, id string) string {
	semid := SemanticIdentifier{
		Type:    typ,
		Country: country,
		ID:      id,
	}
	return semid.String()
}

// SemanticIdentifier is identifier to identify document type, country,
// and civic personal id.
//
// From official guide:
//
// Objects referenced by etsi/:semantics-identifier are persons identified
// by their ETSI Natural Person Sematics Identifier specified in
// ETSI319412-1. See more
// https://github.com/SK-EID/smart-id-documentation#2322-etsisemantics-identifier
type SemanticIdentifier struct {
	Type, Country, ID string
}

func (sd SemanticIdentifier) String() string {
	return fmt.Sprintf("%v%v-%v", sd.Type, sd.Country, sd.ID)
}
