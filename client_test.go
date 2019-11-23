package devinotele_test

import (
	"github.com/kazhuravlev/go-lib-devinotele/devinotele"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var (
	testLogin    = os.GetEnv("TEST_DEVINO_LOGIN")
	testPassword = os.GetEnv("TEST_DEVINO_PASSWORD")
	testNumber   = os.GetEnv("TEST_DEVINO_PHONE")
	testSender   = os.GetEnv("TEST_DEVINO_SENDER")
)

func TestNew(t *testing.T) {
	c, err := devinotele.New(testLogin, testPassword, nil)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestNewWithClient(t *testing.T) {
	c, err := devinotele.New(testLogin, testPassword, http.DefaultClient)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestClient_GetBalance(t *testing.T) {
	c, _ := devinotele.New(testLogin, testPassword, nil)

	balance, err := c.GetBalance()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, float64(balance))
}

func TestClient_SendSms(t *testing.T) {
	c, _ := devinotele.New(testLogin, testPassword, nil)

	messageIDs, err := c.SendSms(testNumber, testSender, "Hello! 😎")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messageIDs))
}
