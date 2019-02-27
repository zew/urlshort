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
	store1, err := storages.NewFS(filepath.Join(".", "url-shorts"))
	if err != nil {
		log.Fatal(err)
	}
	_ = store1

	store2, err := storages.NewBoltDB(filepath.Join(".", "bolt-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store2.DB.Close()
	_ = store2

	store3, err := storages.NewLevelDB(filepath.Join(".", "level-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store3.DB.Close()
	_ = store3

	//
	//
	var activeStore storages.IStore
	activeStore = store3

	http.Handle("/enc", handlers.EncodeHandler(activeStore))
	http.Handle("/dec/", handlers.DecodeHandler(activeStore))

	http.Handle("/dump", handlers.DumpHandler(activeStore))

	http.Handle("/r/", handlers.RedirectHandler(activeStore))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
