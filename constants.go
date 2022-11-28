package smartid

// Authentication by ETSI or by document number.
const (
	AuthTypeEtsi     = "etsi" // (default)
	AuthTypeDocument = "document"
)

const (
	// InteractionDisplayTextAndPIN shows PIN and message in the
	// user's app.
	InteractionDisplayTextAndPIN = "displayTextAndPIN"

	// InteractionVerificationCodeChoice allows user to choice
	// verification code in the user's app.
	InteractionVerificationCodeChoice = "verificationCodeChoice"

	// InteractionConfirmationMessage shows confirmation messages in the user's app.
	InteractionConfirmationMessage = "confirmationMessage"

	// InteractionConfirmationMessageAndVerificationCodeChoice shows both
	// code and message in the user's app.
	InteractionConfirmationMessageAndVerificationCodeChoice = "confirmationMessageAndVerificationCodeChoice"
)

// API endpoints. There are currently supported 2 endpoints for requests.
const (
	EndpointAuthentication = "authentication" // default
	EndpointSignature      = "signature"
	// EndpointPrivate     = "private" // TODO Not implemenented yet
)

// AllowedInteractionsOrder allows you to interact with the user's app.
// For example, display the message or choose a verification code. Not
// all apps can support all interaction types. You can use many of them
// for fallback. The most common one is displayTextAndPIN. It should be
// used as the last fallback if you are going to use AllowedInteractionsOrder.
type AllowedInteractionsOrder struct {
	// Type is used to define interaction type.
	// There are 4 interaction types:
	//
	//	- displayTextAndPIN
	//	- verificationCodeChoice
	//	- confirmationMessage
	//	- confirmationMessageAndVerificationCodeChoice
	Type string `json:"type"`

	// DisplayText60 allows to enter up to 60 characters of text.
	DisplayText60 string `json:"displayText60,omitempty"`

	// DisplayText200 allows to enter up to 200 characters of text.
	DisplayText200 string `json:"displayText200,omitempty"`
}
