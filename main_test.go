package main

import (
	/*"bytes"
	"encoding/binary"
	"io"*/
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*func Testupload(t *testing.T) {

	req, err := http.NewRequest("POST", "/upload", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(upload)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	expected := []byte(`{"status": 500, "error": "Bad connection"}`)
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.Bytes(), expected)
	}
}*/

func TestGet(t *testing.T) {

	/*var rid io.Reader

	var bb bytes.Buffer
	binary.Write(&bb, binary.BigEndian, "Layout")
	rid.Read(bb.Bytes())*/

	var example string
	var correctRes string

	for i := 0; i < 3; i++ {

		switch i {
		case 0:
			example = "States"
			correctRes = `[{"id":"generated","label":"Generated"},{"id":"accepted","label":"Accepted"},{"id":"approved","label":"Approved"},{"id":"implemented","label":"Implemented"},{"id":"benefit_realized","label":"Benefit Realized"},{"id":"rejected","label":"Rejected"}]
`
		case 1:
			example = "Transitions"
			correctRes = `[{"id":"generated-accepted","label":"Accept","from":"generated","to":"accepted","direction":"next"},{"id":"accepted-generated","label":"Back to generated","from":"accepted","to":"generated","direction":"previous"},{"id":"accepted-approved","label":"Approve","from":"accepted","to":"approved","direction":"next"},{"id":"approved-accepted","label":"Back to Accepting","from":"approved","to":"accepted","direction":"previous"},{"id":"approved-implemented","label":"Implement","from":"approved","to":"implemented","direction":"next"},{"id":"implemented-approved","label":"Back to Approving","from":"implemented","to":"approved","direction":"previous"},{"id":"implemented-benefit_realized","label":"Benefit Realize","from":"implemented","to":"benefit_realized","direction":"next"},{"id":"benefit_realized-implemented","label":"Back to Implementation","from":"benefit_realized","to":"implemented","direction":"previous"},{"id":"rejected-generated","label":"Move to Generated","from":"rejected","to":"generated","direction":"next"}]
`
		case 2:
			example = "Layout"
			correctRes = `[{"label":"Tab1","rows":[["General"],["{title}","{description}","{status}","{priority}","{customer}"],["Details"],["{createdate}","{dynpriority}"]]},{"label":"Tab with a long name which should cut, does it?","rows":[["Calculations"],["{amount}","{price}","{total}","{product}"],["{discountprice}","{rebate}","{rebatetotal}","{discounttotal}"]]}]
`
		case 3:
			example = "Actions"
			correctRes = `{"header":[{"label":"Create new item","type":"Create","link":""}],"footer":[{"label":"Edit item","type":"edit","link":"http://localhost:8080/uwf-service/uwf/v1/items/{_id}/inline-edit"},{"label":"Move forward","type":"next","link":""},{"label":"Move backward","type":"previous","link":""}]}
`
		default:
			t.Fatal("for failed")
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/get?Get="+example, nil)
		if err != nil {
			t.Fatal(err)
		}

		get(w, r)
		if w.Code != http.StatusOK {
			t.Fatal(http.StatusOK, w.Code)
		}

		if strings.Compare(w.Body.String(), correctRes) != 0 {
			t.Fatal("FATAL:", w.Body.String(), "iteration:   ", i)
		}
	}
}
