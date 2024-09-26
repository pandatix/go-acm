package api

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type ACMClient struct {
	sub *http.Client
}

func NewACMClient() *ACMClient {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	return &ACMClient{
		sub: &http.Client{
			Jar: jar,
		},
	}
}

type options struct {
	ctx context.Context
}

type Option interface {
	apply(*options)
}

type ctxOption struct {
	ctx context.Context
}

func (opt *ctxOption) apply(opts *options) {
	opts.ctx = opt.ctx
}

func WithContext(ctx context.Context) Option {
	return &ctxOption{
		ctx: ctx,
	}
}

func (client *ACMClient) call(req *http.Request, opts ...Option) (res *http.Response, err error) {
	// Apply functional options
	o := &options{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt.apply(o)
	}

	// Bind options and hardcoded values to request
	req = req.WithContext(o.ctx)
	req.URL, err = url.Parse("https://dl.acm.org" + req.URL.String())
	if err != nil {
		return
	}

	// Execute request
	return client.sub.Do(req)
}
