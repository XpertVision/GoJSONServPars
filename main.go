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

//GLOBAL JSON VAR
var jsn MyJSON
var fileNameMain string
var mut sync.RWMutex

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

func parseJSON(fileName string) error {
	var err error

	mut.RLock()
	jsByte, err := ioutil.ReadFile(fileName)
	mut.RUnlock()

	if err != nil {
		return errors.New("read file error")
	}

	mut.Lock()

	err = json.Unmarshal(jsByte, &jsn)

	mut.Unlock()

	if err != nil {
		fmt.Println("Unmarshal file error:", err)
		return errors.New("Unmarshal file error")
	}

	return nil
}

func upload(w http.ResponseWriter, r *http.Request) {
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

	mut.RLock()
	oldFileBytes, err := ioutil.ReadFile("D:/json/" + fileNameMain)
	mut.RUnlock()

	if err == nil {
		compareResult := bytes.Compare(oldFileBytes, fileBytes)

		if compareResult == 0 {
			return
		}
	}

	mut.Lock()
	fileNameMain = fileType.Filename
	ioutil.WriteFile("D:/json/"+fileNameMain, fileBytes, os.FileMode(os.O_WRONLY))
	mut.Unlock()
}

func get(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()

	if err != nil {
		fmt.Println("parse err:", err)
	}

	getType := r.FormValue("Get")

	/*if len(fileNameMain) == 0 {
		fileNameMain = "example/test.json"
	}*/

	parseJSON("D:/json/" + fileNameMain)

	mut.RLock()

	switch getType {
	case "States":
		json.NewEncoder(w).Encode(jsn.StatesRow)
	case "Transitions":
		json.NewEncoder(w).Encode(jsn.TransitionsRow)
	case "Layout":
		json.NewEncoder(w).Encode(jsn.LayoutRow)
	case "Actions":
		json.NewEncoder(w).Encode(jsn.ActionsRow)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST:" + getType))
	}

	mut.RUnlock()
}

func main() {
	var err error

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/get", get)

	err = server.ListenAndServe()

	if err != nil {
		fmt.Println("error:", err)
	}
}
