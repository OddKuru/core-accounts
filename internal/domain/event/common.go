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

type TestEvent struct {
	Id         vo.ID
	TypeF      Type
	ValueF     any
	TimestampF time.Time
}

func (t TestEvent) ID() vo.ID {
	return t.Id
}

func (t TestEvent) Type() Type {
	return t.TypeF
}

func (t TestEvent) Value() any {
	return t.ValueF
}

func (t TestEvent) Timestamp() time.Time {
	return t.TimestampF
}

func (e Events) Len() int {
	return len(e)
}

func (e Events) Filter(filterFunc func(elem Event) bool) Events {
	newEvents := make(Events, 0)
	for _, ev := range e {
		if filterFunc(ev) {
			newEvents = append(newEvents, ev)
		}
	}
	return newEvents
}

func (e Events) LogData() any {
	sl := make([]string, 0, len(e))
	for _, ev := range e {
		sl = append(sl, ev.ID().Value())
	}
	return sl
}
