package aggregate

import (
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/event"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/pkg/errors"
)

const (
	DefaultEventCapacity = 8
)

type Account struct {
	account *entity.Account
	events  event.Events
}

func NewAccount(acc *entity.Account) (*Account, error) {
	if err := acc.Validate(); err != nil {
		return nil, errx.WrapWithCode(err, codex.InvalidArgument, "[Account] account.Validate")
	}
	events := make([]event.Event, 0, DefaultEventCapacity)

	return &Account{
		account: acc,
		events:  events,
	}, nil
}

func (a *Account) Account() entity.AccountViewer {
	return a.account.View()
}

func (a *Account) Events() event.Events {
	return a.events
}

func (a *Account) HasEvents() bool {
	return len(a.events) > 0
}

func (a *Account) AddEvents(events ...event.Event) {
	a.events = append(a.events, events...)
}

func (a *Account) ClearEvent() {
	a.events = a.events[:0]
}

func (a *Account) ChangeName(name vo.LoginName, t time.Time) error {
	if name.Value() == a.Account().Name().Value() {
		return nil
	}

	if err := a.account.ChangeName(name); err != nil {
		return errors.Wrap(err, "[Account] account.SetName")
	}
	ev := event.NewAccountChangeName(a.account.ID(), a.account.Name(), t)
	a.AddEvents(ev)

	return nil
}

func (a *Account) ChangeEmail(email vo.Email, t time.Time) error {
	if email.Value() == a.Account().Email().Value() {
		return nil
	}

	if err := a.account.ChangeEmail(email); err != nil {
		return errors.Wrap(err, "[Account] account.SetEmail")
	}
	ev := event.NewAccountChangeEmail(a.account.ID(), a.account.Email(), t)
	a.AddEvents(ev)

	return nil
}

func (a *Account) ChangeRole(role vo.Role, t time.Time) error {
	if role.Value() == a.Account().Role().Value() {
		return nil
	}

	if err := a.account.ChangeRole(role); err != nil {
		return errors.Wrap(err, "[Account] account.SetRole")
	}
	ev := event.NewAccountChangeRole(a.account.ID(), a.account.Role(), t)
	a.AddEvents(ev)

	return nil
}

func (a *Account) ChangePassword(password vo.HashedPassword, t time.Time) error {
	if err := a.account.ChangePassword(password); err != nil {
		return errors.Wrap(err, "[Account] account.SetPassword")
	}
	ev := event.NewAccountChangePassword(a.Account().ID(), t)
	a.AddEvents(ev)

	return nil
}
