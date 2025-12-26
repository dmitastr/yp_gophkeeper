package test_data

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

// mock gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// mock getrequest
func MockJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = http.MethodGet
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}

// mock getrequest
func MockJsonDelete(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = http.MethodDelete
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}

// mock postrequest
func MockJsonPost(c *gin.Context, params gin.Params, content any, jsonEncoded bool) {
	c.Request.Method = http.MethodPost
	c.Set("user_id", 1)
	c.Params = params

	if jsonEncoded {
		c.Request.Header.Set("Content-Type", "application/json")

		jsonBytes, err := json.Marshal(content)
		if err != nil {
			panic(err)
		}

		// set query params
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		return
	}

	body, _ := content.([]byte)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

}
