package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/zew/urlshort/handlers"
	"github.com/zew/urlshort/storages"
)

func main() {

	//
	//
	store1, closeLevelDB, err := storages.NewLevelDB(filepath.Join(".", "storages", "level-db"))
	if err != nil {
		log.Fatal(err)
	}

	store2, err := storages.NewBoltDB(filepath.Join(".", "storages", "bolt-db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store2.DB.Close()
	_ = store2

	//
	//
	var activeStore storages.IStore
	activeStore = store1

	http.Handle("/enc", handlers.EncodeHandler(activeStore))
	http.Handle("/dec/", handlers.DecodeHandler(activeStore))

	http.Handle("/dump", handlers.DumpHandler(activeStore))

	http.Handle("/r/", handlers.RedirectHandler(activeStore))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	srv := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	// Graceful shutdown preparations
	sigquit := make(chan os.Signal, 0)
	signal.Notify(sigquit, os.Interrupt, os.Kill)
	closeHTTP := make(chan bool, 0)
	closeRes := make(chan bool, 0)

	// Start HTTP server
	go func() {
		log.Printf("starting HTTP Server. Listening at %q", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Printf("\t\t\t(HTTP server shutdown completed regularly)")
			} else {
				log.Printf("\t\t\t(HTTP server broke irregularly: %v)", err.Error())
			}
			closeHTTP <- true
		}
	}()

	// Activate this to test the freeing of resources, after the HTTP server crashed.
	// Press CTRL+C twice.
	if false {
		go func() {
			<-sigquit
			log.Printf("'crashing' HTTP server")
			srv.Close()
		}()
	}

	// Check for closing signal or HTTP server breakdown
	go func() {
		for {
			select {
			case sig := <-sigquit:
				log.Printf("\tcaught os stop signal: %+v", sig)
				ctx := context.Background()
				timeout, cancel := context.WithTimeout(ctx, 20*time.Second)
				if err := srv.Shutdown(timeout); err != nil {
					log.Printf("\t\tHTTP server could not be stopped or timeout: %v", err)
				} else {
					log.Printf("\t\tHTTP server stopped regularly")
				}
				cancel()
			case <-closeHTTP:
				log.Printf("\t\tcaught http server stop signal")
			}
			break
		}
		closeLevelDB()
		closeRes <- true
	}()

	log.Printf("close resources signal pending...")
	<-closeRes
	log.Printf("close resources signal processed")

}
