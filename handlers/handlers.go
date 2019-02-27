// Package handlers provides HTTP request handlers.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zew/urlshort/storages"
)

// EncodeHandler returns a HandlerFunc for encoding urls
// encapsulating the storage in a closure.
func EncodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if url := r.FormValue("url"); url != "" {
			enc, err := storage.Save(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte(enc))
		}
	}
	return http.HandlerFunc(handleFunc)
}

// DecodeHandler returns a HandlerFunc for retrieving urls by code
// encapsulating the storage in a closure.
func DecodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/dec/"):]
		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}
		w.Write([]byte(url))
	}
	return http.HandlerFunc(handleFunc)
}

// RedirectHandler returns a HandlerFunc for redirecting to the encoded URL,
// encapsulating the storage in a closure.
func RedirectHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/r/"):]
		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}
		http.Redirect(w, r, string(url), 301)
	}
	return http.HandlerFunc(handleFunc)
}

// DumpHandler dumps
func DumpHandler(l *leveldb.DB) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		iter := l.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			k := iter.Key()
			v := iter.Value()
			s := fmt.Sprintf("<a  href='/r/%s' target='red' > key %-20s => val %s  </a> <br>", k, k, v)
			w.Write([]byte(s + "\n"))

		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			w.Write([]byte("iter errors accumulated: " + err.Error() + "\n"))
		}

	}
	return http.HandlerFunc(handleFunc)
}
