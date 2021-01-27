package esi

import (
	"fmt"
	"net/http"
	"net/url"
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

// MODIFY THIS TO TAKE A SLICE OF ETAG and a page
func WithEtag(etag string) OptionFunc {
	if etag == "" {
		return emptyApplicator()
	}

	return WithHeader("if-none-match", etag)

}

func WithAuthorization(token string) OptionFunc {
	if token == "" {
		return emptyApplicator()
	}

	return WithHeader("authorization", fmt.Sprintf("Bearer %s", token))
}

func emptyApplicator() OptionFunc {
	return func(o *options) *options {
		return o
	}
}
