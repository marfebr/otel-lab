package service

import "errors"

// Erros específicos do serviço
var (
	ErrInvalidCEP = errors.New("invalid zipcode")
)
