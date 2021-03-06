package storages

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

const fn = "tmp-leveldb-benchmark"

func BenchmarkLevelDBStore(b *testing.B) {
	storage, closeFunc, err := NewLevelDB(filepath.Join(".", fn))
	if err != nil {
		b.Fatal(err)
	}
	defer closeFunc()

	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
	}

}
func BenchmarkLevelDBRetrieve(b *testing.B) {
	storage, closeFunc, err := NewLevelDB(filepath.Join(".", fn))
	if err != nil {
		b.Fatal(err)
	}
	defer closeFunc()
	b.Logf("Benachmarking %20v", b.N)

	// Fill
	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
	}

	// st, _ := storage.GetProperty("leveldb.stats")
	// log.Printf("stats: %v", st)

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
