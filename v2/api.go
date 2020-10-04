package v2

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type SendSmsRequest struct {
	IncludeSegmentID bool         `json:"-"`
	Messages         []SmsMessage `json:"messages"`
}

type SmsMessage struct {
	From        string            `json:"from"`
	To          string            `json:"to"`
	Text        string            `json:"text"`
	Validity    int               `json:"validity"`
	Priority    int               `json:"priority"`
	CallbackURL string            `json:"callback_url"`
	Options     map[string]string `json:"options"`
}

type SendSmsResponse struct {
	MessageID  string
	SegmentIDs []string
}

func (c *Client) SendSms(ctx context.Context, req SendSmsRequest) (*SendSmsResponse, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, errors.Wrap(err, "cannot encode request")
	}

	query := url.Values{}
	query.Set("includeSegmentId", strconv.FormatBool(req.IncludeSegmentID))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseAddr+"/sms/messages?"+query.Encode(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request")
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do request")
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	switch {
	default:
		return nil, errors.Wrap(ErrBadResponse, "status code is unexpected")
	case resp.StatusCode == http.StatusOK:
		var response struct {
			Code        StatusCode `json:"code"`
			MessageID   string     `json:"messageId"`
			Description string     `json:"description"`
			SegmentIDs  []string   `json:"segmentsId"`
			Reasons     []Reason   `json:"reasons"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, errors.Wrap(err, "cannot decode response")
		}

		switch response.Code {
		default:
			return nil, errors.Wrap(ErrBadResponse, "unexpected response code")
		case StatusCodeOk:
			return &SendSmsResponse{
				MessageID:  response.MessageID,
				SegmentIDs: response.SegmentIDs,
			}, nil
		case StatusCodeRejected:
			return nil, ErrorResponse{
				Description: response.Description,
				Reasons:     response.Reasons,
			}
		}
	}
}
