package v2

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/http"
)

const baseAddr = "https://api.devino.online"

var ErrBadResponse = errors.New("bad response")

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	httpClient  *http.Client
	credentials Credentials
}

func New(opts ...Option) (*Client, error) {
	var c Client
	for i := range opts {
		if err := opts[i](&c); err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	return &c, nil
}
