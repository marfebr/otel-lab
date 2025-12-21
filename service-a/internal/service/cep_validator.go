package service

import (
	"regexp"
	"strings"
)

// CEPValidator interface para validação de CEP
type CEPValidator interface {
	ValidateCEP(cep string) error
}

// cepValidator implementação da validação de CEP
type cepValidator struct{}

// NewCEPValidator cria uma nova instância do validador de CEP
func NewCEPValidator() CEPValidator {
	return &cepValidator{}
}

// ValidateCEP valida se o CEP está no formato correto
// Formato esperado: 8 dígitos numéricos
func (v *cepValidator) ValidateCEP(cep string) error {
	// Remove espaços em branco
	cep = strings.TrimSpace(cep)

	// Verifica se está vazio
	if cep == "" {
		return ErrInvalidCEP
	}

	// Verifica se tem exatamente 8 caracteres
	if len(cep) != 8 {
		return ErrInvalidCEP
	}

	// Verifica se contém apenas dígitos
	matched, err := regexp.MatchString(`^[0-9]{8}$`, cep)
	if err != nil {
		return ErrInvalidCEP
	}

	if !matched {
		return ErrInvalidCEP
	}

	return nil
}
