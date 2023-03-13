package smartid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const pollDown = 1000
const pollUp = 120000

// Client is used to interact with endpoints and make requests and receive,
// responses.
type Client struct {
	// APIUrl sets the base API URL. Check official documentation
	// https://github.com/SK-EID/smart-id-documentation/wiki/Environment-technical-parameters
	APIUrl string

	// Poll defines poll timeout value in milliseconds. The upper bound of
	// timeout is 120000, the minimum is 1000. If not specified by the client,
	// a value halfway between maximum and minimum is used.
	Poll uint32

	httpClient *http.Client
}

// Option interface used for setting optional Client properties.
type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (o optionFunc) apply(c *Client) { o(c) }

// WithHTTPClient specifies which http client to use.
func WithHTTPClient(httpClient *http.Client) Option {
	return optionFunc(func(c *Client) { c.httpClient = httpClient })
}

// NewClient creates a new client instance. Poll will be in range 1000ms to
// 120000ms.
func NewClient(url string, poll uint32, opts ...Option) *Client {
	client := &Client{
		APIUrl: url,
		Poll:   poll,
	}

	for _, v := range opts {
		v.apply(client)
	}

	if client.httpClient == nil {
		client.httpClient = new(http.Client)
	}

	return client
}

// Authenticate does authentication in asynchronous way using channel.
func (c *Client) Authenticate(ctx context.Context, req *AuthRequest) chan *SessionResponse {
	ch := make(chan *SessionResponse)
	go func() {
		resp, err := c.AuthenticateSync(ctx, req)
		if err != nil {
			ch <- &SessionResponse{
				Response: Response{
					Code:    err.(*Error).Code,
					Message: err.(*Error).Message,
				},
			}
		} else {
			ch <- resp
		}
		close(ch)
	}()
	return ch
}

// AuthenticateSync does authentication in synchronous way.
func (c *Client) AuthenticateSync(ctx context.Context, req *AuthRequest) (*SessionResponse, error) {
	session, err := c.newSession(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := session.getResponse(ctx, c)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Sign does signing in asynchronous way using channel. Sign is very similar
// to Authenticate, but uses other endpoint.
func (c *Client) Sign(ctx context.Context, req *AuthRequest) chan *SessionResponse {
	req.endpoint = EndpointSignature
	return c.Authenticate(ctx, req)
}

// SignSync does signing in synchronous way. SignSync is very similar to
// AuthenticateSync, but uses other endpoint.
func (c *Client) SignSync(ctx context.Context, req *AuthRequest) (*SessionResponse, error) {
	req.endpoint = EndpointSignature
	return c.AuthenticateSync(ctx, req)
}

// --------------- unexposed -----------------

// newSession contacts Smart-ID service for authentication to get
// session ID in UUID format. This step also sends interaction order
// to user's app.
func (c *Client) newSession(ctx context.Context, req *AuthRequest) (*Session, error) {
	// Set some defaults fallback
	if req.CertificateLevel == "" {
		req.CertificateLevel = CertLevelQualified
	}
	if req.HashType == "" {
		req.HashType = SHA512
	}
	if req.AuthType == "" {
		req.AuthType = AuthTypeEtsi
	}
	if req.endpoint == "" {
		req.endpoint = EndpointAuthentication
	}
	if len(req.AllowedInteractionsOrder) == 0 {
		req.AllowedInteractionsOrder = []AllowedInteractionsOrder{
			{
				Type:          InteractionDisplayTextAndPIN,
				DisplayText60: "Welcome to Smart-ID!",
			},
		}
	}
	// end of defaults fallback

	resp, err := c.getEndpointResponse(ctx, req)
	if err != nil {
		return nil, err
	}

	return &Session{
		SessionID:        resp.SessionID,
		hash:             req.Hash,
		certificateLevel: req.CertificateLevel,
	}, nil
}

// getEndpointResponse makes authentication request to the endpoint.
func (c *Client) getEndpointResponse(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	url := fmt.Sprintf(
		"%v%v/%v/%v",
		c.APIUrl, req.endpoint, req.AuthType, req.Identifier,
	)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpResp, err := makeHTTPRequest(ctx, c.httpClient, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	resp := AuthResponse{
		Response: Response{
			Code:    httpResp.StatusCode,
			Message: resolveHTTPStatus(httpResp.StatusCode),
		},
	}

	if !resp.IsStatusOK() {
		return nil, &Error{
			Err:     errors.New("Status is NOK"),
			Code:    resp.Code,
			Message: resp.Message,
		}
	}

	body, err := getHTTPResponseBody(httpResp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// getSessionResponse makes request to the session endpoint.
func (c *Client) getSessionResponse(
	ctx context.Context,
	req *SessionRequest, s Session,
) (*SessionResponse, error) {
	url := fmt.Sprintf("%vsession/%v", c.APIUrl, req.SessionID)
	if c.Poll != 0 {
		url += fmt.Sprintf("?timeoutMs=%v", c.Poll)
	}

	httpResp, err := makeHTTPRequest(
		ctx, c.httpClient, http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, err
	}

	resp := SessionResponse{
		Response: Response{
			Code:    httpResp.StatusCode,
			Message: resolveHTTPStatus(httpResp.StatusCode),
		},
		Session: s,
	}

	body, err := getHTTPResponseBody(httpResp)
	if err != nil {
		return nil, err
	}

	errJSON := json.Unmarshal(body, &resp)
	if errJSON != nil {
		return nil, err
	}

	if resp.IsCompleted() && resp.IsFailed() {
		resp.Message = resp.GetFailureReason()
		return &resp, nil
	}

	// Make this expensive operation here, to make certificate available
	// for all required methods in Cert.
	resp.Cert.createX509CertIfNeeded()
	return &resp, nil
}

// getHTTPResponseBody extracts response body from HTTP response.
func getHTTPResponseBody(r *http.Response) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

// makeHTTPRequest makes just a HTTP request.
func makeHTTPRequest(
	ctx context.Context,
	httpClient *http.Client,
	mthd, url string,
	payld []byte,
) (*http.Response, error) {
	rd := bytes.NewReader(payld)
	httpReq, err := http.NewRequestWithContext(ctx, mthd, url, rd)
	if err != nil {
		return nil, err
	}
	cLen := strconv.Itoa(len(payld))
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Content-Length", cLen)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// resolveHTTPStatus resolves custom HTTP statuses. These non-standard
// statuses, only related to Smart-ID services.
// See more:
// https://github.com/SK-EID/smart-id-documentation/blob/master/README.md#211-http-status-codes
// https://github.com/SK-EID/smart-id-documentation/blob/master/README.md#233-http-status-code-usage
func resolveHTTPStatus(st int) string {
	switch st {
	case 471:
		return "No suitable account of requested type found, " +
			"but user has some other accounts."
	case 472:
		return "Person should view Smart-ID app or Smart-ID self-service " +
			"portal now."
	case 480:
		return "The client is too old and not supported any more."
	case 580:
		return "System is under maintenance, retry again later."
	default:
		return http.StatusText(st)
	}
}
