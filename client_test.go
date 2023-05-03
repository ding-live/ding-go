package ding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParsePhoneNumber(t *testing.T) {
	hc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client, err := NewClient(Config{
		CustomHTTPClient: hc.Client(),
		CustomerUUID:     uuid.New().String(),
	})
	require.NoError(t, err)

	_, authErr := client.Authenticate(AuthenticateOptions{
		PhoneNumber: "invalid_phone_number",
	})
	require.ErrorIs(t, authErr, ErrInvalidPhoneNumber)
}

func TestParseCustomerUUID(t *testing.T) {
	_, err := NewClient(Config{
		CustomerUUID: "invalid_customer_uuid",
	})
	require.Error(t, err, ErrInvalidCustomerUUID)
}
