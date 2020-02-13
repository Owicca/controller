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

	"github.com/Owicca/controller/models/dir"
	"github.com/Owicca/controller/models/file"
	"github.com/Owicca/controller/models/response"
	"github.com/Owicca/controller/models/walker"

	"github.com/gorilla/mux"
)

var (
	Dir         *string
	Port        *string
	Environment *string
	Walker      *walker.Walker
)

func main() {
	Dir = flag.String("d", "./", "Directory to serve")
	Port = flag.String("p", "8080", "Port to serve")
	Environment = flag.String("e", "DEVEL", "Environment")
	flag.Parse()

	if port, check := os.LookupEnv("CONTROLLER_PORT"); check == true {
		*Port = port
	}
	if dir, check := os.LookupEnv("CONTROLLER_DIR"); check == true {
		*Dir = dir
	}
	if env, check := os.LookupEnv("CONTROLLER_ENV"); check == true {
		*Environment = env
	}

	Walker = walker.NewWalker()
	Walker.ParsePath(Dir)

	r := mux.NewRouter()
	r.Use(RefreshDirList)
	r.HandleFunc("/", Index)
	items := r.PathPrefix("/items/").Subrouter()
	items.Use(SetJson)
	items.HandleFunc("/", ServeList).Methods("GET")
	items.HandleFunc("/{pseudoname}/", ServeFile).Methods("GET")
	items.HandleFunc("/{pseudoname}/", DeleteFile).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	s := http.Server{
		Addr:           "0.0.0.0:" + *Port,
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
			Walker.TTL = 10
		}
		Walker.TTL--
		next.ServeHTTP(w, r)
	})
}

//func SetMime(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Println(*r.URL)
//		w.Header().Set("Content-Type", "application/json")
//		next.ServeHTTP(w, r)
//	})
//}

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

func ServeFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["pseudoname"]
	log.Println(param)
	for _, child := range Walker.FSTree.(dir.Dir).Children {
		if fl, ok := child.(file.File); ok {
			if fl.Info.PseudoName == param {
				log.Println(string(child.GetPath()))
				http.ServeFile(w, r, string(child.GetPath()))
				break
			}
		} else if child.(dir.Dir).Info.PseudoName == param {
			log.Println(string(child.GetPath()))
			http.ServeFile(w, r, string(child.GetPath()))
			break
		}
	}
	http.NotFound(w, r)
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
