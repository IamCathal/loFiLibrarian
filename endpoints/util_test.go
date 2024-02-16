package endpoints

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	EXPECTED_IP = "ipHere"
)

func TestMain(m *testing.M) {
	// setup

	code := m.Run()

	os.Exit(code)
}

func TestSetupCORS(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	responseWriter := httptest.NewRecorder()

	SetupCORS(responseWriter, req)

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}

	for header, expectedValue := range expectedHeaders {
		assert.Equal(t, expectedValue, responseWriter.Header().Get(header))
	}
}

func TestGetRealIP_xRealIpHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	req.Header.Set("X-Real-Ip", EXPECTED_IP)

	assert.Equal(t, EXPECTED_IP, getRealIP(req))
}

func TestGetRealIP_remoteAddrHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	req.Header.Set(REMOTE_ADDR_HEADER, EXPECTED_IP)

	assert.Equal(t, EXPECTED_IP, getRealIP(req))
}

func TestGetRealIP_requestFieldFallback(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	req.RemoteAddr = EXPECTED_IP

	assert.Equal(t, EXPECTED_IP, getRealIP(req))
}

func TestIsValidInt(t *testing.T) {
	assert.True(t, isValidInt("5"))
}

func TestIsValidInt_zero(t *testing.T) {
	assert.True(t, isValidInt("0"))
}

func TestIsValidInt_negativeInteger(t *testing.T) {
	assert.True(t, isValidInt("-5"))
}

func TestIsValidInt_empty(t *testing.T) {
	assert.False(t, isValidInt(""))
}

func TestIsValidInt_float(t *testing.T) {
	assert.False(t, isValidInt("2.6"))
}

func TestIsValidInt_word(t *testing.T) {
	assert.False(t, isValidInt("chips"))
}
