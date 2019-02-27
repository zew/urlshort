package storages

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
)

// This is ridiculous "benchmark" since it operates
// with only a single file in the directory.
// No surprise that this is fast.
func BenchmarkCode(b *testing.B) {
	storage, err := NewFS(filepath.Join(".", "url-shorts-benchmark"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		storage.enc()
	}
}

// A more realistic benchmark reveals,
// that iterating all existing files each time,
// to find the new maximum is a expletive deleted idea.
func BenchmarkStore(b *testing.B) {
	storage, err := NewFS(filepath.Join(".", "url-shorts-benchmark"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com", i))
	}
}
