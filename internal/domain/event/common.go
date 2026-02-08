package event

import (
	"time"

	"github.com/OddKuru/core-accounts/internal/domain/vo"
)

type (
	Type string

	Event interface {
		ID() vo.ID
		Type() Type
		Value() any
		Timestamp() time.Time
	}

	Events []Event
)
