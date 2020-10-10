package devinotele_test

import (
	"context"
	v2 "github.com/dev-services42/go-lib-devinotele/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

var (
	testApiKey = os.Getenv("TEST_DEVINO_API_KEY")
	testNumber = os.Getenv("TEST_DEVINO_PHONE")
	testSender = os.Getenv("TEST_DEVINO_SENDER")
)

func TestClient_SendSms(t *testing.T) {
	client, err := v2.New(
		v2.WithHttpClient(http.DefaultClient),
		v2.WithApiKey(testApiKey),
	)
	require.NoError(t, err)

	res, err := client.SendSms(context.Background(), v2.SendSmsRequest{
		IncludeSegmentID: true,
		Messages: []v2.SmsMessage{
			{
				From:        testSender,
				To:          testNumber,
				Text:        "Hello",
				Validity:    0,
				Priority:    0,
				CallbackURL: "",
				Options:     nil,
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, 1, len(res))
	assert.Equal(t, v2.MessageStatusOK, res[0].Status)
	assert.NotEmpty(t, res[0].SegmentIDs)
	assert.NotEmpty(t, res[0].MessageID)
	assert.Empty(t, res[0].Error)
}
