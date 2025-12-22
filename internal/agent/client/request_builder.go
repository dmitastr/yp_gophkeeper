package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"gophkeep/internal/agent/compression"
)

type RequestBuilder struct {
	req        *retryablehttp.Request
	compressor compression.ICompressor
}

func NewRequestBuilder(method, baseURL, suffix string) (*RequestBuilder, error) {
	path, err := url.JoinPath(baseURL, suffix)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	return &RequestBuilder{req: req, compressor: compression.NewCompressor()}, nil
}

func (r *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	return r.WithRawBody(data)
}

func (r *RequestBuilder) WithRawBody(data []byte) *RequestBuilder {
	if r.compressor != nil {
		var err error
		data, err = r.compressor.Compress(data)
		if err != nil {
			panic(err)
		}
		r.req.Header.Set("Content-Encoding", "gzip")
	}

	r.req.Body = io.NopCloser(bytes.NewReader(data))
	r.req.ContentLength = int64(len(data))
	r.req.Header.Set("Content-Type", "application/json")

	r.req.SetBody(data)

	return r
}

func (r *RequestBuilder) WithBearer(token string) *RequestBuilder {
	r.req.Header.Add("Authorization", "Bearer "+token)
	return r
}

func (r *RequestBuilder) Build() *retryablehttp.Request {
	return r.req
}
