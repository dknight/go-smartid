package smartid

import "testing"

type semIDTestPair struct {
	identifier SemanticIdentifier
	result     string
}

func TestSemanticIdentifier_String(t *testing.T) {
	testdata := []semIDTestPair{
		semIDTestPair{
			SemanticIdentifier{
				Type:    IdentifierTypePNO,
				Country: CountryEE,
				ID:      "30303039914",
			},
			"PNOEE-30303039914",
		},
		semIDTestPair{
			SemanticIdentifier{
				Type:    IdentifierTypeIDC,
				Country: CountryLV,
				ID:      "030303-10012",
			},
			"IDCLV-030303-10012",
		},
		semIDTestPair{
			SemanticIdentifier{
				Type:    IdentifierTypePAS,
				Country: CountryKZ,
				ID:      "1234567890",
			},
			"PASKZ-1234567890",
		},
	}

	for _, pair := range testdata {
		if pair.identifier.String() != pair.result {
			t.Error("expected", pair.result, "got", pair.identifier.String())
		}
	}
}
