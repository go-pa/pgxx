package sqx

import (
	"context"
	"net/url"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func query() sq.SelectBuilder {
	return NewStatementBuilder().Select("*").From("table")
}

func TestWithPaging(t *testing.T) {
	ctx := context.Background()
	t.Run("default", func(t *testing.T) {
		sql, _, err := WithPaging(ctx, query(), url.Values{}, nil).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table LIMIT 100 OFFSET 0", sql)
	})
	t.Run("default size, page 10", func(t *testing.T) {
		sql, _, err := WithPaging(ctx, query(), url.Values{"page": {"10"}}, &PagingOptions{DefaultPageSize: 100}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table LIMIT 100 OFFSET 1000", sql)
	})
	t.Run("size 10, page 10", func(t *testing.T) {
		sql, _, err := WithPaging(ctx, query(), url.Values{"page": {"10"}, "size": {"10"}}, &PagingOptions{DefaultPageSize: 10}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table LIMIT 10 OFFSET 100", sql)
	})
}

func TestWithOrderBy(t *testing.T) {
	ctx := context.Background()
	t.Run("default", func(t *testing.T) {
		sql, _, err := WithOrderBy(ctx, query(), url.Values{}, "default", []OrderByMapping{}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table ORDER BY default", sql)
	})

	t.Run("default2", func(t *testing.T) {
		sql, _, err := WithOrderBy(ctx, query(), url.Values{}, "default",
			[]OrderByMapping{
				{
					Key:    "key1",
					Fields: []string{"c1", "c2"},
				},
			}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table ORDER BY default", sql)
	})

	t.Run("asc", func(t *testing.T) {
		sql, _, err := WithOrderBy(ctx, query(), url.Values{"order_by": {"key1,ignored"}}, "default",
			[]OrderByMapping{
				{
					Key:    "key1",
					Fields: []string{"c1", "c2"},
				},
			}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table ORDER BY c1 ASC, c2 ASC", sql)
	})

	t.Run("desc", func(t *testing.T) {
		sql, _, err := WithOrderBy(ctx, query(), url.Values{"order_by": {"-key1"}}, "default",
			[]OrderByMapping{
				{
					Key:    "key1",
					Fields: []string{"c1", "c2"},
				},
			}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table ORDER BY c1 DESC, c2 DESC", sql)
	})

	t.Run("multi", func(t *testing.T) {
		sql, _, err := WithOrderBy(ctx, query(), url.Values{"order_by": {"key1,-key2,key3"}}, "default",
			[]OrderByMapping{
				{
					Key:    "key1",
					Fields: []string{"c1", "c2"},
				},
				{
					Key:    "key2",
					Fields: []string{"c3"},
				},
				{
					Key:    "key3",
					Fields: []string{"c4", "c5"},
				},
			}).ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, "SELECT * FROM table ORDER BY c1 ASC, c2 ASC, c3 DESC, c4 ASC, c5 ASC", sql)
	})
}
