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

	"github.com/gobuffalo/packr"
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
	Dir          *string
	Port         *string
	StaticBox    = packr.NewBox("./static/")
)

func main() {
	Dir = flag.String("d", "./", "Directory to serve")
	Port = flag.String("p", "8080", "Port to serve")
	flag.Parse()

	WalkTheWalk()

	r := mux.NewRouter()
	r.Use(RefreshDirList)
	r.HandleFunc("/", Index)
	serveFile := r.PathPrefix("/items/").Subrouter()
	serveFile.HandleFunc("/{id}/", ServeFile).Methods("GET")

	items := r.PathPrefix("/items/").Subrouter()
	items.Use(SetJson)
	items.HandleFunc("/", ServeList).Methods("GET")
	items.HandleFunc("/{id}/", DeleteFile).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(StaticBox)))

	s := http.Server{
		Addr:           ":" + *Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Serving on %s dir %s", *Port, *Dir)
	log.Fatal(s.ListenAndServe())
}

func RefreshDirList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WalkTheWalk()
		next.ServeHTTP(w, r)
	})
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

func WalkTheWalk() {
	FileNameList = nil
	err := filepath.Walk(*Dir, Listdir)
	if err != nil {
		log.Fatal(err)
	}
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
	tpl, err := StaticBox.FindString("index.html")
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	} else {
		t, _ := template.New("index").Parse(tpl)
		t.Execute(w, nil)
	}
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

	if len(FileNameList) >= intId {
		http.ServeFile(w, r, FileNameList[intId].Href)
	} else {
		http.NotFound(w, r)
	}
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
			WalkTheWalk()
		}
	}
}
