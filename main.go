package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"errors"

	"github.com/Owicca/controller/models/dir"
	"github.com/Owicca/controller/models/file"
	"github.com/Owicca/controller/models/response"
	"github.com/Owicca/controller/models/walker"
	"github.com/Owicca/controller/config"

	"github.com/gorilla/mux"
	"path/filepath"
)

var (
	Dir         *string
	Port        *string
	Host 		*string
	Environment *string
	MimeTypes   map[string]string
	Walker      *walker.Walker
	DummyPseudo = "o4FWKxKQkj"
)

func main() {
	MimeTypes = config.NewMimeTypes()

	if port, check := os.LookupEnv("CONTROLLER_PORT"); check == true {
		*Port = port
	} else {
		*Port = "8080"
	}
	if dir, check := os.LookupEnv("CONTROLLER_DIR"); check == true {
		*Dir = dir
	} else {
		*Dir = "./"
	}
	if env, check := os.LookupEnv("CONTROLLER_ENV"); check == true {
		*Environment = env
	} else {
		*Environment = "DEVEL"
	}
	if host, check := os.LookupEnv("CONTROLLER_HOST"); check == true {
		*Host = host
	} else {
		*Host = "127.0.0.1"
	}

	Dir = flag.String("d", *Dir, "Directory to serve")
	Port = flag.String("p", *Port, "Port to serve")
	Environment = flag.String("e", *Environment, "Environment")
	Host = flag.String("h", *Host, "Host")
	flag.Parse()

	Walker = walker.NewWalker()
	Walker.ParsePath(Dir)

	r := mux.NewRouter()
	r.Use(RefreshDirList)
	r.HandleFunc("/", Index)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	items := r.PathPrefix("/items/").Subrouter()
	items.HandleFunc("/{pseudoname}/", ServeFile).Methods("GET")//ServeFile returns a byte stream, not a json
	items.Use(SetJson)
	items.HandleFunc("/", ServeList).Methods("GET")
	items.HandleFunc("/{pseudoname}/", DeleteFile).Methods("DELETE")

	s := http.Server{
		Addr:           *Host + ":" + *Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Serving on %s dir %s", *Port, *Dir)
	log.Fatal(s.ListenAndServe())
}

// refresh Walker.FSTree every 10 requests
func RefreshDirList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Walker.TTL == 0 {
			Walker.ParsePath(Dir)
			if *Environment == "DEVEL" {
				log.Println("Refreshed Walker.FSTree")
			}
			Walker.TTL = 10
		}
		Walker.TTL--
		next.ServeHTTP(w, r)
	})
}

func SetJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// security risk to allow cors from "*" TODO
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, nil)
}

// json the FSTree and format the response
func ServeList(w http.ResponseWriter, r *http.Request) {
	res := response.Res{Success: true, Data: nil, Error: nil}
	res.Data = Walker.FSTree

	js, err := json.Marshal(res)
	if err != nil {
		res.Success = false
	}
	w.Write(js)
}

/*
* find a file
* set response mimetype based on file extension
* server file content
*/
func ServeFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["pseudoname"]
	splitPseudo := Splitter(param)
	filePath, fErr := FindFile(Walker.FSTree.(dir.Dir), splitPseudo)
	if fErr != nil {
		log.Println(fErr)
		http.NotFound(w, r)
	} else {
		ext := filepath.Ext(filePath)
		mimeType := "application/octet-stream"
		if ext != "" {
			if mime, ok := MimeTypes[ext]; ok {
				mimeType = mime
			}
		}

		log.Println("Serve: ", filePath, "\nMime: ", mimeType)
		w.Header().Set("Content-Type", mimeType)
		http.ServeFile(w, r, filePath)
	}
}

/*
* return file path
* or recurse in folder to look for file path
*/
func FindFile(Parent dir.Dir, Paths []string) (string, error) {
	if len(Paths) > 0 {
		child, ok := Parent.Children[Paths[0]]
		if ok {
			file, check := child.(file.File)
			if check {
				// log.Println("Found: ", file.Info.Name, "\nPath: ", string(file.GetPath()))
				return string(file.GetPath()), nil
			} else {
				// log.Println("Recurse in dir: ", child.(dir.Dir).Info.Name)
				FindFile(child.(dir.Dir), Paths[1:])
			}
		}
	}

	return "", errors.New("File not found")
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	response := response.Res{Success: false, Data: nil, Error: nil}
	//params := mux.Vars(r)
	//intId, _ := strconv.Atoi(params["id"])

	js, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", js)
	//	pathErr := os.Remove(FileNameList[intId].Href)
	//	if pathErr != nil {
	//		response.Error = pathErr.Error()
	//		js, _ := json.Marshal(response)
	//		fmt.Fprintf(w, "%s", js)
	//	} else {
	//		response.Success = true
	//		res, err := json.Marshal(response)
	//		if err != nil {
	//			response.Success = false
	//			response.Error = err.Error()
	//			log.Println(response)
	//		} else {
	//			response.Success = true
	//			response.Data = res
	//			js, _ := json.Marshal(response)
	//			fmt.Fprintf(w, "%s", js)
	//			WalkTheWalk()
	//		}
	//	}
}

func Splitter(Pseudo string) []string {
	var items []string
	lowerLimit := 1
	length := 5
	for lowerLimit - 1 + length <= len(Pseudo) {
		items = append(items, Pseudo[lowerLimit-1:lowerLimit-1 + length])
		lowerLimit = lowerLimit+5
	}

	return items
}