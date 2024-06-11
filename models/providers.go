package models

import "fmt"

type ProviderEntity string
type ProviderConfig map[string]interface{}

// Errors
var (
	ErrProviderNotSupported = fmt.Errorf("provider not supported")
)
