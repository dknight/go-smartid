package smartid

import (
	"fmt"
)

// Session response status. There are only 2 statuses available for Smart-ID
// service.
const (
	// SessionStatusComplete indicates that request is completed.
	SessionStatusComplete = "COMPLETE"

	// SessionStatusRunning indicates that request is still running.
	SessionStatusRunning = "RUNNING"
)

// Session response result codes.
const (
	SessionResultOK                                         = "OK"
	SessionResultUserRefusedCertChoice                      = "USER_REFUSED_CERT_CHOICE"
	SessionResultUserRefusedDisplayTextAndPIN               = "USER_REFUSED_DISPLAYTEXTANDPIN"
	SessionResultUserRefusedVCChoice                        = "USER_REFUSED_VC_CHOICE"
	SessionResultUserRefusedConfirmationMessage             = "USER_REFUSED_CONFIRMATIONMESSAGE"
	SessionResultUserRefusedConfirmationMessageWithVCChoice = "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE"
	SessionResultWrongVC                                    = "WRONG_VC"
	SessionResultTimeout                                    = "TIMEOUT"
	SessionResultRequiredInteractionNotSupportedByApp       = "REQUIRED_INTERACTION_NOT_SUPPORTED_BY_APP"
)

// Session represents information about session.
type Session struct {
	// SessionID is the session identified in UUID format.
	SessionID string `json:"sessionID"`

	// Hash for authentication and verification code.
	hash AuthHash

	// certificateLevel for certificate level check.
	certificateLevel string
}

// getResponse makes response to session endpoint API. It also polls from
// the service session status.
func (s Session) getResponse(c *Client) (*SessionResponse, error) {
	req := &SessionRequest{
		SessionID:        s.SessionID,
		hash:             s.hash,
		certificateLevel: s.certificateLevel,
	}
	for {
		resp, err := c.getSessionResponse(req, s)
		if err != nil {
			return nil, err
		}

		if resp.IsCompleted() {
			return resp, nil
		}
		// fmt.Println(resp.State) // DEBUG
	}
}

// SessionRequest represents structure that contains session properties required
// for request.
type SessionRequest struct {
	// SessionID is the session identified in UUID format.
	SessionID string

	// AuthHash hash for authentication check.
	hash AuthHash

	// certificatelevel for checking certificate level.
	certificateLevel string
}

// SessionResponse is used for session endpoint response.
type SessionResponse struct {
	// State is the status of session. There are only 2 statuses:
	//	- RUNNING
	//	- COMPLETE
	State string `json:"state"`

	// Result shows what the end result of the response. See result codes.
	Result Result `json:"result"`

	// Signature contains signature (Value and Algorithm).
	// Empty if not OK.
	Signature Signature `json:"signature"`

	// Cert contains signature (Value and CertificateLevel).
	// Empty if not OK.
	Cert Cert `json:"cert"`

	// InteractionFlowUsed show which interaction flow was used.
	InteractionFlowUsed string `json:"interactionFlowUsed,omitempty"`

	// DeviceIpAddress give IP address for the device where user's used
	// the app.
	DeviceIPAddress string `json:"deviceIpAddress"`

	// Session contains request parameters like SessionID and Hash.
	Session

	// Response contains code and message.
	Response
}

// GetFailureReason returns result code for the session failure.
// If the return value is SESSION_RESULT_OK, this means there is no failure.
//
// Possible result codes are:
//	- SessionResultOK
//	- SessionResultUserRefused
//	- SessionResultUserRefusedDisplayTextAndPIN
//	- SessionResultUserRefusedVCChoice
//	- SessionResultUserRefusedConfirmationMessage
//	- SessionResultUserRefusedConfirmationMessageWithVCChoice
//	- SessionResultUserWrongVC
//	- SessionResultUserTimeout
func (r *SessionResponse) GetFailureReason() string {
	if r.Result.EndResult == "" {
		return r.Message
	}
	return r.Result.EndResult
}

// Validate checks is session response is valid.
func (r *SessionResponse) Validate() (bool, error) {
	if !r.IsCompleted() {
		return false, fmt.Errorf("Response is not completed")
	}
	if r.IsFailed() {
		return false, fmt.Errorf(r.GetFailureReason())
	}
	if !r.IsValidSignature() {
		return false, fmt.Errorf("Invalid signature")
	}
	if r.Cert.IsExpired() {
		return false, fmt.Errorf("Certificate has expired")
	}
	if r.Cert.IsNotActive() {
		return false, fmt.Errorf("Certificate is not yet active")
	}
	if !r.Cert.IsSameLevel(r.certificateLevel) {
		return false, fmt.Errorf("Certificate level does not match")
	}
	return true, nil
}

// IsValidSignature checks validity of the signature.
func (r *SessionResponse) IsValidSignature() bool {
	return r.Signature.IsValid(r.Cert, r.hash)
}

// IsCompleted checks that response has completed. If the return value is
// empty, it also means not completed.
func (r *SessionResponse) IsCompleted() bool {
	return r.State == SessionStatusComplete || r.State == ""
}

// IsFailed checks that response is not successful.
func (r *SessionResponse) IsFailed() bool {
	return r.Result.EndResult != SessionResultOK
}

// GetIdentity gets user identity based on certificated.
func (r *SessionResponse) GetIdentity() *Identity {
	if r.IsFailed() {
		return nil
	}
	return newIdentity(r.Cert.GetSubject())
}

// GetIdentity gets user identity based on certificated.
func (r *SessionResponse) GetIssuerIdentity() *Identity {
	if r.IsFailed() {
		return nil
	}
	return newIdentity(r.Cert.GetIssuer())
}
