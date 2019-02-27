package storages

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func BenchmarkLevelDBStore(b *testing.B) {
	storage, err := NewLevelDB(filepath.Join(".", "tmp-leveldb-benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	defer storage.DB.Close()

	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
	}
}
func BenchmarkLevelDBRetrieve(b *testing.B) {
	storage, err := NewLevelDB(filepath.Join(".", "tmp-leveldb-benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	defer storage.DB.Close()
	b.Logf("Benachmarking %20v", b.N)

	// Fill
	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
	}

	// We only want to test retrieval
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = storage.Load(fmt.Sprintf("%v", b.N-i-1))
		if err == leveldb.ErrNotFound {
			b.Logf("Element not there %v", i)
			break
		}
		if err != nil {
			b.Fatal(err)
		}
	}
}
