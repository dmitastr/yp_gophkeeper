package validation

import "errors"

var ErrorSecretTypeNotFound = errors.New("secretType not exist")
var ErrorInvalidSecretContent = errors.New("invalid secret content")
var ErrorEmptySecretContent = errors.New("empty secret content")
