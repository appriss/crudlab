package crudlab

import (
	"net/http/httptest"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"fmt"
	"io"
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
	Body            io.ReadCloser
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
			Body: test.Body,
		}
		fmt.Printf("Running method %s on resource %s\n", req.Method, req.URL)
		h(rec, req)
		if rec.Code != test.ValidStatus {
			return fmt.Errorf("Run method %s on resource %s: expected %d but got %d.\n Request Object: %q\n Response: %q\n", 
				req.Method, req.URL, test.ValidStatus, rec.Code, req, rec)
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

func ( t CRUDTest) AddJSONBody( obj interface{} ) CRUDTest {
	b, _ := json.Marshal(obj)
	t.Body = ioutil.NopCloser(bytes.NewReader(b))
	if t.Header == nil {
		t.Header = make(http.Header)
	}
	t.Header.Add("Content-Type","application/json")
	return t
}

func ( t CRUDTest) AddXMLBody( obj interface{} ) CRUDTest {
	b, _ := xml.Marshal(obj)
	t.Body = ioutil.NopCloser(bytes.NewReader(b))
	if t.Header == nil {
		t.Header = make(http.Header)
	}
	t.Header.Add("Content-Type","application/xml")
	return t
}
