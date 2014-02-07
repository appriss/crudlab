package crudlab

import "net/http"
import "errors"
// import "net/http/httptest"
import "fmt"
import "testing"

func goodHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This request succeeded!")
}

func badHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This request failed!", http.StatusInternalServerError)
}

func goodBodyValidator([]byte) error {
	return nil
}

func badBodyValidator( []byte ) error {
	return errors.New("Body Validation Failed")
}

func goodHeaderValidator( http.Header ) error {
	return nil
}

func badHeaderValidator( http.Header ) error {
	return errors.New("Header Validation Failed")
}

type testableResponse struct {
}

func TestResponseProper(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200}}
	err := HandlerTest(mlist, goodHandler)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestResponseError(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 500}}
	err := HandlerTest(mlist, badHandler)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestResponseNoErrorMatch(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200}}
	err := HandlerTest(mlist, badHandler)
	if err == nil {
		t.Error("Expected failure but got success")
	}
}

func TestUseBodyValidator(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200, BodyValidator: goodBodyValidator }}
	err := HandlerTest(mlist, goodHandler)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestBodyValidationFailed(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200, BodyValidator: badBodyValidator }}
	err := HandlerTest(mlist, goodHandler)
	if err == nil {
		t.Error(err.Error())
	}
}

func TestUseHeaderValidator(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200, HeaderValidator: goodHeaderValidator }}
	err := HandlerTest(mlist, goodHandler)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestHeaderValidationFailed(t *testing.T) {
	mlist := []CRUDTest{CRUDTest{Method: "GET", ValidStatus: 200, HeaderValidator: badHeaderValidator }}
	err := HandlerTest(mlist, goodHandler)
	if err == nil {
		t.Error(err.Error())
	}
}
