package status

import (
	"encoding/json"
	"fmt"
)

type Check string

const (
	CheckUnknown Check = "unknown"

	CheckValid            Check = "valid"
	CheckInvalid          Check = "invalid"
	CheckWithoutAttempt   Check = "without_attempt"
	CheckRateLimited      Check = "rate_limited"
	CheckAlreadyValidated Check = "already_validated"
	CheckExpiredAuth      Check = "expired_auth"
)

func (d *Check) UnmarshalJSON(b []byte) error {
	var res string
	err := json.Unmarshal(b, &res)
	if err != nil {
		return fmt.Errorf("unmarshal Check: %w", err)
	}

	switch res {
	case "valid", "invalid", "without_attempt", "rate_limited", "already_validated", "expired_auth":
		*d = Check(res)
	default:
		*d = CheckUnknown
	}

	return nil
}

// ----------------------------------------------------------------------------

type Auth string

const (
	AuthUnknown Auth = "unknown"

	AuthPending      Auth = "pending"
	AuthRateLimited  Auth = "rate_limited"
	AuthSpamDetected Auth = "spam_detected"
	AuthApproved     Auth = "approved"
	AuthCanceled     Auth = "canceled"
	AuthExpired      Auth = "expired"
)

func (d *Auth) UnmarshalJSON(b []byte) error {
	var res string
	err := json.Unmarshal(b, &res)
	if err != nil {
		return fmt.Errorf("unmarshal Auth: %w", err)
	}

	switch res {
	case "pending", "rate_limited", "spam_detected", "approved", "canceled", "expired":
		*d = Auth(res)
	default:
		*d = AuthUnknown
	}

	return nil
}

// ----------------------------------------------------------------------------

type Retry string

const (
	RetryUnknown Retry = "unknown"

	RetryApproved         Retry = "approved"
	RetryDenied           Retry = "denied"
	RetryNoAttempt        Retry = "no_attempt"
	RetryRateLimited      Retry = "rate_limited"
	RetryExpiredAuth      Retry = "expired_auth"
	RetryAlreadyValidated Retry = "already_validated"
)

func (d *Retry) UnmarshalJSON(b []byte) error {
	var res string
	err := json.Unmarshal(b, &res)
	if err != nil {
		return fmt.Errorf("unmarshal Retry: %w", err)
	}

	switch res {
	case "approved", "denied", "no_attempt", "rate_limited", "expired_auth", "already_validated":
		*d = Retry(res)
	default:
		*d = RetryUnknown
	}

	return nil
}
