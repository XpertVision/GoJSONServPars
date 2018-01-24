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

//GLOBAL STRUCT

type API struct {
	jsn MyJSON
	mut sync.RWMutex
}

//

//JSON SUBSTRUCT & STRUCT
type States struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type Transitions struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"`
}

type Layout struct {
	Label string     `json:"label"`
	Rows  [][]string `json:"rows"`
}

type Button struct {
	Label string `json:"label"`
	Mtype string `json:"type"`
	Link  string `json:"link, omitempty"` //,omitempty
}

type Actions struct {
	Header []Button `json:"header"`
	Footer []Button `json:"footer"`
}

//STRUCT
type MyJSON struct {
	States      []States      `json:"states"`
	Transitions []Transitions `json:"transitions"`
	Layout      []Layout      `json:"layout"`
	Actions     Actions       `json:"actions"`
}

//
//

/*func parseJSON(fileName string, a *API) error {
	var err error

	a.mut.RLock()
	file, err := ioutil.ReadFile(fileName)
	a.mut.RUnlock()

	if err != nil {
		return errors.New("read file error")
	}

	a.mut.Lock()
	err = json.Unmarshal(file, &a.jsn)

	if err != nil {
		fmt.Println("Unmarshal file error:", err)
		return errors.New("Unmarshal file error")
	}

	a.mut.Unlock()

	return nil
}*/

func (a *API) upload(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpAPI API

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: Parse form error"))
	}

	file, fileType, err := r.FormFile("file")

	if err != nil {
		fmt.Println("ERROR FILE UPLOAD:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: File upload error"))
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("ERROR FILE READ:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: File read error"))
		return
	}

	err = json.Unmarshal(fileBytes, &tmpAPI.jsn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error: File parcing error"))
		return
	}

	a.mut.Lock()
	if !reflect.DeepEqual(tmpAPI.jsn, a.jsn) {
		a.jsn = tmpAPI.jsn
		ioutil.WriteFile("D:/json/"+fileType.Filename, fileBytes, os.FileMode(os.O_WRONLY))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("You tried upload actual file. Reject."))
	}
	a.mut.Unlock()
}

func (a *API) get(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: Parse form error"))
	}

	getType := r.FormValue("Get")

	a.mut.RLock()

	switch getType {
	case "States":
		json.NewEncoder(w).Encode(a.jsn.States)
	case "Transitions":
		json.NewEncoder(w).Encode(a.jsn.Transitions)
	case "Layout":
		json.NewEncoder(w).Encode(a.jsn.Layout)
	case "Actions":
		json.NewEncoder(w).Encode(a.jsn.Actions)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST:" + getType))
	}

	a.mut.RUnlock()
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
