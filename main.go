package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type file struct {
	Name string
	Href string
}

type JR struct {
	Success bool
	Error   string
}

var (
	FileNameList []file
)

func main() {
	dir := flag.String("d", "./", "Directory to serve")
	port := flag.String("p", "8080", "Port to serve")
	flag.Parse()

	err := filepath.Walk(*dir, Listdir)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/home/", Index)
	r.HandleFunc("/{id}/", ServeFile).PathPrefix("/v/")
	r.HandleFunc("/{id}/", DeleteFile).PathPrefix("/d/")

	s := http.Server{
		Addr:           ":" + *port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Serving on %s dir %s", *port, *dir)
	log.Fatal(s.ListenAndServe())
}

func Listdir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		FileNameList = append(FileNameList, file{
			Name: info.Name(),
			Href: path,
		})
	}

	return nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := FileNameList
	t, _ := template.ParseFiles("index.tpl")
	t.Execute(w, data)
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	intId, _ := strconv.Atoi(params["id"])

	http.ServeFile(w, r, FileNameList[intId].Href)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	response := JR{Success: false, Error: nil}
	params := mux.Vars(r)
	intId, _ := strconv.Atoi(params["id"])

	pathErr := os.Remove(FileNameList[intId].Href)
	if pathErr != nil {
		response.Error = pathErr
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
	}
	fmt.Fprintf(w, "%s", res)
}
