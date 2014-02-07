package crudlab

import (
	"fmt"
	// "log"
	"net/http"
	"net/http/httptest"
	"net/url"
	// "io"
	//"regexp"
)

type LabRoom struct {
}

type TestRunner interface {
	RunTest() error
}

type CRUDTest struct {
	Resource        *url.URL
	Method          string
	ValidStatus     int
	Header          http.Header
	BodyValidator   func([]byte) error
	HeaderValidator func(http.Header) error
}

func HandlerTest(t []CRUDTest, h func(http.ResponseWriter, *http.Request)) error {
	for _, test := range t {
		rec := httptest.NewRecorder()
		req := &http.Request{
			Method: test.Method,
			URL:    test.Resource,
			Header: test.Header,
		}
		h(rec, req)
		if rec.Code != test.ValidStatus {
			return fmt.Errorf("Run method %s on resource %s: expected %d but got %d", req.Method, req.URL, test.ValidStatus, rec.Code)
		}
		var err error

		for _, v := range []func(){
			func() {
				if test.HeaderValidator != nil {
					err = test.HeaderValidator(rec.HeaderMap)
				}
			},
			func() {
				if test.BodyValidator != nil {
					err = test.BodyValidator(rec.Body.Bytes())
				}
			},
		} {
			v()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
