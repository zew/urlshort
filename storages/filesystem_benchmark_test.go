package storages

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
)

//
func BenchmarkCode(b *testing.B) {
	storage, err := New(filepath.Join(".", "url-shorts-benchmark"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		storage.enc()
	}
}
func BenchmarkStore(b *testing.B) {
	storage, err := New(filepath.Join(".", "url-shorts-benchmark"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com", i))
	}
}
