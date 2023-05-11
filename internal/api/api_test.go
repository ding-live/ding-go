package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ding-live/ding-go/pkg/status"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeFmt = "2006-01-02T15:04:05.999999999Z"

func TestInvalidAPIKey(t *testing.T) {
	ts := testServer("")

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testInvalidApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	_, authErr := a.Authentication(context.Background(), AuthRequest{})
	assert.ErrorIs(t, authErr, ErrUnauthorized)
}

func TestParseAuthSuccess(t *testing.T) {
	testUUID := uuid.New()
	testCreatedAt := time.Now().Add(-time.Minute).UTC()
	testExpiresAt := time.Now().UTC()

	rawRes := fmt.Sprintf(`{
    	"authentication_uuid": "%s",
    	"status": "pending",
    	"created_at": "%s",
    	"expires_at": "%s"
	}`, testUUID.String(), testCreatedAt.Format(timeFmt), testExpiresAt.Format(timeFmt))

	ts := testServer(rawRes)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	res, err := a.Authentication(context.Background(), AuthRequest{})
	require.NoError(t, err)

	assert.Equal(t, &AuthenticationResponse{
		Success: &AuthSuccessResponse{
			AuthenticationUUID: testUUID.String(),
			Status:             status.AuthPending,
			CreatedAt:          testCreatedAt,
			ExpiresAt:          testExpiresAt,
		},
	}, res)
}

func TestParseError(t *testing.T) {
	rawRes := `{
    	"code": "invalid_phone_number",
    	"message": "+invalid is not a valid phone number",
    	"doc_url": "https://docs.example.com/api/error-handling#invalid_phone_number"
	}`

	ts := testServer(rawRes, http.StatusBadRequest)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	res, err := a.Authentication(context.Background(), AuthRequest{})
	require.NoError(t, err)

	assert.Equal(t, &AuthenticationResponse{
		Error: &ErrorResponse{
			Code:    ErrorCodeInvalidPhoneNumber,
			Message: "+invalid is not a valid phone number",
			DocURL:  "https://docs.example.com/api/error-handling#invalid_phone_number",
		},
	}, res)
}

func TestParseCheckSuccess(t *testing.T) {
	testUUID := uuid.New()

	rawRes := fmt.Sprintf(`{
    	"authentication_uuid": "%s",
    	"status": "valid"
	}`, testUUID.String())

	ts := testServer(rawRes)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	res, err := a.Check(context.Background(), CheckRequest{})
	require.NoError(t, err)

	assert.Equal(t, &CheckResponse{
		Success: &CheckSuccessResponse{
			AuthenticationUUID: testUUID.String(),
			Status:             status.CheckValid,
		},
	}, res)
}

func TestParseRetrySuccess(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now().Add(-time.Minute).UTC()
	nextRetryAt := time.Now().UTC()

	rawRes := fmt.Sprintf(`{
    	"authentication_uuid": "%s",
    	"status": "expired_auth",
    	"created_at": "%s",
    	"next_retry_at": "%s",
    	"remaining_retry": 0
	}`, id.String(), createdAt.Format(timeFmt), nextRetryAt.Format(timeFmt))

	ts := testServer(rawRes)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	res, err := a.Retry(context.Background(), RetryRequest{})
	require.NoError(t, err)

	assert.Equal(t, &RetryResponse{
		Success: &RetrySuccessResponse{
			AuthenticationUUID: id.String(),
			Status:             status.RetryExpiredAuth,
			CreatedAt:          createdAt,
			NextRetryAt:        nextRetryAt,
			RemainingRetry:     0,
		},
	}, res)
}

func TestParseInvalidResponse(t *testing.T) {
	rawRes := "this is not json"

	ts := testServer(rawRes)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	_, authErr := a.Authentication(context.Background(), AuthRequest{})
	require.ErrorIs(t, authErr, ErrInternal)
}

func TestUnknownAuthStatus(t *testing.T) {
	testUUID := uuid.New()
	testCreatedAt := time.Now().Add(-time.Minute).UTC()
	testExpiresAt := time.Now().UTC()

	rawRes := fmt.Sprintf(`{
    	"authentication_uuid": "%s",
    	"status": "--------------------------------",
    	"created_at": "%s",
    	"expires_at": "%s"
	}`, testUUID.String(), testCreatedAt.Format(timeFmt), testExpiresAt.Format(timeFmt))

	ts := testServer(rawRes)

	a, err := New(Config{
		BaseURL:          ts.URL,
		APIKey:           testApiKey,
		CustomHTTPClient: ts.Client(),
		LeveledLogger:    testLogger{},
	})
	require.NoError(t, err)

	res, err := a.Authentication(context.Background(), AuthRequest{})
	require.NoError(t, err)

	assert.Equal(t, &AuthenticationResponse{
		Success: &AuthSuccessResponse{
			AuthenticationUUID: testUUID.String(),
			Status:             status.AuthUnknown,
			CreatedAt:          testCreatedAt,
			ExpiresAt:          testExpiresAt,
		},
	}, res)
}

// ----------------------------------------------------------------------------

type testLogger struct{}

func (testLogger) Debugf(string, ...interface{}) {}
func (testLogger) Infof(string, ...interface{})  {}
func (testLogger) Warnf(string, ...interface{})  {}
func (testLogger) Errorf(string, ...interface{}) {}
