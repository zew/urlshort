package storages

import (
	"fmt"
	"log"
	"net/url"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type levelDBT struct {
	Root string
	*leveldb.DB
}

// NewLevelDB creates and returns a new level db
func NewLevelDB(filename string) (*levelDBT, error) {
	db, err := leveldb.OpenFile(filename, nil)
	if err != nil {
		log.Fatal(err)
	}
	l := &levelDBT{}
	l.Root = filename
	l.DB = db
	// defer l.DB.Close()
	log.Printf("levelDB creation finished")
	return l, nil
}

// Enc returns the h argument of a URL
func (l *levelDBT) enc(strURL string) (string, error) {
	urll, err := url.Parse(strURL)
	if err != nil {
		return "", err
	}
	p := urll.Query().Get("h")
	if p == "" {
		return "", errors.Errorf("parameter h cannot be empty")
	}
	// encode := strconv.FormatUint(uint64(next+1), 36)
	return p, nil
}

// Save writes the url to database and returns the hash
func (l *levelDBT) Save(url string) (string, error) {
	encode, err := l.enc(url)
	if err != nil {
		return encode, err
	}
	err = l.DB.Put([]byte(encode), []byte(url), nil)
	return encode, err
}

// Load takes hash and retrieves the corresponding url.
func (l *levelDBT) Load(encode string) (string, error) {
	// urlBytes, err := ioutil.ReadFile(filepath.Join(s.Root, encode))
	data, err := l.DB.Get([]byte(encode), nil)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s", data) // force new allocation
	return s, nil

}

// Dump dumps
func (l *levelDBT) Dump() {

}
