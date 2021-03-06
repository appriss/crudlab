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
	"strings"
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
	Header          *RequestHeader
	ContentType 	string
	Body            io.ReadCloser
	BodyValidator   func(content_type string, body []byte) error
	HeaderValidator func(http.Header) error
}

func (t CRUDTest) AddObjToBody(obj interface{}) CRUDTest {
	switch t.Header.M.Get("Content-Type") {
	case "application/json":
		b, _ := json.Marshal(obj)
		t.Body = ioutil.NopCloser(bytes.NewReader(b))
	case "application/xml":
		b, _ := xml.Marshal(obj)
		t.Body = ioutil.NopCloser(bytes.NewReader(b))
	default:
		fmt.Println("Debug, serializing could not determine content-type")
	}
	return t
}

type RequestHeader struct {
	M 	http.Header
}

func (h *RequestHeader) AddHeader(key, value string) *RequestHeader {
	h.M.Add(key, value)
	return h
}

func NewHeader() *RequestHeader {
	rq := &RequestHeader{}
	rq.M = make(http.Header)
	return rq
}

func ParseResponseBody(content_type string, body []byte, v interface{}) error {
	switch content_type {
	case "application/json":
		if err := json.Unmarshal(body, v); err != nil {
			fmt.Println("Error unserializing json")
			return err
		}
	case "application/xml":
		if err := xml.Unmarshal(body, v); err != nil {
			fmt.Println("Error unserializing xml")
			return err
		}
	default:
		return fmt.Errorf("Content type is not recognized, %s", content_type)
	}

	return nil
}

func HandlerTest(t []CRUDTest, h func(http.ResponseWriter, *http.Request)) error {
	for _, test := range t {
		rec := httptest.NewRecorder()
		if test.Header == nil {
			test.Header = NewHeader()
		}
		req := &http.Request{
			Method: test.Method,
			URL:    test.Resource,
			Header: test.Header.M,
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
					err = test.BodyValidator(test.Header.M.Get("Content-Type"), rec.Body.Bytes())
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

func VerifyHTTPMethods(methods []string, s string) bool {
	m := strings.Split(s, ",")
	for i,v := range m {
		m[i] = strings.TrimSpace(v)
	}
	for _,v := range methods {
		found := false
		for _, vv := range m {
			if vv == v {
				found = true 
				break
			}
		}
		if !found {
			return false
		}
	}	
	return true
}