package sqx

import (
	"fmt"
	"net/url"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// TODO: options

// OrderByMapping defines the relation between the order_by url parameter and database ordering.
type OrderByMapping struct {
	Key    string
	Fields []string
}

func WithOrderBy(sb sq.SelectBuilder, q url.Values, defaultOrdering string, mapping []OrderByMapping) sq.SelectBuilder {
	if q.Get("order_by") == "" {
		if defaultOrdering == "" {
			return sb
		}
		return sb.OrderBy(defaultOrdering)
	}
	var orderBy []string
	for _, orderKey := range strings.Split(q.Get("order_by"), ",") {
		dir := "ASC"
		if strings.HasPrefix(orderKey, "-") {
			orderKey = strings.TrimPrefix(orderKey, "-")
			dir = "DESC"
		}
		for _, v := range mapping {
			if v.Key == orderKey {
				for _, f := range v.Fields {
					orderBy = append(orderBy, fmt.Sprintf("%s %s", f, dir))
				}
			}
		}
	}
	if len(orderBy) > 0 {
		return sb.OrderBy(orderBy...)
	}
	return sb
}
