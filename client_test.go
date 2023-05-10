package ding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
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

func TestParseCallbackURL(t *testing.T) {
	hc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client, err := NewClient(Config{
		CustomHTTPClient: hc.Client(),
		CustomerUUID:     uuid.New().String(),
	})
	require.NoError(t, err)

	pn := phonenumbers.GetExampleNumber("US")

	_, authErr := client.Authenticate(AuthenticateOptions{
		PhoneNumber: phonenumbers.Format(pn, phonenumbers.E164),
		CallbackURL: String("invalid_callback_url"),
	})
	require.ErrorIs(t, authErr, ErrInvalidCallbackURL)

	_, authErr2 := client.Authenticate(AuthenticateOptions{
		PhoneNumber: phonenumbers.Format(pn, phonenumbers.E164),
		CallbackURL: String("https://example.com/callback"),
	})
	require.NotErrorIs(t, authErr2, ErrInvalidCallbackURL)
}
