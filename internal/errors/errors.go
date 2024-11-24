package errors

import (
	"github.com/pkg/errors"
)

var NotFound = errors.New("not found")
var BadRequest = errors.New("bad request")
