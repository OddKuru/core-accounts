package rco_test

import (
	"testing"

	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	"github.com/OddKuru/core-accounts/internal/domain/rco"
	"github.com/stretchr/testify/assert"
)

func TestQueryOrderValidate(t *testing.T) {
	t.Run("Should pass for ASC", func(t *testing.T) {
		err := rco.Asc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Should pass for DESC", func(t *testing.T) {
		err := rco.Desc.Validate()
		assert.NoError(t, err)
	})

	t.Run("Should error for other values", func(t *testing.T) {
		order := rco.QueryOrder("OTHER")
		err := order.Validate()
		assert.Error(t, err)
	})
}

func TestNewQuery(t *testing.T) {
	t.Run("Should construct query with ASC", func(t *testing.T) {
		q, err := rco.NewQuery(2, 10, "created_at", rco.Asc)
		assert.NoError(t, err)

		assert.Equal(t, uint(2), q.Page())
		assert.Equal(t, uint(10), q.Limit())
		assert.Equal(t, "created_at", q.SortBy())
		assert.Equal(t, rco.Asc, q.SortOrder())
		assert.Equal(t, uint(10), q.PaginationOffset())
	})

	t.Run("Should construct query with DESC", func(t *testing.T) {
		q, err := rco.NewQuery(3, 5, "id", rco.Desc)
		assert.NoError(t, err)

		assert.Equal(t, uint(3), q.Page())
		assert.Equal(t, uint(5), q.Limit())
		assert.Equal(t, "id", q.SortBy())
		assert.Equal(t, rco.Desc, q.SortOrder())
	})

	t.Run("Should error with invalid order", func(t *testing.T) {
		order := rco.QueryOrder("OTHER")
		q, err := rco.NewQuery(1, 20, "id", order)
		assert.Error(t, err)
		assert.Equal(t, rco.Query{}, q)
		assert.Equal(t, codex.InvalidArgument, errx.Code(err))
	})
}

func TestPaginationOffset(t *testing.T) {
	order := rco.Asc

	q, err := rco.NewQuery(1, 10, "id", order)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), q.PaginationOffset())

	q, err = rco.NewQuery(3, 25, "id", order)
	assert.NoError(t, err)
	assert.Equal(t, uint(50), q.PaginationOffset())
}

func TestDataWithPageCount(t *testing.T) {
	t.Run("Should hold data and page count with slice", func(t *testing.T) {
		data := []int{1, 2, 3}
		dwp := rco.NewDataWithPageCount(data, 5)

		assert.Equal(t, data, dwp.Data())
		assert.Equal(t, uint(5), dwp.PageCount())
	})

	t.Run("Should work with struct type", func(t *testing.T) {
		type item struct {
			ID   int
			Name string
		}

		data := []item{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
		}

		dwp := rco.NewDataWithPageCount(data, 2)
		assert.Equal(t, data, dwp.Data())
		assert.Equal(t, uint(2), dwp.PageCount())
	})
}
