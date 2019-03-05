package storages

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type boltDBT struct {
	Root string
	*bolt.DB
}

// NewBoltDB creates and returns a new level db
func NewBoltDB(filename string) (*boltDBT, func(), error) {
	db, err := bolt.Open(filepath.Join(".", filename), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	l := &boltDBT{}
	l.Root = filename
	l.DB = db

	// instead of defer b.DB.Close()
	closingFunc := func() {
		// log.Printf("closing boltdb files start")
		l.DB.Close()
		log.Printf("closing boltdb files stop")
	}

	ufunc := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			if err != bolt.ErrBucketExists {
				return fmt.Errorf("bucket creation failed: %v", err)
			}
		}
		return nil
	}
	err = l.DB.Update(ufunc)

	log.Printf("boltDB  creation finished")
	return l, closingFunc, err
}

// Enc returns the h argument of a URL
func (l *boltDBT) enc(strURL string) (string, error) {
	urll, err := url.Parse(strURL)
	if err != nil {
		return "", err
	}
	p := urll.Query().Get("h")
	if p == "" {
		return "", errors.Errorf("parameter h cannot be empty")
	}
	p = strings.Join(strings.Fields(p), "") // remove all spaces
	// encode := strconv.FormatUint(uint64(next+1), 36)
	return p, nil
}

// Save writes the url to database and returns the hash
func (l *boltDBT) Save(url string) (string, error) {
	encode, err := l.enc(url)
	if err != nil {
		return encode, err
	}
	ufunc := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		k, v := []byte(encode), []byte(url)
		err = b.Put(k, v)
		if err != nil {
			return fmt.Errorf("put error %s - %s : %v", k, v, err)
		}
		// err = tx.Commit() // this is an auto-commit environment; no tx handling allowed or necessary
		return nil
	}
	err = l.DB.Update(ufunc)
	// log.Printf("boltdb updated %v - %v - %v ", url, encode, err)
	return encode, err
}

// Load takes hash and retrieves the corresponding url.
func (l *boltDBT) Load(encode string) (string, error) {
	var s string
	ufunc := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		v := b.Get([]byte(encode))
		s = fmt.Sprintf("%s", v) // force new allocation
		return nil
	}
	err := l.DB.View(ufunc)
	// log.Printf("boltdb found %v for %v", s, encode)
	return s, err
}

// Dump writes
func (l *boltDBT) Dump(from, to int) (string, error) {
	// implementation; todo
	var bf bytes.Buffer

	ufunc := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		c := b.Cursor()
		c.First()
		ctr := 0
		for {
			ctr++
			k, v := c.Next()
			if k == nil && v == nil {
				break
			}
			bf.Write(k)
			bf.Write([]byte(" "))
			bf.Write(v)
			bf.Write([]byte("<br>\n"))
			log.Printf("boltdb dump %s for %s", k, v)
			log.Printf("boltdb ctr %2v", ctr)
		}
		return nil
	}
	err := l.DB.View(ufunc)
	return bf.String(), err
}
