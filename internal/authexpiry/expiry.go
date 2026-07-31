package authexpiry

import (
	"errors"
	"time"
)

// ErrExpired is returned when a token expiry is not after the current time.
var ErrExpired = errors.New("authexpiry: token expired")

// ValidateExpiry fixes the refresh path bug where JWT expiry was not validated.
func ValidateExpiry(exp time.Time, now time.Time) error {
	if exp.IsZero() {
		return errors.New("authexpiry: missing expiry")
	}
	if !exp.After(now) {
		return ErrExpired
	}
	return nil
}
