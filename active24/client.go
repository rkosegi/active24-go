/*
Copyright 2023 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package active24

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SandboxToken = "123456qwerty-ok"
)

type Option func(c *client)

func ApiEndpoint(ep string) Option {
	return func(c *client) {
		c.h.apiEndpoint = ep
	}
}

func Timeout(to time.Duration) Option {
	return func(c *client) {
		c.h.c.Timeout = to
	}
}

type ApiError interface {
	Error() error
	Response() *http.Response
}

type Client interface {
	//Dns provides interface to interact with DNS records
	Dns() Dns
	//Domains provides interface to interact with domains
	Domains() Domains
}

func New(apiKey string, opts ...Option) Client {
	c := &client{
		h: helper{
			apiEndpoint: "https://api.active24.com",
			auth:        apiKey,
			c: http.Client{
				Timeout: time.Second * 10,
			},
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type client struct {
	h helper
}

func (c *client) Dns() Dns {
	return &dns{
		h: c.h,
	}
}

func (c *client) Domains() Domains {
	return &domains{
		h: c.h,
	}
}

type apiError struct {
	err  error
	resp *http.Response
}

func (a *apiError) Error() error {
	return a.err
}

func (a *apiError) Response() *http.Response {
	return a.resp
}

type helper struct {
	apiEndpoint string
	auth        string
	c           http.Client
}

func (ch *helper) do(method string, suffix string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, fmt.Sprintf("%s/%s", ch.apiEndpoint, suffix), body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ch.auth))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "*/*")
	return ch.c.Do(r)
}

func apiErr(resp *http.Response, err error) ApiError {
	if err == nil {
		return nil
	}
	return &apiError{
		err:  err,
		resp: resp,
	}
}
