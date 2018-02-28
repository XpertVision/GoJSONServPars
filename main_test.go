package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var testAPI API

func TestUpload(t *testing.T) {
	var err error

	file, err := os.Open("D:/json/example/test.json")

	if err != nil {
		t.Fatal("file didn't open")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.json")

	if err != nil {
		t.Fatal("CreateFromFile Error")
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal("copy")
	}

	err = writer.Close()
	if err != nil {
		t.Fatal("writer.Close() error")
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/upload", body)
	if err != nil {
		t.Fatal("Make new request error: ", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	testAPI.upload(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("response error:", w.Code)
	}

	fileExample, err := ioutil.ReadFile("D:/json/example/test.json")
	if err != nil {
		t.Fatal("Read example file to buffer error")
	}

	fileUpload, err := ioutil.ReadFile("D:/json/test.json")
	if err != nil {
		t.Fatal("Read uploaded file to buffer error")
	}

	if !bytes.Equal(fileExample, fileUpload) {
		t.Fatal("File uploaded incorrect")
	}
}

func TestGet(t *testing.T) {
	var example string
	var correctRes string

	for i := 0; i < 4; i++ {

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
			example = "incorrect request"
			correctRes = "BAD REQUEST:incorrect request"
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/get?Get="+example, nil)
		if err != nil {
			t.Fatal(err)
		}

		testAPI.get(w, r)
		if w.Code != http.StatusOK {
			t.Fatal(http.StatusOK, w.Code)
		}

		if strings.Compare(w.Body.String(), correctRes) != 0 {
			t.Fatal("FATAL:", w.Body.String(), "iteration:   ", i, correctRes)
		}
	}
}
