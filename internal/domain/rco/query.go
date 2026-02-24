package rco

import (
	"github.com/OddEer0/errx"
	"github.com/OddEer0/errx/codex"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	Asc  QueryOrder = "ASC"
	Desc QueryOrder = "DESC"
)

type (
	QueryOrder string

	Query struct {
		page      uint
		limit     uint
		sortBy    string
		sortOrder QueryOrder
	}

	DataWithPageCount[DataType any] struct {
		data      DataType
		pageCount uint
	}
)

func (q QueryOrder) Validate() error {
	return validation.Validate(string(q), validation.Required,
		validation.In(string(Asc), string(Desc)))
}

func (q Query) Page() uint {
	return q.page
}

func (q Query) Limit() uint {
	return q.limit
}

func (q Query) SortBy() string {
	return q.sortBy
}

func (q Query) SortOrder() QueryOrder {
	return q.sortOrder
}

func NewQuery(page, limit uint, sortBy string, sortOrder QueryOrder) (Query, error) {
	if err := sortOrder.Validate(); err != nil {
		return Query{}, errx.WrapWithCode(err, codex.InvalidArgument, "Query.Validate")
	}
	return Query{
		page:      page,
		limit:     limit,
		sortBy:    sortBy,
		sortOrder: sortOrder,
	}, nil
}

func (q Query) PaginationOffset() uint {
	return (q.page - 1) * q.limit
}

func NewDataWithPageCount[T any](data T, pageCount uint) *DataWithPageCount[T] {
	return &DataWithPageCount[T]{
		data:      data,
		pageCount: pageCount,
	}
}

func (d DataWithPageCount[DataType]) Data() DataType {
	return d.data
}

func (d DataWithPageCount[DataType]) PageCount() uint {
	return d.pageCount
}
