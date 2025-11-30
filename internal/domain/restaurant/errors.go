package restaurant

import "errors"

var (
	ErrInvalidAddress     = errors.New("Invalid address")
	ErrInvalidWebsite     = errors.New("Invalid website")
	ErrInvalidLogo        = errors.New("Invalid logo")
	ErrInvalidCity        = errors.New("Invalid city")
	ErrInvalidBanner      = errors.New("Invalid banner")
	ErrInvalidDistrict    = errors.New("Invalid district")
	ErrInvalidPhoneNumber = errors.New("Invalid phone number")
)
