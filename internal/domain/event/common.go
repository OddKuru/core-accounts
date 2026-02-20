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

func (e Events) LogData() any {
	sl := make([]string, 0, len(e))
	for _, ev := range e {
		sl = append(sl, ev.ID().Value())
	}
	return sl
}
