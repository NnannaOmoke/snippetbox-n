package models

import "errors"

// ErrNoRecord shows that the record does not exist in the database
var ErrNoRecord = errors.New("Models: No matching record found")
var ErrInvalidCredentails = errors.New("Models: Invalid Credentials")
var ErrDuplicateEmail = errors.New("Models: Duplicate Email")
