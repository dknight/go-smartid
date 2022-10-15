package smartid

import (
	"fmt"
	"log"
)

func ExampleGenerateAuthHash() {
	hash := GenerateAuthHash(SHA512)
	fmt.Println(hash)
	// Output: [46 117 219 69 255 193 115 74 0 96 133 66 216 167 99 93 127 89 158 75 218 203 252 240 196 213 171 133 188 200 23 170 70 31 27 209 213 109 225 183 46 78 169 27 148 118 58 120 142 199 100 164 235 69 107 157 219 201 143 1 112 244 171 183]
}

func ExampleAuthHash_ToBase64String() {
	hash := GenerateAuthHash(SHA512)
	fmt.Println(hash.ToBase64String())
	// Output: LnXbRf/Bc0oAYIVC2KdjXX9Znkvay/zwxNWrhbzIF6pGHxvR1W3hty5OqRuUdjp4jsdkpOtFa53byY8BcPSrtw==
}

func ExampleAuthHash_EncodeBase64() {
	hash := GenerateAuthHash(SHA512)
	fmt.Println(hash.EncodeBase64())
	// Output: [76 110 88 98 82 102 47 66 99 48 111 65 89 73 86 67 50 75 100 106 88 88 57 90 110 107 118 97 121 47 122 119 120 78 87 114 104 98 122 73 70 54 112 71 72 120 118 82 49 87 51 104 116 121 53 79 113 82 117 85 100 106 112 52 106 115 100 107 112 79 116 70 97 53 51 98 121 89 56 66 99 80 83 114 116 119 61 61]
}

func ExampleAuthHash_CalculateVerificationCode() {
	hash := GenerateAuthHash(SHA512)
	fmt.Println(hash.CalculateVerificationCode())
	// Output: 3174
}

func ExampleNewSemanticIdentifier() {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "12345678901")
	fmt.Println(semid)
	// Output: PNOEE-12345678901
}

func ExampleClient_AuthenticateSync() {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
	request := AuthRequest{
		// Replace in production with real RelyingPartyUUID.
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		// Replace in production with real RelyingPartyName.
		RelyingPartyName: "DEMO",
		// It is good to generate new has for security reasons.
		Hash: GenerateAuthHash(SHA512),
		// We use personal ID as Identifier, also possible to use document
		// number.
		Identifier: semid,
		AuthType:   AuthTypeEtsi,
	}

	resp, err := client.AuthenticateSync(&request)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := resp.Validate(); err != nil {
		log.Fatalln(err)
	}

	// It is also good to verify the certificate over secure. But it isn't
	// mandatory, but strongly recommended.
	//
	// If you always get an error: x509: certificate signed by unknown
	// authority. Most probably you need install ca-certificates for
	// example for GNU Linux.
	//
	// sudo apt-get install ca-certificates
	// sudo dnf install ca-certificates
	certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
	if ok, err := resp.Cert.Verify(certPaths); !ok {
		log.Fatalln(err)
	}

	identity := resp.GetIdentity()
	fmt.Println("Name:", identity.CommonName)
	fmt.Println("Personal ID:", identity.SerialNumber)
	fmt.Println("Country:", identity.Country)
	// Output:
	// Name: TESTNUMBER,QUALIFIED OK1
	// Personal ID: PNOEE-30303039914
	// Country: EE
}
func ExampleClient_Authenticate() {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
	request := AuthRequest{
		// Replace in production with real RelyingPartyUUID.
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		// Replace in production with real RelyingPartyName.
		RelyingPartyName: "DEMO",
		// It is good to generate new has for security reasons.
		Hash: GenerateAuthHash(SHA512),
		// We use personal ID as Identifier, also possible to use document
		// number.
		Identifier: semid,
		AuthType:   AuthTypeEtsi,
	}

	resp := <-client.Authenticate(&request)
	if _, err := resp.Validate(); err != nil {
		log.Fatalln(err)
	}

	// It is also good to verify the certificate over secure. But it isn't
	// mandatory, but strongly recommended.
	//
	// If you always get an error: x509: certificate signed by unknown
	// authority. Most probably you need install ca-certificates for
	// example for GNU Linux.
	//
	// sudo apt-get install ca-certificates
	// sudo dnf install ca-certificates
	certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
	if ok, err := resp.Cert.Verify(certPaths); !ok {
		log.Fatalln(err)
	}

	identity := resp.GetIdentity()
	fmt.Println("Name:", identity.CommonName)
	fmt.Println("Personal ID:", identity.SerialNumber)
	fmt.Println("Country:", identity.Country)
	// Output:
	// Name: TESTNUMBER,QUALIFIED OK1
	// Personal ID: PNOEE-30303039914
	// Country: EE
}

func ExampleClient_SignSync() {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
	request := AuthRequest{
		// Replace in production with real RelyingPartyUUID.
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		// Replace in production with real RelyingPartyName.
		RelyingPartyName: "DEMO",
		// It is good to generate new has for security reasons.
		Hash: GenerateAuthHash(SHA512),
		// We use personal ID as Identifier, also possible to use document
		// number.
		Identifier: semid,
		AuthType:   AuthTypeEtsi,
	}

	resp, err := client.AuthenticateSync(&request)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := resp.Validate(); err != nil {
		log.Fatalln(err)
	}

	// It is also good to verify the certificate over secure. But it isn't
	// mandatory, but strongly recommended.
	//
	// If you always get an error: x509: certificate signed by unknown
	// authority. Most probably you need install ca-certificates for
	// example for GNU Linux.
	//
	// sudo apt-get install ca-certificates
	// sudo dnf install ca-certificates
	certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
	if ok, err := resp.Cert.Verify(certPaths); !ok {
		log.Fatalln(err)
	}

	identity := resp.GetIdentity()
	fmt.Println("Name:", identity.CommonName)
	fmt.Println("Personal ID:", identity.SerialNumber)
	fmt.Println("Country:", identity.Country)
	// Output:
	// Name: TESTNUMBER,QUALIFIED OK1
	// Personal ID: PNOEE-30303039914
	// Country: EE
}
func ExampleClient_Sign() {
	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
	request := AuthRequest{
		// Replace in production with real RelyingPartyUUID.
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		// Replace in production with real RelyingPartyName.
		RelyingPartyName: "DEMO",
		// It is good to generate new has for security reasons.
		Hash: GenerateAuthHash(SHA512),
		// We use personal ID as Identifier, also possible to use document
		// number.
		Identifier: semid,
		AuthType:   AuthTypeEtsi,
	}

	resp := <-client.Authenticate(&request)
	if _, err := resp.Validate(); err != nil {
		log.Fatalln(err)
	}

	// It is also good to verify the certificate over secure. But it isn't
	// mandatory, but strongly recommended.
	//
	// If you always get an error: x509: certificate signed by unknown
	// authority. Most probably you need install ca-certificates for
	// example for GNU Linux.
	//
	// sudo apt-get install ca-certificates
	// sudo dnf install ca-certificates
	certPaths := []string{"./certs/sid_demo_sk_ee_2022_PEM.crt"}
	if ok, err := resp.Cert.Verify(certPaths); !ok {
		log.Fatalln(err)
	}

	identity := resp.GetIdentity()
	fmt.Println("Name:", identity.CommonName)
	fmt.Println("Personal ID:", identity.SerialNumber)
	fmt.Println("Country:", identity.Country)
	// Output:
	// Name: TESTNUMBER,QUALIFIED OK1
	// Personal ID: PNOEE-30303039914
	// Country: EE
}
