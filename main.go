package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sync"
)

//API struct than contains main mutex for actual JSON and JSON file
type API struct {
	jsn MyJSON
	mut sync.RWMutex
}

//States JSON SUBSTRUCT
type States struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

//Transitions JSON SUBSTRUCT
type Transitions struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"`
}

//Layout JSON SUBSTRUCT
type Layout struct {
	Label string     `json:"label"`
	Rows  [][]string `json:"rows"`
}

//Button JSON->Actions SUBSTRUCT
type Button struct {
	Label string `json:"label"`
	Mtype string `json:"type"`
	Link  string `json:"link"`
}

//Actions JSON SUBSTRUCT
type Actions struct {
	Header []Button `json:"header"`
	Footer []Button `json:"footer"`
}

//MyJSON JSON STRUCT
type MyJSON struct {
	States      []States      `json:"states"`
	Transitions []Transitions `json:"transitions"`
	Layout      []Layout      `json:"layout"`
	Actions     Actions       `json:"actions"`
}

// upload func for uploading JSON file to server
func (a *API) upload(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpJSON MyJSON

	defer func() {
		if err != nil {
			fmt.Println("error: ", err)
		}
	}()

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("BAD REQUEST: Parse form error"))
		return
	}

	file, fileType, err := r.FormFile("file")

	if err != nil {
		fmt.Println("ERROR FILE UPLOAD:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("BAD REQUEST: File upload error"))
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("ERROR FILE READ:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("BAD REQUEST: File read error"))
		return
	}

	err = json.Unmarshal(fileBytes, &tmpJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Internal Server Error: File parcing error"))
		return
	}

	a.mut.Lock()
	defer a.mut.Unlock()

	if reflect.DeepEqual(tmpJSON, a.jsn) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("You tried upload actual file. Reject."))
	}

	a.jsn = tmpJSON
	err = ioutil.WriteFile("D:/json/"+fileType.Filename, fileBytes, os.FileMode(os.O_WRONLY))
}

//get func for returning part of JSON file from server
func (a *API) get(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("BAD REQUEST: Parse form error"))
		return
	}

	getType := r.FormValue("Get")

	defer func() {
		if err != nil {
			fmt.Println("Send unswer error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	a.mut.RLock()
	defer a.mut.RUnlock()

	switch getType {
	case "States":
		err = json.NewEncoder(w).Encode(a.jsn.States)
	case "Transitions":
		err = json.NewEncoder(w).Encode(a.jsn.Transitions)
	case "Layout":
		err = json.NewEncoder(w).Encode(a.jsn.Layout)
	case "Actions":
		err = json.NewEncoder(w).Encode(a.jsn.Actions)
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("BAD REQUEST:" + getType))
	}
}

func main() {
	var err error

	var mainAPI API

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/upload", mainAPI.upload)
	http.HandleFunc("/get", mainAPI.get)

	err = server.ListenAndServe()

	if err != nil {
		fmt.Println("error:", err)
	}
}
