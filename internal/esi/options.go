package esi

import (
	"bytes"
	"net/http"
	"net/url"
)

type options struct {
	// This is the User Agent to use in our HTTP Request
	method       string
	path         string
	query        url.Values
	headers      http.Header
	body         *bytes.Buffer
	retryOnError bool
	maxattempts  int
}

type OptionsFunc func(*options) *options

func (s *service) opts(optionFuncs []OptionsFunc) *options {
	opts := &options{
		maxattempts: 3,
	}

	for _, optionFunc := range optionFuncs {
		opts = optionFunc(opts)
	}

	opts.headers.Set("user-agent", s.ua)
	opts.headers.Set("accept", "application/json")
	opts.headers.Set("content-type", "application/json; charset=UTF-8")

	return opts
}

func WithMethod(method string) OptionsFunc {
	return func(o *options) *options {
		o.method = method
		return o
	}
}

func WithPath(path string) OptionsFunc {
	return func(o *options) *options {
		o.path = path
		return o
	}
}

func WithQuery(key, value string) OptionsFunc {
	return func(o *options) *options {
		if o.query == nil {
			o.query = url.Values{}
		}

		o.query.Set(key, value)

		return o
	}
}

func WithHeaders(key, value string) OptionsFunc {
	return func(o *options) *options {
		if o.headers == nil {
			o.headers = http.Header{}
		}

		o.headers.Set(key, value)

		return o
	}
}

func WithBody(d []byte) OptionsFunc {
	return func(o *options) *options {
		o.body = bytes.NewBuffer(d)
		return o
	}
}

// Helper Funcs Below here cause i'm lazy. These will contain
// some sort of if check to ensure the value is not empty before
// calling one of the func above

func WithEtag(etag string) OptionsFunc {
	if etag == "" {
		return emptyApplicator()
	}

	return WithHeaders("if-none-match", etag)

}

func emptyApplicator() OptionsFunc {
	return func(o *options) *options {
		return o
	}
}
