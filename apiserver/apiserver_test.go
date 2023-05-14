package apiserver

import (
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
