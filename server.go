package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// parse command line arguments
	dbname := flag.String("name", "./testdb", "Open the database with the specified name")
	addr := flag.String("addr", ":5050", "listen address")

	flag.Parse()

	// open database and start backend
	backend, err := NewBackend(*dbname)
	if err != nil {
		log.Fatal("can't open database", err)
	}

	backend.Start()
	defer backend.Shutdown()

	// listen and serve http
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			HandleGet(w, r, backend)
		case "POST", "PUT":
			HandlePut(w, r, backend)
		case "DELETE":
			HandleDelete(w, r, backend)
		}
	})

	log.Printf("Server listening on : %s", *addr)
	http.ListenAndServe(*addr, nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request, backend *Backend) {
	key := r.URL.Path[len("/"):]

	val, err := backend.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if val == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func HandlePut(w http.ResponseWriter, r *http.Request, backend *Backend) {
	key := r.URL.Path[len("/"):]

	defer r.Body.Close()
	val, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	backend.Put(key, string(val))

	w.WriteHeader(http.StatusCreated)
}

func HandleDelete(w http.ResponseWriter, r *http.Request, backend *Backend) {
	key := r.URL.Path[len("/"):]

	backend.Delete(key)

	w.WriteHeader(http.StatusNoContent)
}
