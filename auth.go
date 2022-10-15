package smartid

// AuthResponse is returned from authentication and signing endpoint
// responses.
type AuthResponse struct {
	// SessionID is UUID format.
	SessionID string

	// Response is embedded response to store return code and message.
	Response
}

// AuthRequest represents structure that contains authentication and
// signing properties required for request.
type AuthRequest struct {
	// RelyingPartyUUID is UUID of relying party. Issued by SK.
	RelyingPartyUUID string `json:"relyingPartyUUID"`

	// RelyingPartyName is the name of the relying party. Issued by SK.
	RelyingPartyName string `json:"relyingPartyName"`

	// CertificateLevel is the level of certificate. Used possible values:
	//	QUALIFIED (recommended, default)
	//	ADVANCED
	CertificateLevel string `json:"certificateLevel"`

	// Base64 encoded hash function output to be signed (base64 encoding
	// according to rfc4648).
	Hash AuthHash `json:"hash"`

	// Hash algorithm. At the moment used only SHA512
	HashType string `json:"hashType"`

	// Nonce set behavior when requester wants, it can override the
	// idempotent behavior inside of this timeframe using an optional
	// nonce parameter present for all POST requests. Normally, that
	// parameter can be omitted. Read more
	// https://github.com/SK-EID/smart-id-documentation#235-idempotent-behaviour
	Nonce string `json:"nonce,omitempty"`

	// Used only when agreed with Smart-ID provider.
	// When omitted request capabilities are derived from
	// CertificateLevel parameter.
	Capabilities []string `json:"capabilities,omitempty"`

	// An app can support different interaction flows and a relying party can
	// demand a particular flow with or without a fall back possibility.
	// Different interaction flows can support different amount of data to
	// display information to user.
	AllowedInteractionsOrder []AllowedInteractionsOrder `json:"allowedInteractionsOrder,omitempty"`

	// AuthType is the type of authentication. Can be etsi or document.
	//	 AuthTypeEtsi is for authentication by semantic identifier (default).
	//	 AuthTypeDocument is for authentication by document number.
	AuthType string

	// Identifier is the semantic identifier or document number. This
	// is identifier used for person's identication.
	Identifier string

	// endpoint is the API endendpoint
	endpoint string

	// Request Request object contains code and message.
	Request
}
