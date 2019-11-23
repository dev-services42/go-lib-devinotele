package devinotele

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	BasePath     = "https://integrationapi.net/rest"
	TokenTimeout = time.Minute * 120
)

type Client struct {
	login         string
	password      string
	token         string
	tokenM        *sync.Mutex
	tokenDeadline time.Time
	httpClient    *http.Client
}

func New(login, password string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		login:      login,
		password:   password,
		token:      "",
		tokenM:     new(sync.Mutex),
		httpClient: httpClient,
	}, nil
}

func (c *Client) fetchToken() (string, error) {
	params := url.Values{}
	params.Set("login", c.login)
	params.Set("password", c.password)
	u := BasePath + "/user/sessionid?" + params.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var token string
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token, nil
}

func (c *Client) refreshToken() error {
	c.tokenM.Lock()
	defer c.tokenM.Unlock()

	token, err := c.fetchToken()
	if err != nil {
		return err
	}

	c.token = token
	c.tokenDeadline = time.Now().Add(TokenTimeout)
	return nil
}

func (c *Client) getToken() (string, error) {
	if c.token == "" || time.Now().After(c.tokenDeadline) {
		if err := c.refreshToken(); err != nil {
			return "", err
		}
	}

	return c.token, nil
}

func (c *Client) getBaseQuery() (url.Values, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("SessionId", token)
	return q, nil
}

// SendSms отправляет сообщение на указанный номер
// text Текст сообщения, не более 2000 символов
// source Адрес отправителя, не более 11 латинских символов или 15 цифр
func (c *Client) SendSms(phoneNumber string, source string, text string) ([]string, error) {
	params, err := c.getBaseQuery()
	if err != nil {
		return nil, err
	}

	params.Set("DestinationAddress", phoneNumber)
	params.Set("SourceAddress", source)
	params.Set("Data", text)

	u := BasePath + "/Sms/Send?" + params.Encode()

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	//Обработка статусов ответа 400 - 451
	if (http.StatusBadRequest <= resp.StatusCode) && (resp.StatusCode <= http.StatusUnavailableForLegalReasons) {
		var errorResponse ErrorResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errorResponse); err != nil {
			return nil, BadRequest
		}

		switch err := errorResponse.Err(); err {
		case nil:
			panic("error response does not have error")
		case ErrUnauthorizedAccess:
			c.refreshToken()
			return nil, err
		default:
			return nil, err
		}
	}

	var messageIDs []string
	if err := json.NewDecoder(resp.Body).Decode(&messageIDs); err != nil {
		return nil, err
	}

	return messageIDs, nil
}

// GetMessageInfo получает состояние сообщения
func (c *Client) GetMessageInfo(id string) (*MessageInfo, error) {
	params, err := c.getBaseQuery()
	if err != nil {
		return nil, err
	}

	params.Set("messageId", id)

	u := BasePath + "/Sms/State?" + params.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//Обработка статусов ответа 400 - 451
	if (http.StatusBadRequest <= resp.StatusCode) && (resp.StatusCode <= http.StatusUnavailableForLegalReasons) {
		var errorResponse ErrorResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errorResponse); err != nil {
			return nil, BadRequest
		}

		switch err := errorResponse.Err(); err {
		case nil:
			panic("error response does not have error")
		case ErrUnauthorizedAccess:
			c.refreshToken()
			return nil, err
		default:
			return nil, err
		}
	}

	var messageState MessageInfo
	if err := json.NewDecoder(resp.Body).Decode(&messageState); err != nil {
		return nil, err
	}

	return &messageState, nil
}

func (c *Client) GetBalance() (float64, error) {
	params, err := c.getBaseQuery()
	if err != nil {
		return 0, err
	}

	u := BasePath + "/User/Balance?" + params.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	//Обработка статусов ответа 400 - 451
	if (http.StatusBadRequest <= resp.StatusCode) && (resp.StatusCode <= http.StatusUnavailableForLegalReasons) {
		var errorResponse ErrorResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errorResponse); err != nil {
			return 0, BadRequest
		}

		switch err := errorResponse.Err(); err {
		case nil:
			panic("error response does not have error")
		case ErrUnauthorizedAccess:
			c.refreshToken()
			return 0, err
		default:
			return 0, err
		}
	}

	var balance float64
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
