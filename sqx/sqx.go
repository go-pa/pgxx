package sqx

import (
	sq "github.com/Masterminds/squirrel"
)

// NewStatementBuilder returns a squirrel statement builder with postgres
// dollar placeholder format.
func NewStatementBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar)
}
