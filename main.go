package main

import (
	"fmt"
	"github.com/sasidakh/crypt-server/pkg/crypt"
	"github.com/sasidakh/crypt-server/pkg/router"
	"log"
	"net/http"
	"os"
)

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	dcrptr := crypt.NewDecrypter(os.Getenv("PRI_KEY"), os.Getenv("PASS"))
	// 1 MB small machine
	r.ParseMultipartForm(1 << 20)
	file, _, err := r.FormFile("file")
	must(err)
	defer file.Close()
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	decStr := dcrptr.Decrypt(file)
	w.Write([]byte(decStr))
}

func main() {
	r := router.NewRouter()
	s := http.Server{
		Addr:    ":8035",
		Handler: r,
	}
	r.SetHandlerFunc("POST", "/upload", uploadFile)
	fmt.Printf("Starting server on %s\n", s.Addr)
	must(s.ListenAndServe())
}
