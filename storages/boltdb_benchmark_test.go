package storages

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
	bolt "go.etcd.io/bbolt"
)

const dbName = "tmp-boltdb-benchmark"

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

// BenchmarkBoltDBStoreBulk illustrates batching of inserts
// and syncing only every thousand operations.
// performance reaches 300.000 inserts per second
func BenchmarkBoltDBStoreBulk(b *testing.B) {
	storage, closeFunc, err := NewBoltDB(filepath.Join(".", dbName))
	if err != nil {
		b.Fatal(err)
	}
	defer closeFunc()

	for i := 0; i < b.N; i++ {
		ufunc := func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			if b == nil {
				return fmt.Errorf("bucket does not exist")
			}
			for j := 0; j < 1000; j++ {
				k, v := []byte(fmt.Sprintf(" %04v-%04v", i, j)), []byte("blupp")
				err = b.Put(k, v)
				if err != nil {
					return fmt.Errorf("put error %s - %s : %v", k, v, err)
				}
			}
			return nil
		}
		err = storage.Update(ufunc)
		if err != nil {
			b.Fatal(err)
		}
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
