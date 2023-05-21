package apiserver

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"testing"
)

func TestBase(t *testing.T) {

	dir, err := dashboard.ReadDir("dashboard/build")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		fmt.Println(file.Name())
	}

	serverRoot, err := fs.Sub(dashboard, "dashboard/build")
	if err != nil {
		log.Fatal(err)
	}

	index, err := serverRoot.Open("index.html")
	if err != nil {
		log.Fatal(err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatal(err)
	}

	name := stat.Name()
	got := name

	expected := "index.html"

	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestJSON(t *testing.T) {

	type message struct {
		Id    string `json:"id"`
		Value string `json:"value"`
	}

	type rlist struct {
		R []message
	}

	var Response rlist
	Response.R = append(Response.R, message{"monitorOpCompleted", "0"})
	Response.R = append(Response.R, message{"monitorOpFailed", "0"})
	Response.R = append(Response.R, message{"monitorOpAborted", "0"})
	Response.R = append(Response.R, message{"monitorOpCancelled", "0"})
	Response.R = append(Response.R, message{"monitorOpTimeout", "0"})

	for _, v := range Response.R {
		vjson, _ := json.Marshal(v)
		fmt.Println(v.Id, v.Value, string(vjson))
	}

	expected := `[{"id":"monitorOpCompleted","value":"0"},{"id":"monitorOpFailed","value":"0"},{"id":"monitorOpAborted","value":"0"},{"id":"monitorOpCancelled","value":"0"},{"id":"monitorOpTimeout","value":"0"}]`
	got, _ := json.Marshal(Response.R)

	if string(got) != expected {
		t.Errorf("got %s, expected %s", string(got), expected)
	}
}
