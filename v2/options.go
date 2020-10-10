package devinotele

import (
	"github.com/pkg/errors"
	"net/http"
)

var ErrBadConfiguration = errors.New("bad configuration")

type Option func(client *Client) error

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) error {
		if httpClient == nil {
			return errors.Wrap(ErrBadConfiguration, "httpclient must be set")
		}

		client.httpClient = httpClient

		return nil
	}
}

func WithApiKey(apiKey string) Option {
	return func(client *Client) error {
		if apiKey == "" {
			return errors.Wrap(ErrBadConfiguration, "apiKey must be set")
		}

		client.credentials = ApiKeyCredentials(apiKey)

		return nil
	}
}

func WithBasicAuth(username, password string) Option {
	return func(client *Client) error {
		if username == "" || password == "" {
			return errors.Wrap(ErrBadConfiguration, "username and password must be set")
		}

		client.credentials = BasicAuthCredentials{
			username: username,
			password: password,
		}

		return nil
	}
}
