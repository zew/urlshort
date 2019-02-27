package storages

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

type boltDBT struct {
	Root string
	*bolt.DB
}

// NewBoltDB creates and returns a new level db
func NewBoltDB(filename string) (*boltDBT, error) {
	db, err := bolt.Open(filepath.Join(".", filename), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	b := &boltDBT{}
	b.Root = filename
	b.DB = db
	// defer b.DB.Close()
	log.Printf("boltDB creation finished")
	return b, nil
}

func (b *boltDBT) run() {
	ufunc := func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		k, v := []byte("answer"), []byte("42")
		err = b.Put(k, v)
		if err != nil {
			return fmt.Errorf("put error %s - %s : %s", k, v, err)
		}
		v2 := b.Get(k)
		log.Printf("found %v", v2)
		return nil
	}
	b.DB.Update(ufunc)
	log.Printf("boltDB update finished")
}
