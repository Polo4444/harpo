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
	Entity   string       `json:"entity"`
	Level    MessageLevel `json:"level"`
	Location []string     `json:"location"`
	Subject  string       `json:"subject"`
	Details  error        `json:"details"`
}

func NewMessage(entity string, level MessageLevel, location []string, subject string, details error) *Message {
	return &Message{Entity: entity, Level: level, Location: location, Subject: subject, Details: details}
}

func (m *Message) LocationToString() string {
	return strings.Join(m.Location, " > ")
}
func (m *Message) Validate() error {

	if m.Subject == "" {
		return ErrNoSubjectProvided
	}

	if len(m.Location) == 0 {
		return fmt.Errorf("no message location provided")
	}

	return nil
}
