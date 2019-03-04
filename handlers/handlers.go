// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"

	"github.com/zew/urlshort/storages"
)

// EncodeHandler returns a HandlerFunc for encoding urls
// encapsulating the storage in a closure.
func EncodeHandler(st storages.IStore) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if url := r.FormValue("url"); url != "" {
			enc, err := st.Save(url)
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
func DecodeHandler(st storages.IStore) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		code := r.URL.Path[len("/dec/"):]
		url, err := st.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}
		w.Write([]byte(url))
		w.Write([]byte("<br>\n"))

		str, err := st.Dump(0, 100)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br>\n"))
		}
		w.Write([]byte(str))

	}
	return http.HandlerFunc(handleFunc)
}

// RedirectHandler returns a HandlerFunc for redirecting to the encoded URL,
// encapsulating the storage in a closure.
func RedirectHandler(st storages.IStore) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/r/"):]
		url, err := st.Load(code)
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
func DumpHandler(st storages.IStore) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		str, err := st.Dump(0, 100)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br>\n"))
		}
		w.Write([]byte(str))
	}
	return http.HandlerFunc(handleFunc)
}
