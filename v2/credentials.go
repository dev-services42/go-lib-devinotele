package v2

import "net/http"

type Credentials interface {
	Propagate(req *http.Request)
}

type ApiKeyCredentials string

func (c ApiKeyCredentials) Propagate(req *http.Request) {
	req.Header.Set("Authorization", "Key "+string(c))
}

type BasicAuthCredentials struct {
	username string
	password string
}

func (c BasicAuthCredentials) Propagate(req *http.Request) {
	req.SetBasicAuth(c.username, c.password)
}
