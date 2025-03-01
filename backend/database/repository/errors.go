package repository

import "errors"

var ErrNotFound = errors.New("record not found")
var ErrInvalidTransaction = errors.New("invalid transaction")
