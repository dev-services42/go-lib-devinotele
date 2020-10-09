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

type MessageStatus int8

const (
	MessageStatusUnknown MessageStatus = iota
	MessageStatusOK
	MessageStatusError
)

type MessageResult struct {
	Status     MessageStatus
	MessageID  string
	SegmentIDs []string
	Error      ErrorResponse
}

func (c *Client) SendSms(ctx context.Context, req SendSmsRequest) ([]MessageResult, error) {
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

	httpReq.Header.Set("content-type", "application/json")
	c.credentials.Propagate(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "cannot do request")
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	// Do not check the status code, because devinotele not describe http
	// status codes.
	var response struct {
		Result []struct {
			Code        StatusCode `json:"code"`
			MessageID   string     `json:"messageId"`
			Description string     `json:"description"`
			SegmentIDs  []string   `json:"segmentsId"`
			Reasons     []Reason   `json:"reasons"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "cannot decode response")
	}

	statuses := make([]MessageResult, len(response.Result))
	for i := range response.Result {
		result := &response.Result[i]
		switch result.Code {
		default:
			statuses[i] = MessageResult{
				Status:     MessageStatusUnknown,
				MessageID:  "",
				SegmentIDs: nil,
				Error: ErrorResponse{
					Description: "Unknown message status code",
					Reasons:     nil,
				},
			}
		case StatusCodeOK:
			statuses[i] = MessageResult{
				Status:     MessageStatusOK,
				MessageID:  result.MessageID,
				SegmentIDs: result.SegmentIDs,
				Error: ErrorResponse{
					Description: "",
					Reasons:     nil,
				},
			}
		case StatusCodeRejected:
			statuses[i] = MessageResult{
				Status:     MessageStatusError,
				MessageID:  "",
				SegmentIDs: nil,
				Error: ErrorResponse{
					Description: result.Description,
					Reasons:     result.Reasons,
				},
			}
		}
	}

	return statuses, nil
}
