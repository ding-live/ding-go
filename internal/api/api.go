package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ding-live/ding-go/pkg/status"
	"github.com/hashicorp/go-retryablehttp"
)

type API struct {
	baseURL       string
	apiKey        string
	hc            *http.Client
	leveledLogger LeveledLogger
}

type Config struct {
	BaseURL           string
	APIKey            string
	MaxNetworkRetries *int
	CustomHTTPClient  *http.Client
	LeveledLogger     LeveledLogger
}

func New(cfg Config) (*API, error) {
	client := retryablehttp.NewClient()

	if cfg.MaxNetworkRetries != nil {
		client.RetryMax = *cfg.MaxNetworkRetries
	}

	if cfg.CustomHTTPClient != nil {
		client.HTTPClient = cfg.CustomHTTPClient
	}

	if cfg.LeveledLogger != nil {
		client.Logger = convertLogger(cfg.LeveledLogger)
	}

	a := &API{
		baseURL:       cfg.BaseURL,
		apiKey:        cfg.APIKey,
		hc:            client.StandardClient(),
		leveledLogger: cfg.LeveledLogger,
	}

	if a.leveledLogger == nil {
		return nil, fmt.Errorf("missing logger")
	}

	return a, nil
}

const APIKeyHeader = "x-api-key"

type AuthSuccessResponse struct {
	AuthenticationUUID string      `json:"authentication_uuid"`
	Status             status.Auth `json:"status"`
	CreatedAt          time.Time   `json:"created_at"`
	ExpiresAt          time.Time   `json:"expires_at"`
}

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	DocURL  string    `json:"doc_url"`
}

type ErrorCode string

const (
	ErrorCodeInternalServer     ErrorCode = "internal_server_error"
	ErrorCodeBadRequest         ErrorCode = "bad_request"
	ErrorCodeInvalidPhoneNumber ErrorCode = "invalid_phone_number"
	ErrorCodeAccountInvalid     ErrorCode = "account_invalid"
	ErrorCodeNegativeBalance    ErrorCode = "negative_balance"
	ErrorCodeInvalidLine        ErrorCode = "invalid_line"
	ErrorCodeUnsupportedRegion  ErrorCode = "unsupported_region"
	ErrorCodeInvalidAuthUUID    ErrorCode = "invalid_auth_uuid"
)

type AuthRequest struct {
	PhoneNumber     string  `json:"phone_number,omitempty"`
	CustomerUUID    string  `json:"customer_uuid,omitempty"`
	IP              *string `json:"ip,omitempty"`
	DeviceID        *string `json:"device_id,omitempty"`
	DeviceType      *string `json:"device_type,omitempty"`
	AppVersion      *string `json:"app_version,omitempty"`
	CallbackURL     *string `json:"callback_url,omitempty"`
	IsReturningUser *bool   `json:"is_returning_user,omitempty"`
}

type CheckRequest struct {
	CustomerUUID       string `json:"customer_uuid"`
	AuthenticationUUID string `json:"authentication_uuid"`
	CheckCode          string `json:"check_code"`
}

type CheckSuccessResponse struct {
	AuthenticationUUID string       `json:"authentication_uuid"`
	Status             status.Check `json:"status"`
}

type RetryRequest struct {
	CustomerUUID       string `json:"customer_uuid"`
	AuthenticationUUID string `json:"authentication_uuid"`
}

type RetrySuccessResponse struct {
	AuthenticationUUID string       `json:"authentication_uuid"`
	Status             status.Retry `json:"status"`
	CreatedAt          time.Time    `json:"created_at"`
	NextRetryAt        time.Time    `json:"next_retry_at"`
	RemainingRetry     int          `json:"remaining_retry"`
}

type GatewayErrorMessage struct {
	Message string `json:"message"`
}

var (
	ErrInternal     = fmt.Errorf("internal error")
	ErrUnauthorized = fmt.Errorf("unauthorized")
)

// ----------------------------------------------------------------------------

// TODO(2024-25) -> abstract each endpoint parsing logic using generics
type AuthenticationResponse struct {
	Error   *ErrorResponse
	Success *AuthSuccessResponse
}

func (a *API) Authentication(ctx context.Context, req AuthRequest) (*AuthenticationResponse, error) {
	res, err := a.post(ctx, "authentication", req)
	if err != nil {
		return nil, ErrInternal
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		a.leveledLogger.Errorf("received a non-200 HTTP status %d", res.StatusCode)

		if res.StatusCode == http.StatusForbidden {
			return nil, ErrUnauthorized
		}

		var resp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			a.leveledLogger.Errorf("unable to decode response: %s", err)
			return nil, ErrInternal
		}

		return &AuthenticationResponse{
			Error: &resp,
		}, nil
	}

	var resp AuthSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		a.leveledLogger.Errorf("unable to decode response of HTTP OK status: %s", err)
		return nil, ErrInternal
	}

	return &AuthenticationResponse{
		Success: &resp,
	}, nil

}

type CheckResponse struct {
	Error   *ErrorResponse
	Success *CheckSuccessResponse
}

func (a *API) Check(ctx context.Context, req CheckRequest) (*CheckResponse, error) {
	res, err := a.post(ctx, "check", req)
	if err != nil {
		return nil, ErrInternal
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		a.leveledLogger.Errorf("received a non-200 HTTP status %d", res.StatusCode)

		if res.StatusCode == http.StatusForbidden {
			return nil, ErrUnauthorized
		}

		var resp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			a.leveledLogger.Errorf("unable to decode response: %s", err)
			return nil, ErrInternal
		}

		return &CheckResponse{
			Error: &resp,
		}, nil
	}

	var resp CheckSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		a.leveledLogger.Errorf("unable to decode response of HTTP OK status: %s", err)
		return nil, ErrInternal
	}

	return &CheckResponse{
		Success: &resp,
	}, nil
}

type RetryResponse struct {
	Error   *ErrorResponse
	Success *RetrySuccessResponse
}

func (a *API) Retry(ctx context.Context, req RetryRequest) (*RetryResponse, error) {
	res, err := a.post(ctx, "retry", req)
	if err != nil {
		return nil, ErrInternal
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		a.leveledLogger.Errorf("received a non-200 HTTP status %d", res.StatusCode)

		if res.StatusCode == http.StatusForbidden {
			return nil, ErrUnauthorized
		}

		var resp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			a.leveledLogger.Errorf("unable to decode response: %s", err)
			return nil, ErrInternal
		}

		return &RetryResponse{
			Error: &resp,
		}, nil
	}

	var resp RetrySuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		a.leveledLogger.Errorf("unable to decode response of HTTP OK status: %s", err)
		return nil, ErrInternal
	}

	return &RetryResponse{
		Success: &resp,
	}, nil
}

// ----------------------------------------------------------------------------

func (a *API) post(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		a.leveledLogger.Errorf("marshal request payload %v: %v", payload, err)
		return nil, ErrInternal
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s", a.baseURL, url),
		bytes.NewBuffer(body),
	)
	if err != nil {
		a.leveledLogger.Errorf("create HTTP request %v", err)
		return nil, ErrInternal
	}

	req.Header.Set(APIKeyHeader, a.apiKey)
	req.Header.Set("content-type", "application/json")

	res, err := a.hc.Do(req)
	if err != nil {
		a.leveledLogger.Errorf("perform HTTP request %v", err)
		return nil, ErrInternal
	}

	return res, nil
}

// ----------------------------------------------------------------------------

// LeveledLogger is an interface that can be implemented by any logger or a
// logger wrapper to provide leveled logging. The methods accept a message
// string and a variadic number of key-value pairs.
type LeveledLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type loggerShim struct {
	baseLogger LeveledLogger
}

func (l loggerShim) Error(msg string, keysAndValues ...interface{}) {
	l.baseLogger.Errorf(fmt.Sprint(msg, keysAndValues))
}

func (l loggerShim) Info(msg string, keysAndValues ...interface{}) {
	l.baseLogger.Infof(fmt.Sprint(msg, keysAndValues))
}

func (l loggerShim) Debug(msg string, keysAndValues ...interface{}) {
	l.baseLogger.Debugf(fmt.Sprint(msg, keysAndValues))
}

func (l loggerShim) Warn(msg string, keysAndValues ...interface{}) {
	l.baseLogger.Warnf(fmt.Sprint(msg, keysAndValues))
}

func convertLogger(logger LeveledLogger) retryablehttp.LeveledLogger {
	return loggerShim{logger}
}
