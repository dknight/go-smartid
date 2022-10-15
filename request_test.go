package smartid

import "testing"

func TestRequestError(t *testing.T) {
	msg := "Error"
	req := &RequestError{msg}
	if req.Error() != msg {
		t.Error("expected", msg, "got", req.Error())
	}
}
