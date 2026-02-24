package psql

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/aggregate"
	"github.com/OddKuru/core-accounts/internal/domain/entity"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/OddKuru/core-accounts/internal/domain/repository"
	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var _ repository.AccountQuery = (*AccountQuery)(nil)

const (
	AccountSliceCap = 16
)

type AccountQuery struct {
	pool *pgxpool.Pool
}

func NewAccountQuery(pool *pgxpool.Pool) (*AccountQuery, error) {
	return &AccountQuery{
		pool: pool,
	}, nil
}

func (a *AccountQuery) HasByEmail(ctx context.Context, email vo.Email) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(ctx, AccountHasByEmailQuery, email.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) HasByName(ctx context.Context, name vo.LoginName) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(ctx, AccountHasByNameQuery, name.Value()).Scan(&exists)
	if err != nil {
		return false, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}
	return exists, nil
}

func (a *AccountQuery) GetById(ctx context.Context, id vo.ID) (*aggregate.Account, error) {
	var (
		name     string
		email    string
		password string
		role     string
		version  uint
		created  time.Time
		updated  time.Time
	)

	err := a.pool.QueryRow(ctx, AccountGetByIdQuery, id.Value()).
		Scan(&name, &email, &password, &version, &created, &updated, &role)
	if err != nil {
		return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}

	voName, err := vo.NewLoginName(name)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewLoginName")
	}
	voEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewEmail")
	}
	voPassword := vo.NewHashedPassword([]byte(password))
	voRole, err := vo.NewRole(RoleValueFromDB[role])
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] vo.NewRole")
	}

	accEntity, err := entity.NewAccount(id, voName, voEmail, voPassword, voRole, version, created, updated)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] entity.NewAccount")
	}
	agg, err := aggregate.NewAccount(accEntity)
	if err != nil {
		return nil, errors.Wrap(err, "[AccountQuery] aggregate.NewAccount")
	}

	return agg, nil
}

var allowedSortFields = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"login":      "login",
	"email":      "email",
}

func (a *AccountQuery) GetByQuery(
	ctx context.Context, query rco.Query,
) (*rco.DataWithPageCount[[]*aggregate.Account], error) {
	limit := query.Limit()
	if limit == 0 {
		limit = 10
	}

	page := query.Page()
	if page == 0 {
		page = 1
	}

	sortBy := "created_at"
	if v, ok := allowedSortFields[query.SortBy()]; ok {
		sortBy = v
	}

	sortOrder := "DESC"
	if query.SortOrder() == rco.Asc {
		sortOrder = "ASC"
	}

	var total uint
	if err := a.pool.QueryRow(ctx, AccountCountQuery).Scan(&total); err != nil {
		return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.QueryRow")
	}

	pageCount := uint(math.Ceil(float64(total) / float64(limit)))

	if page > pageCount {
		page = pageCount
	}

	offset := (page - 1) * limit

	dataQuery := fmt.Sprintf(`
		SELECT a.id, a.login, a.email, a.password_hash, a.version,
		       r.value, a.updated_at, a.created_at
		FROM accounts a
		JOIN roles r ON r.id = a.role_id
		WHERE deleted_at IS NULL
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	rows, err := a.pool.Query(ctx, dataQuery, limit, offset)
	if err != nil {
		return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] pool.Query")
	}
	defer rows.Close()

	accounts := make([]*aggregate.Account, 0, AccountSliceCap)

	for rows.Next() {
		var (
			id        string
			login     string
			email     string
			password  string
			version   uint
			role      string
			createdAt time.Time
			updatedAt time.Time
		)
		err := rows.Scan(
			&id,
			&login,
			&email,
			&password,
			&version,
			&role,
			&updatedAt,
			&createdAt,
		)
		if err != nil {
			return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] rows.Scan")
		}

		voID, err := vo.NewID(id)
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] vo.NewID")
		}

		voLogin, err := vo.NewLoginName(login)
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] vo.NewLoginName")
		}
		voEmail, err := vo.NewEmail(email)
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] vo.NewEmail")
		}
		voPassword := vo.NewHashedPassword([]byte(password))
		voRole, err := vo.NewRole(RoleValueFromDB[role])
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] vo.NewRole")
		}

		accEntity, err := entity.NewAccount(
			voID,
			voLogin,
			voEmail,
			voPassword,
			voRole,
			version,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] entity.NewAccount")
		}

		agg, err := aggregate.NewAccount(accEntity)
		if err != nil {
			return nil, errors.Wrap(err, "[AccountQuery] aggregate.NewAccount")
		}

		accounts = append(accounts, agg)
	}

	if err := rows.Err(); err != nil {
		return nil, errx.WrapWithCode(err, codex.Internal, "[AccountQuery] rows.Err")
	}

	return rco.NewDataWithPageCount(accounts, pageCount), nil
}
