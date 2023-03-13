![Smart-ID in Go language](https://github.com/dknight/go-smartid/blob/main/files/smartid-gopher.png?raw=true)

# Intro

Package smartid implements an interface in Go to work with the
Smart-ID API (https://www.smart-id.com). Smart-ID is used to easily and
safely authenticate and sign documents online using only a smart phone.
Smart-ID is a popular method in the Baltic countries of Estonia, Latvia,
and Lithuania for authenticating and signing documents online for banks,
social media, government offices, and other institutions.

Official Smart-ID [technical documentation](https://github.com/SK-EID/smart-id-documentation/wiki).

## Installation

```sh
go get github.com/dknight/go-smartid
```

## Usage

The bare minimum required to make an authentication request. Demonstarates
synchronous way.

For more examples [see full docs](https://pkg.go.dev/github.com/dknight/go-smartid).

### Sync request

```go
semid := NewSemanticIdentifier(IdentifierTypePNO, CountryEE, "30303039914")
client := NewClient("https:sid.demo.sk.ee/smart-id-rp/v2/", 5000)
request := AuthRequest{
	// Replace in production with real RelyingPartyUUID.
	RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
	// Replace in production with real RelyingPartyName.
	RelyingPartyName: "DEMO",
	// It is good to generate new has for security reasons.
	Hash: GenerateAuthHash(SHA512),
 	// We use personal ID as Identifier, also possible to use document number.
	Identifier: semid,
}

// This blocks thread until it completes
resp, err := client.AuthenticateSync(context.TODO(), &request)
if err != nil {
	log.Fatalln(err)
}

if _, err := resp.Validate(); err != nil {
	log.Fatalln(err)
}

// It is also good to verify the certificate over secure. But it isn't
// mandatory, but strongly recommended.
//
certPaths := []string{"./certs/TEST_of_EID-SK_2016.pem.crt"}
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
```

### Async way using channel

Another example contains many more quest parameters for the signing method.
Sign and Authenticate methods are similar and you can use the same
AuthRequest parameters for both of them.

This examples is asynchronous uses channel.

```go
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
	Identifier: semid,
	AuthType:   AuthTypeEtsi,
	CertificateLevel: CertLevelQualified,
	AllowedInteractionsOrder: []AllowedInteractionsOrder{
		{
			Type:          InteractionVerificationCodeChoice,
			DisplayText60: "Welcome to Smart-ID!",
		},
		{
			Type:          InteractionDisplayTextAndPIN,
			DisplayText200: "Welcome to Smart-ID! A bit longer text.",
		},
	},
}

resp := <-client.Sign(context.TODO(), &request)
if _, err := resp.Validate(); err != nil {
	log.Fatalln(err)
}

// It is also good to verify the certificate over secure. But it isn't
// mandatory, but strongly recommended.
//
certPaths := []string{"./certs/TEST_of_EID-SK_2016.pem.crt"}
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
```

For more examples [see docs](http://missing-yet).

## What is not included

1. `:private/` endpoint.
2. Better certificated parsing and data extraction. You can get certificated
from response, verify and parse it in own way `response.Cert.GetX509Cert()`.
3. Smart-ID API version v1 is not supported, only v2.

## Testing

SK test environment **is very unstable**. Possible technical problems might be:

1. Problems with service availability. Doesn't work too often.
2. They change test data without any warning.
3. Problems with certificates.
4. Problems with performance requests can last very long time, sometimes
408 Timeout will be given as response.

```go
go test
```

# Troubleshooting

## Problems with certificates

### x509: certificate signed by unknown authority

If in development you get an error `x509: certificate signed by unknown
authority`. Then you need to install SK test certificates to your system.
Install certificates from directory `./certs` to your operating system.

Fedora Linux example:

```sh
sudo cp ./certs/TEST_of_* /usr/share/pki/ca-trust-source/anchors/
sudo update-ca-trust
```

Then you can verify your certificate, but don't forget to replace with your
personal certificate in production.

```go
certPaths := []string{"./certs/TEST_of_EID-SK_2016.pem.crt"}
if ok, err := resp.Cert.Verify(certPaths); !ok {
 	log.Fatalln(err)
}
```

## Contribution

Any help is appreciated. Found a bug, typo, inaccuracy, etc.? Please do not
hesitate and make pull request or issue.

## License

MIT 2022-2023
