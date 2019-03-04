package storages

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

const dbName = "tmp-boltdb-benchmark-4"

func BenchmarkBoltDBStore(b *testing.B) {
	storage, closeFunc, err := NewBoltDB(filepath.Join(".", dbName))
	if err != nil {
		b.Fatal(err)
	}
	defer closeFunc()

	for i := 0; i < b.N; i++ {
		storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
	}
}
func BenchmarkBoltDBRetrieve(b *testing.B) {
	storage, closeFunc, err := NewBoltDB(filepath.Join(".", dbName))
	if err != nil {
		b.Fatal(err)
	}
	defer closeFunc()
	b.Logf("Benachmarking %20v", b.N)

	// Fill
	if true {
		// For some reason, Stop- and Start makes it run forever
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			storage.Save(fmt.Sprintf("http://www.abc%v.com?h=%v", i, i))
		}
		b.StartTimer()
	}

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
