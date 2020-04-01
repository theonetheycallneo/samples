package beater

import "net/http"

const authHeader = "x-api-key"

type TokenAuth struct {
	token string
}

func (a TokenAuth) Enabled() bool { return a.token != "" }

func (a TokenAuth) Verify(req *http.Request) bool {
	if req.Header.Get(authHeader) == a.token {
		return true
	}
	return false
}
