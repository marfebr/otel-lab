package service

import "errors"

// Erros específicos do serviço
var (
	ErrInvalidCEP  = errors.New("invalid zipcode")
	ErrCEPNotFound = errors.New("can not find zipcode")
)
