package esi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eveisesi/athena"
	"github.com/volatiletech/null"
)

type options struct {
	method       string
	path         string
	query        url.Values
	headers      http.Header
	body         []byte
	retryOnError bool
	maxattempts  int
}

type OptionFunc func(*options) *options

func (s *service) opts(optionFuncs []OptionFunc) *options {
	opts := &options{
		maxattempts:  3,
		query:        make(url.Values),
		headers:      make(http.Header),
		body:         nil,
		retryOnError: true,
	}

	for _, optionFunc := range optionFuncs {
		opts = optionFunc(opts)
	}

	opts.headers.Set("user-agent", s.ua)
	opts.headers.Set("accept", "application/json")
	opts.headers.Set("content-type", "application/json; charset=UTF-8")

	return opts
}

func WithMethod(method string) OptionFunc {
	return func(o *options) *options {
		o.method = method
		return o
	}
}

func WithPath(path string) OptionFunc {
	return func(o *options) *options {
		o.path = path
		return o
	}
}

func WithQuery(key, value string) OptionFunc {
	return func(o *options) *options {
		if o.query == nil {
			o.query = url.Values{}
		}

		o.query.Set(key, value)

		return o
	}
}

func WithHeader(key, value string) OptionFunc {
	return func(o *options) *options {
		if o.headers == nil {
			o.headers = make(http.Header)
		}

		o.headers.Set(key, value)

		return o
	}
}

func WithBody(d []byte) OptionFunc {
	return func(o *options) *options {
		o.body = d
		return o
	}
}

// Helper Funcs below here cause i'm lazy. These will contain
// some sort of if check to ensure the value is not empty before
// calling one of the func above

func WithEtag(etag *athena.Etag) OptionFunc {

	if etag == nil || etag.Etag == "" {
		return emptyApplicator()
	}

	return WithHeader("if-none-match", etag.Etag)

}

func WithPage(i int) OptionFunc {
	return WithQuery("page", strconv.Itoa(i))
}

func WithAuthorization(token null.String) OptionFunc {
	if !token.Valid {
		return emptyApplicator()
	}

	return WithHeader("authorization", fmt.Sprintf("Bearer %s", token.String))
}

func emptyApplicator() OptionFunc {
	return func(o *options) *options {
		return o
	}
}
