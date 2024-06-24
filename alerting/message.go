package alerting

import (
	"fmt"
	"strings"
)

type MessageLevel string

const (
	DebugMessage   MessageLevel = "debug"
	InfoMessage    MessageLevel = "info"
	WarningMessage MessageLevel = "warning"
	ErrorMessage   MessageLevel = "error"
	FatalMessage   MessageLevel = "fatal"
)

type Message struct {
	Entity  string       `json:"entity"`
	Level   MessageLevel `json:"level"`
	Extras  []string     `json:"location"`
	Subject string       `json:"subject"`
	Details string       `json:"details"`
	Err     error        `json:"error"`
}

func NewMessage(entity string, level MessageLevel, extras []string, subject, details string, err error) *Message {
	return &Message{
		Entity:  entity,
		Level:   level,
		Extras:  extras,
		Subject: subject,
		Details: details,
		Err:     err,
	}
}

func (m *Message) ExtrasToString() string {
	return strings.Join(m.Extras, ", ")
}
func (m *Message) Validate() error {

	if m.Subject == "" {
		return ErrNoSubjectProvided
	}

	if len(m.Extras) == 0 {
		return fmt.Errorf("no message extras provided")
	}

	return nil
}
