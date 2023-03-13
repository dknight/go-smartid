// Package smartid implements an interface in Go to work with the
// Smart-ID API (https://www.smart-id.com). Smart-ID is used to easily and
// safely authenticate and sign documents online using only a smart phone.
// Smart-ID is a popular method in the Baltic countries of Estonia, Latvia,
// and Lithuania for authenticating and signing documents online for banks,
// social media, government offices, and other institutions.
//
// The bare minimum required to make an authentication request. Demonstarates
// synchronous way.
//
//	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
//	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
//	request := AuthRequest{
//		// Replace in production with real RelyingPartyUUID.
//		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
//		// Replace in production with real RelyingPartyName.
//		RelyingPartyName: "DEMO",
//		// It is good to generate new has for security reasons.
//		Hash: GenerateAuthHash(SHA512),
//		// We use personal ID as Identifier, also possible to use document
//		// number.
//		Identifier: semid,
//		AuthType:   AuthTypeEtsi,
//	}
//
//	// This blocks thread until it completes
//	resp, err := client.AuthenticateSync(&request)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	if _, err := resp.Validate(); err != nil {
//		log.Fatalln(err)
//	}
//
//	// It is also good to verify the certificate over secure. But it isn't
//	// mandatory, but strongly recommended.
//	//
//	// If you always get an error: x509: certificate signed by unknown
//	// authority. Most probably you need install ca-certificates for
//	// example for GNU Linux.
//	//
//	// sudo apt-get install ca-certificates
//	// sudo dnf install ca-certificates
//	certPaths := []string{"./certs/TEST_of_EID-SK_2016.pem.crt"}
//	 if ok, err := resp.Cert.Verify(certPaths); !ok {
//		t.Error(err)
//	}
//
//	identity := resp.GetIdentity()
//	fmt.Println("Name:", identity.CommonName)
//	fmt.Println("Personal ID:", identity.SerialNumber)
//	fmt.Println("Country:", identity.Country)
//	// Output:
//	// Name: TESTNUMBER,QUALIFIED OK1
//	// Personal ID: PNOEE-30303039914
//	// Country: EE
//
// Another example contains many more quest parameters for the signing method.
// Sign and Authenticate methods are similar and you can use the same
// AuthRequest parameters for both of them.
//
// This examples is asynchronous uses channel.
//
//	semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
//	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
//	request := AuthRequest{
//		// Replace in production with real RelyingPartyUUID.
//		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
//		// Replace in production with real RelyingPartyName.
//		RelyingPartyName: "DEMO",
//		// It is good to generate new has for security reasons.
//		Hash: GenerateAuthHash(SHA384),
//		// HashType should be the same as in GenerateAuthHash.
//		HashType: SHA384,
//		// We use personal ID as Identifier, also possible to use document
//		// number.
//		Identifier: semid,
//		AuthType:   AuthTypeEtsi,
//		CertificateLevel: CertLevelQualified,
//		AllowedInteractionsOrder: []AllowedInteractionsOrder{
//			{
//				Type:          InteractionVerificationCodeChoice,
//				DisplayText60: "Welcome to Smart-ID!",
//			},
//			{
//				Type:          InteractionDisplayTextAndPIN,
//				DisplayText200: "Welcome to Smart-ID! A bit longer text."
//			},
//		},
//	}
//
//	resp := <-client.Sign(&request)
//
//	if _, err := resp.Validate(); err != nil {
//		log.Fatalln(err)
//	}
//
//	// It is also good to verify the certificate over secure. But it isn't
//	// mandatory, but strongly recommended.
//	//
//	// If you always get an error: x509: certificate signed by unknown
//	// authority. Most probably you need install ca-certificates for
//	// example for GNU Linux.
//	//
//	// sudo apt-get install ca-certificates
//	// sudo dnf install ca-certificates
//	certPaths := []string{"./certs/TEST_of_EID-SK_2016.pem.crt"}
//	 if ok, err := resp.Cert.Verify(certPaths); !ok {
//		t.Error(err)
//	}
//
//	identity := resp.GetIdentity()
//	fmt.Println("Name:", identity.CommonName)
//	fmt.Println("Personal ID:", identity.SerialNumber)
//	fmt.Println("Country:", identity.Country)
//	// Output:
//	// Name: TESTNUMBER,QUALIFIED OK1
//	// Personal ID: PNOEE-30303039914
//	// Country: EE
//
// Demonstration of siging with document number.
//
//	docid := "PNOEE-30303039914-1Q3P-Q"
//	client := NewClient("https://sid.demo.sk.ee/smart-id-rp/v2/", 5000)
//	request := AuthRequest{
//		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
//		RelyingPartyName: "DEMO",
//		Hash: GenerateAuthHash(SHA512),
//		Identifier: docid,
//		AuthType:   AuthTypeDocument,
//	}
//	resp := <-client.Authenticate(&request)
//	fmt.Printf("%+v\n", resp)
package smartid
