package models

import "errors"

var ErrorInvalidBankCardNumber = errors.New("card number is invalid")
var ErrorEmptyAuthData = errors.New("username and password cannot be empty")
var ErrorSecretTypeNotFound = errors.New("secretType not exist")
var ErrorUserNotFound = errors.New("user not found")
var ErrorInvalidPassword = errors.New("invalid password")
