package event_test

import (
	"testing"
	"time"

	"github.com/OddKuru/core-accounts/internal/domain/event"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEvents(t *testing.T) {
	voID1, _ := vo.NewID(uuid.New().String())
	voID2, _ := vo.NewID(uuid.New().String())
	voID3, _ := vo.NewID(uuid.New().String())
	var typeVal event.Type = "test-type"
	val := "exdee"
	ti := time.Now()

	ev1 := event.TestEvent{
		Id:         voID1,
		TypeF:      typeVal,
		ValueF:     val,
		TimestampF: ti,
	}

	assert.Equal(t, ev1.ID().Value(), voID1.Value())
	assert.Equal(t, ev1.Type(), typeVal)
	assert.Equal(t, ev1.Value(), any(val))
	assert.Equal(t, ev1.Timestamp(), ti)

	ev := event.Events{
		ev1,
		event.TestEvent{
			Id:         voID2,
			TypeF:      typeVal,
			ValueF:     val,
			TimestampF: ti,
		},
		event.TestEvent{
			Id:         voID3,
			TypeF:      typeVal,
			ValueF:     val,
			TimestampF: ti,
		},
	}

	assert.Equal(t, 3, ev.Len())
	assert.Equal(t, []string{voID1.Value(), voID2.Value(), voID3.Value()}, ev.LogData())
	assert.Equal(t, []string{voID1.Value(), voID2.Value()}, ev.Filter(func(ev event.Event) bool {
		return ev.ID().Value() != voID3.Value()
	}).LogData())
	assert.Equal(t, []string{voID1.Value(), voID3.Value()}, ev.Filter(func(ev event.Event) bool {
		return ev.ID().Value() != voID2.Value()
	}).LogData())
}
