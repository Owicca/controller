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
	Name string `json:"name"`
	Href string `json:"href"`
}

type JR struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
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
	r.HandleFunc("/", Index)
	items := r.PathPrefix("/items/").Subrouter()
	items.Use(SetJson)

	items.HandleFunc("/", ServeList).Methods("GET")
	items.HandleFunc("/{id}/", ServeFile).Methods("GET")
	items.HandleFunc("/{id}/", DeleteFile).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

func SetMime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(*r.URL)
		//w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func SetJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
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
	t, _ := template.ParseFiles("index.tpl")
	t.Execute(w, nil)
}

func ServeList(w http.ResponseWriter, r *http.Request) {
	res := JR{Success: true, Data: FileNameList, Error: nil}

	js, err := json.Marshal(res)
	if err != nil {
		res.Success = false
		log.Println(res)
	} else {
		w.Write(js)
	}
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	intId, _ := strconv.Atoi(params["id"])

	http.ServeFile(w, r, FileNameList[intId].Href)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	response := JR{Success: false, Data: nil, Error: nil}
	params := mux.Vars(r)
	intId, _ := strconv.Atoi(params["id"])

	pathErr := os.Remove(FileNameList[intId].Href)
	if pathErr != nil {
		response.Error = pathErr.Error()
		js, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", js)
	} else {
		response.Success = true
		res, err := json.Marshal(response)
		if err != nil {
			response.Success = false
			response.Error = err.Error()
			log.Println(response)
		} else {
			response.Success = true
			response.Data = res
			js, _ := json.Marshal(response)
			fmt.Fprintf(w, "%s", js)
		}
	}
}
