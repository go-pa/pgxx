package sqx

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-pa/pgxx/logging"
)

// PagingOptions.
// To completley disable logging it is recommended to set
type PagingOptions struct {
	DefaultPageSize     int64
	MaxPageSize         int64 // set to -1 to disable max page size
	PageSizeParameter   string
	PageNumberParameter string
	Logger              logging.LoggerFunc
}

var DefaultPagingOptions = PagingOptions{
	DefaultPageSize:     100,
	MaxPageSize:         1000,
	PageSizeParameter:   "size",
	PageNumberParameter: "page",
	Logger:              logging.Get,
}

// WithPaging
func WithPaging(ctx context.Context, sb sq.SelectBuilder, q url.Values, opts *PagingOptions) sq.SelectBuilder {
	logger := opts.getLogger(ctx)
	size := opts.getDefaultPageSize()
	pageSizeParam := opts.getPageSizeParameter()
	if q.Get(pageSizeParam) != "" {
		s := q.Get(pageSizeParam)
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil && logger != nil {
			fmt.Fprintf(logger, "could not parse page size paramenter %s: %v", s, err)
		} else {
			size = n
		}
	}
	maxSize := opts.getMaxPageSize()
	if maxSize > 0 && size > maxSize {
		size = maxSize
	}
	sb = sb.Limit(uint64(size))
	pageParam := opts.getPageNumberParameter()
	var page int64
	if q.Get(pageParam) != "" {
		s := q.Get(pageParam)
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			// TODO: used to be application specific logging here
		} else {
			page = n
		}
	}
	sb = sb.Offset(uint64(page * size))
	return sb
}

func (p *PagingOptions) getDefaultPageSize() int64 {
	if p == nil || p.DefaultPageSize < 1 {
		return DefaultPagingOptions.DefaultPageSize
	}
	return p.DefaultPageSize
}

func (p *PagingOptions) getMaxPageSize() int64 {
	if p == nil || p.MaxPageSize == 0 {
		return DefaultPagingOptions.MaxPageSize
	}
	return p.DefaultPageSize
}

func (p *PagingOptions) getPageSizeParameter() string {
	if p == nil || p.PageSizeParameter == "" {
		return DefaultPagingOptions.PageSizeParameter
	}
	return p.PageSizeParameter
}

func (p *PagingOptions) getPageNumberParameter() string {
	if p == nil || p.PageNumberParameter == "" {
		return DefaultPagingOptions.PageNumberParameter
	}
	return p.PageNumberParameter
}

func (p *PagingOptions) getLogger(ctx context.Context) io.Writer {
	if p == nil || p.Logger == nil {
		if DefaultPagingOptions.Logger == nil {
			return nil
		}
		return DefaultPagingOptions.Logger(ctx)
	}
	return p.Logger(ctx)
}
