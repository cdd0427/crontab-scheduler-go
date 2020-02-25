package common

import "errors"

var (
	ERR_LOCK_READY_REQUIRED = errors.New("Lock is already occupied")
)
