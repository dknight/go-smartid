package smartid

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

type ClientTestTable map[string]struct {
	request AuthRequest
	result  ClientTestResult
}
type ClientTestResult struct {
	Identity
	Response
}

const demoPartyUUID = "00000000-0000-0000-0000-000000000000"
const demoPartyName = "DEMO"

var client = NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 10000)
var clientTestTableAuth = ClientTestTable{
	"client_ee_id_ok": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30303039914"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,OK",
				SerialNumber: "PNOEE-30303039914",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
	"client_ee_id_other": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA256),
			HashType:         SHA256,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30303039816"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,MULTIPLE OK",
				SerialNumber: "PNOEE-30303039816",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
	"client_ee_id_no_numbers": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA384),
			HashType:         SHA384,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30303039903"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,QUALIFIED OK,PNOEE-30303039903",
				SerialNumber: "PNOEE-30303039903",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
	"client_lt_id_other": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryLT,
				"49912318881"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "LT",
				CommonName:   "TESTNUMBER, OK",
				SerialNumber: "PNOLT-49912318881",
			},
			Response{
				Code:    471,
				Message: "No suitable account of requested type found, but user has some other accounts.",
			},
		},
	},
	"client_lv_id_other": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryLV,
				"311299-18886"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "LV",
				CommonName:   "TESTNUMBER,OK",
				SerialNumber: "PNOLV-311299-18886",
			},
			Response{
				Code:    471,
				Message: "No suitable account of requested type found, but user has some other accounts.",
			},
		},
	},
	"client_ee_doc_new_cert": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier:       "PNOEE-39912319997-AAAA-Q",
			AuthType:         AuthTypeDocument,
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,BOD",
				SerialNumber: "PNOEE-39912319997",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
	"client_ee_doc_ageu18": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier:       "PNOEE-50701019992-9ZN6-Q",
			AuthType:         AuthTypeDocument,
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,MINOR",
				SerialNumber: "PNOEE-50701019992",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
	"client_ee_id_refuse1": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039928"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultUserRefusedDisplayTextAndPIN,
			},
		},
	},
	"client_ee_id_refuse2": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039939"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultUserRefusedVCChoice,
			},
		},
	},
	"client_ee_id_refuse3": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039946"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultUserRefusedConfirmationMessage,
			},
		},
	},
	"client_ee_id_refuse4": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039950"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultUserRefusedConfirmationMessageWithVCChoice,
			},
		},
	},
	"client_ee_id_refuse5": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039961"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultUserRefusedCertChoice,
			},
		},
	},
	"client_ee_id_wrong_vc": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039972"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultWrongVC,
			},
		},
	},
	"client_ee_id_timeout": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30403039983"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultTimeout,
			},
		},
	},
	"client_ee_id_not_found": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"01234567891"),
		},
		result: ClientTestResult{
			Identity{},
			Response{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			},
		},
	},
}

func TestAuthenticate(t *testing.T) {
	for key, test := range clientTestTableAuth {
		fmt.Println("Testing auth", key)
		ch := client.Authenticate(&test.request)
		resp := <-ch
		if resp.Code != test.result.Code {
			t.Error("expected HTTP code", test.result.Code, "got", resp.Code)
		}
		_, err := resp.Validate()
		if err != nil && err.Error() != test.result.Message {
			t.Error("expected name", test.result.Message, "got", err.Error())
		}

		identity := resp.GetIdentity()
		if identity == nil { // if no identify no point to check further
			continue
		}
		if identity.Country != test.result.Country {
			t.Error("expected country", test.result.Country, "got",
				identity.Country)
		}
		if identity.CommonName != test.result.CommonName {
			t.Error("expected name", test.result.CommonName, "got",
				identity.CommonName)
		}
		if identity.SerialNumber != test.result.SerialNumber {
			t.Error("expected personal id", test.result.SerialNumber, "got",
				identity.SerialNumber)
		}
		certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
		// Might be problems with ca-certificates
		// if ok, err := resp.Cert.Verify(certPaths); !ok {
		// 	t.Error(err)
		// }
		_, _ = resp.Cert.Verify(certPaths)
	}
}

var clientTestTableSign = ClientTestTable{
	"client_ee_id_ok": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			Identifier: NewSemanticIdentifier(
				IdentifierTypePNO,
				CountryEE,
				"30303039914"),
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,OK",
				SerialNumber: "PNOEE-30303039914",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
}

func TestSign(t *testing.T) {
	for key, test := range clientTestTableSign {
		fmt.Println("Testing sign", key)
		ch := client.Sign(&test.request)
		resp := <-ch
		if resp.Code != test.result.Code {
			t.Error("expected HTTP code", test.result.Code, "got", resp.Code)
		}
		identity := resp.GetIdentity()
		if identity == nil { // if no identify no point to check further
			continue
		}
		if identity.Country != test.result.Country {
			t.Error("expected country", test.result.Country, "got",
				identity.Country)
		}
		if identity.CommonName != test.result.CommonName {
			t.Error("expected name", test.result.CommonName, "got",
				identity.CommonName)
		}
		if identity.SerialNumber != test.result.SerialNumber {
			t.Error("expected personal id", test.result.SerialNumber, "got",
				identity.SerialNumber)
		}
		_, err := resp.Validate()
		if err != nil {
			t.Error("Invalid response", err.Error())
		}
		certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
		// Might be problems with ca-certificates
		// if ok, err := resp.Cert.Verify(certPaths); !ok {
		// 	t.Error(err)
		// }
		_, _ = resp.Cert.Verify(certPaths)
	}
}

var clientTestTableSignFailed = ClientTestTable{
	"client_ee_doc_ageu18": {
		request: AuthRequest{
			RelyingPartyUUID: demoPartyUUID,
			RelyingPartyName: demoPartyName,
			Hash:             GenerateAuthHash(SHA512),
			HashType:         SHA512,
			Identifier:       "PNOEE-50701019992-9ZN6-Q",
			AuthType:         AuthTypeDocument,
		},
		result: ClientTestResult{
			Identity{
				Country:      "EE",
				CommonName:   "TESTNUMBER,MINOR",
				SerialNumber: "PNOEE-50701019992",
			},
			Response{
				Code:    http.StatusOK,
				Message: SessionResultOK,
			},
		},
	},
}

func TestSignFailedCert(t *testing.T) {
	for key, test := range clientTestTableSignFailed {
		fmt.Println("Testing sign", key)
		ch := client.Sign(&test.request)
		resp := <-ch
		if ok, _ := resp.Validate(); ok == true {
			t.Error("Certificate is valid. Expcted", ok, "got", !ok)
		}
	}
}

func TestSignExtended(t *testing.T) {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
	request := AuthRequest{
		// Replace in production with real RelyingPartyUUID.
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		// Replace in production with real RelyingPartyName.
		RelyingPartyName: "DEMO",
		// It is good to generate new has for security reasons.
		Hash: GenerateAuthHash(SHA384),
		// HashType should be the same as in GenerateAuthHash.
		HashType: SHA384,
		// We use personal ID as Identifier, also possible to use document
		// number.
		Identifier:       semid,
		AuthType:         AuthTypeEtsi,
		CertificateLevel: CertLevelQualified,
		AllowedInteractionsOrder: []AllowedInteractionsOrder{
			{
				Type:          InteractionVerificationCodeChoice,
				DisplayText60: "Welcome to Smart-ID!",
			},
			{
				Type:          InteractionDisplayTextAndPIN,
				DisplayText60: "Welcome to Smart-ID!",
			},
		},
	}

	resp := <-client.Sign(&request)

	if _, err := resp.Validate(); err != nil {
		log.Fatalln(err)
	}

	identity := resp.GetIdentity()
	fmt.Println("Name:", identity.CommonName)
	fmt.Println("Personal ID:", identity.SerialNumber)
	fmt.Println("Country:", identity.Country)
}
