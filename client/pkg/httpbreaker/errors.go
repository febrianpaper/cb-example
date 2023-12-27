package httpbreaker

import "errors"

var (
	ErrUnkownVendor       = errors.New("unkown vendor")
	ErrCreateNewRequest   = errors.New("create new request failed")
	ErrSendingRequest     = errors.New("sending request failed")
	ErrReadBody           = errors.New("read body failed")
	ErrEmptyUrl           = errors.New("url is empty")
	ErrVendorNotAvailable = errors.New("vendor not available")
)
