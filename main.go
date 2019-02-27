package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zew/urlshort/handlers"
	"github.com/zew/urlshort/storages"
)

func main() {

	//
	//

	//
	storage1, err := storages.NewFS(filepath.Join(".", "url-shorts"))
	if err != nil {
		log.Fatal(err)
	}
	_ = storage1

	storage2, err := storages.NewBoltDB(filepath.Join(".", "bolt-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer storage2.DB.Close()
	_ = storage2

	storage3, err := storages.NewLevelDB(filepath.Join(".", "level-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer storage3.DB.Close()
	_ = storage3

	//
	//
	var activeStorage storages.IStorage
	activeStorage = storage3

	http.Handle("/enc", handlers.EncodeHandler(activeStorage))
	http.Handle("/dec", handlers.DecodeHandler(activeStorage))
	http.Handle("/dump", handlers.DumpHandler(storage3.DB))

	http.Handle("/r/", handlers.RedirectHandler(activeStorage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
