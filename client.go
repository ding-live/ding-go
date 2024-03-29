package ding

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ding-live/ding-go/internal/api"
	"github.com/ding-live/ding-go/pkg/status"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

// Client is the main Ding client.
// It is used to interact with the Ding API.
type Client struct {
	api          api.API
	customerUUID string
}

// Config is the configuration required to instanciate a new client
type Config struct {
	// CustomerUUID is the UUID that was given to you during your onboarding
	CustomerUUID string

	// APIKey is your secret API key
	APIKey string

	// MaxNetworkRetries sets maximum number of times that the library will
	// retry requests that appear to have failed due to an intermittent
	// problem.
	//
	// This value is a pointer to allow us to differentiate an unset versus
	// empty value. Use ding.Int for an easy way to set this value.
	//
	// Defaults to 3.
	MaxNetworkRetries *int

	// CustomHTTPClient is an HTTP client instance to use when making API requests.
	//
	// If left unset, it'll be set to a default HTTP client for the package.
	CustomHTTPClient *http.Client

	// LeveledLogger is the logger that the will be used to log errors,
	// warnings, and informational messages.
	//
	// LeveledLogger is implemented by Logger, and one can be
	// initialized at the desired level of logging. LeveledLogger
	// also provides out-of-the-box compatibility with a Logrus Logger, but may
	// require a thin shim for use with other logging libraries that use less
	// standard conventions like Zap.
	//
	// Defaults to ding.Logger.
	LeveledLogger LeveledLogger
}

var (
	ErrUnauthorized        = errors.New("unauthorized, please check your API key")
	ErrInternal            = errors.New("an unhandled error occured")
	ErrInvalidPhoneNumber  = errors.New("invalid phone number")
	ErrInvalidCustomerUUID = errors.New("invalid account UUID")
	ErrNegativeBalance     = errors.New("negative balance")
	ErrUnsupportedRegion   = errors.New("unsupported region")
	ErrInvalidAuthUUID     = errors.New("invalid authentication UUID")
	ErrInvalidCallbackURL  = errors.New("invalid callback URL")
)

const apiBaseURL = "https://api.ding.live/v1"

// NewClient returns a new Ding client with a `Config` object.
func NewClient(cfg Config) (*Client, error) {
	if !isValidUUID(cfg.CustomerUUID) {
		return nil, ErrInvalidCustomerUUID
	}

	logger := cfg.LeveledLogger

	if logger == nil {
		logger = &DefaultLeveledLogger
	}

	api, err := api.New(api.Config{
		BaseURL:           apiBaseURL,
		APIKey:            cfg.APIKey,
		MaxNetworkRetries: cfg.MaxNetworkRetries,
		CustomHTTPClient:  cfg.CustomHTTPClient,
		LeveledLogger:     logger,
	})
	if err != nil {
		return nil, fmt.Errorf("create API client: %w", err)
	}

	return &Client{
		customerUUID: cfg.CustomerUUID,
		api:          *api,
	}, nil
}

// DeviceType is the type of device used to authenticate.
type DeviceType string

var (
	DeviceTypeAndroid DeviceType = "ANDROID"
	DeviceTypeIOS     DeviceType = "IOS"
	DeviceTypeWeb     DeviceType = "WEB"
)

func (d DeviceType) String() string {
	return string(d)
}

// AuthenticateOptions are the options used to authenticate a user. Only PhoneNumber
// is required. Other options are optional but recommended because they are used by
// the Ding antispam system.
type AuthenticateOptions struct {
	PhoneNumber     string
	IP              *string
	DeviceID        *string
	DeviceType      *DeviceType
	AppVersion      *string
	CallbackURL     *string
	IsReturningUser bool
}

// Authentication is the result of an authentication request.
type Authentication struct {
	AuthenticationUUID string
	Status             status.Auth
	CreatedAt          time.Time
	ExpiresAt          time.Time
}

// AuthenticateWithContext performs an authentication request against the Ding API that
// can be cancelled with a context. Authentication requests allow you to send a message to
// a given phone number with a code that the user will have to enter in your app.
func (c *Client) AuthenticateWithContext(ctx context.Context, opt AuthenticateOptions) (*Authentication, error) {
	if !isValidNumber(opt.PhoneNumber) {
		return nil, ErrInvalidPhoneNumber
	}

	if opt.CallbackURL != nil {
		if !isValidURL(*opt.CallbackURL) {
			return nil, ErrInvalidCallbackURL
		}
	}

	req := api.AuthRequest{
		PhoneNumber:     opt.PhoneNumber,
		CustomerUUID:    c.customerUUID,
		IP:              opt.IP,
		DeviceID:        opt.DeviceID,
		AppVersion:      opt.AppVersion,
		CallbackURL:     opt.CallbackURL,
		IsReturningUser: &opt.IsReturningUser,
	}

	if opt.DeviceType != nil {
		req.DeviceType = String(opt.DeviceType.String())
	}

	res, err := c.api.Authentication(ctx, req)
	if err != nil {
		return nil, apiErrToErr(err)
	}

	if res.Error != nil {
		return nil, apiErrorCodeToErr(res.Error.Code)
	}

	return &Authentication{
		AuthenticationUUID: res.Success.AuthenticationUUID,
		Status:             res.Success.Status,
		CreatedAt:          res.Success.CreatedAt,
		ExpiresAt:          res.Success.ExpiresAt,
	}, nil
}

// Authenticate performs an authentication request against the Ding API. Authentication
// requests allow you to send a message to a given phone number with a code that the user
// will have to enter in your app.
func (c *Client) Authenticate(opt AuthenticateOptions) (*Authentication, error) {
	return c.AuthenticateWithContext(context.Background(), opt)
}

// ----------------------------------------------------------------------------

// Check is the result of a check request.
type Check struct {
	AuthenticationUUID string
	Status             status.Check
}

// CheckWithContext performs a check request against the Ding API that can be cancelled
// with a context. Check requests allow you to enter the code that the user entered in
// your app to check if it is valid.
func (c *Client) CheckWithContext(ctx context.Context, authUUID string, code string) (*Check, error) {
	if !isValidUUID(authUUID) {
		return nil, ErrInvalidAuthUUID
	}

	res, err := c.api.Check(ctx, api.CheckRequest{
		CustomerUUID:       c.customerUUID,
		AuthenticationUUID: authUUID,
		CheckCode:          code,
	})
	if err != nil {
		return nil, apiErrToErr(err)
	}

	if res.Error != nil {
		return nil, apiErrorCodeToErr(res.Error.Code)
	}

	return &Check{
		AuthenticationUUID: res.Success.AuthenticationUUID,
		Status:             res.Success.Status,
	}, nil
}

// Check performs a check request against the Ding API.
// Check requests allow you to enter the code that the user entered in your app to check if it is valid.
func (c *Client) Check(authUUID string, code string) (*Check, error) {
	return c.CheckWithContext(context.Background(), authUUID, code)
}

// ----------------------------------------------------------------------------

// Retry is the result of a retry request.
type Retry struct {
	AuthenticationUUID string
	Status             status.Retry
}

// RetryWithContext performs a retry request against the Ding API that can be cancelled with
// a context.
func (c *Client) RetryWithContext(ctx context.Context, authUUID string) (*Retry, error) {
	if !isValidUUID(authUUID) {
		return nil, ErrInvalidAuthUUID
	}

	res, err := c.api.Retry(ctx, api.RetryRequest{
		CustomerUUID:       c.customerUUID,
		AuthenticationUUID: authUUID,
	})
	if err != nil {
		return nil, apiErrToErr(err)
	}

	if res.Error != nil {
		return nil, apiErrorCodeToErr(res.Error.Code)
	}

	return &Retry{
		AuthenticationUUID: res.Success.AuthenticationUUID,
		Status:             res.Success.Status,
	}, nil
}

// Retry performs a retry request against the Ding API. Retry requests allow you to send
// a new SMS to the user with a new code, using the initial authentication UUID.
func (c *Client) Retry(authUUID string) (*Retry, error) {
	return c.RetryWithContext(context.Background(), authUUID)
}

// ----------------------------------------------------------------------------

func apiErrToErr(err error) error {
	switch err {
	case api.ErrUnauthorized:
		return ErrUnauthorized
	default:
		return ErrInternal
	}
}

func apiErrorCodeToErr(code api.ErrorCode) error {
	switch code {
	case api.ErrorCodeInvalidPhoneNumber:
		return ErrInvalidPhoneNumber
	case api.ErrorCodeAccountInvalid:
		return ErrInvalidCustomerUUID
	case api.ErrorCodeNegativeBalance:
		return ErrNegativeBalance
	case api.ErrorCodeInvalidLine:
		return ErrInvalidPhoneNumber
	case api.ErrorCodeUnsupportedRegion:
		return ErrUnsupportedRegion
	case api.ErrorCodeInvalidAuthUUID:
		return ErrInvalidAuthUUID
	default:
		return ErrInternal
	}
}

// ----------------------------------------------------------------------------

func isValidNumber(phoneNumber string) bool {
	num, err := phonenumbers.Parse(phoneNumber, "ZZ")
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	return true
}

func isValidUUID(customerUUID string) bool {
	if _, err := uuid.Parse(customerUUID); err != nil {
		return false
	}

	return true
}

func isValidURL(u string) bool {
	if _, err := url.ParseRequestURI(u); err != nil {
		return false
	}

	return true
}
