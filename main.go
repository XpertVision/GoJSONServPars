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
		Link  string `json:"link,omitempty"`
	} `json:"header"`

	Footer []struct {
		Label string `json:"label"`
		Mtype string `json:"type"`
		Link  string `json:"link,omitempty"`
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

	mut.RLock()
	jsByte, err := ioutil.ReadFile(fileName)
	mut.RUnlock()

	if err != nil {
		return errors.New("read file error")
	}

	mut.Lock()

	errUnm := json.Unmarshal(jsByte, &jsn)

	mut.Unlock()

	if errUnm != nil {
		fmt.Println("Unmarshal file error:", errUnm)
		return errors.New("Unmarshal file error")
	}

	return nil
}

func upload(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()

	if parseErr != nil {
		fmt.Println("parse err:", parseErr)
	}

	file, fileType, er := r.FormFile("file")

	if er != nil {
		fmt.Println("ERROR FILE UPLOAD:", er)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if er != nil {
		fmt.Println("ERROR FILE READ:", err)
		return
	}

	mut.RLock()
	oldFileBytes, errRF := ioutil.ReadFile("D:/json/" + fileNameMain)
	mut.RUnlock()

	if errRF == nil {
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
	parseErr := r.ParseForm()

	if parseErr != nil {
		fmt.Println("parse err:", parseErr)
	}

	getType := r.FormValue("Get")

	if len(fileNameMain) == 0 {
		fileNameMain = "test.json"
	}

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
		fmt.Println("Invalid value or empty")
		return
	}

	mut.RUnlock()
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/get", get)

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("error:", err)
	}
}
