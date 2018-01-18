package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

//GLOBAL STRUCT

type API struct {
	jsn          MyJSON
	fileNameMain string
	mut          sync.RWMutex
	file         []byte
	isParsed     bool
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

type Actions struct {
	Header []struct {
		Label string `json:"label"`
		Mtype string `json:"type"`
		Link  string `json:"link"` //,omitempty
	} `json:"header"`

	Footer []struct {
		Label string `json:"label"`
		Mtype string `json:"type"`
		Link  string `json:"link"`
	} `json:"footer"`
}

//STRUCT
type MyJSON struct {
	StatesRow      []States      `json:"states"`
	TransitionsRow []Transitions `json:"transitions"`
	LayoutRow      []Layout      `json:"layout"`
	ActionsRow     Actions       `json:"actions"`
}

//
//

func parseJSON(fileName string, a *API) error {
	var err error

	a.mut.RLock()
	jsByte, err := ioutil.ReadFile(fileName)
	a.mut.RUnlock()

	if err != nil {
		return errors.New("read file error")
	}

	a.mut.Lock()
	err = json.Unmarshal(jsByte, &a.jsn)
	a.mut.Unlock()

	if err != nil {
		fmt.Println("Unmarshal file error:", err)
		return errors.New("Unmarshal file error")
	}

	a.mut.Lock()
	a.isParsed = true
	a.mut.Unlock()

	return nil
}

func (a *API) upload(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
	}

	file, fileType, err := r.FormFile("file")

	if err != nil {
		fmt.Println("ERROR FILE UPLOAD:", err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("ERROR FILE READ:", err)
		return
	}

	/*a.mut.RLock()
	oldFileBytes, err := ioutil.ReadFile("D:/json/" + a.fileNameMain)
	a.mut.RUnlock()*/

	if err == nil {
		compareResult := bytes.Compare(a.file /*oldFileBytes*/, fileBytes)

		if compareResult == 0 {
			return
		}
	}

	a.mut.Lock()
	a.fileNameMain = fileType.Filename
	ioutil.WriteFile("D:/json/"+a.fileNameMain, fileBytes, os.FileMode(os.O_WRONLY))
	a.file = fileBytes
	a.mut.Unlock()
}

func (a *API) get(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
	}

	getType := r.FormValue("Get")

	/*if len(fileNameMain) == 0 {
		fileNameMain = "example/test.json"
	}*/

	if !a.isParsed {
		parseJSON("D:/json/"+a.fileNameMain, a)
	}

	a.mut.RLock()

	switch getType {
	case "States":
		json.NewEncoder(w).Encode(a.jsn.StatesRow)
	case "Transitions":
		json.NewEncoder(w).Encode(a.jsn.TransitionsRow)
	case "Layout":
		json.NewEncoder(w).Encode(a.jsn.LayoutRow)
	case "Actions":
		json.NewEncoder(w).Encode(a.jsn.ActionsRow)
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
