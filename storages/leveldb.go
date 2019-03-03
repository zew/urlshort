package storages

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type levelDBT struct {
	Root string
	*leveldb.DB

	Lf *os.File    // "transaction log"
	L  *log.Logger // "transaction logger"
}

// NewLevelDB creates and returns a new level db
func NewLevelDB(filename string) (*levelDBT, func(), error) {
	db, err := leveldb.OpenFile(filename, nil)
	if err != nil {
		log.Fatal(err)
	}
	l := &levelDBT{}
	l.Root = filename
	l.DB = db
	l.Lf, err = os.OpenFile(filepath.Join(filename, "insert.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile: %v", err)
	}
	l.L = log.New(l.Lf, "", log.Ldate)

	//
	// Instead of calling:
	//     defer l.DB.Close()
	//     defer l.Lf.Close()
	//
	// Only *one* receiving channel is triggered by the OS.
	// And we already need such channel in main().
	// Thus we cannot spawn a goroutine here, waiting for it.
	// Instead we return a 'cancel' func to be called later from main().
	closingFunc := func() {
		log.Printf("closing leveldb files start")
		l.DB.Close()
		l.Lf.Close()
		log.Printf("closing leveldb files stop")
	}

	log.Printf("levelDB creation finished")
	return l, closingFunc, nil

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
	p = strings.Join(strings.Fields(p), "") // remove all spaces
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
	if err == nil {
		l.L.Printf("%v %v", encode, url)
	}
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

// Dump writes
func (l *levelDBT) Dump(from, to int) (string, error) {
	var bf bytes.Buffer
	iter := l.NewIterator(nil, nil)
	ctr := -1
	for iter.Next() {
		ctr++
		if ctr < from {
			continue
		}
		if ctr > to {
			break
		}

		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		k := iter.Key()
		v := iter.Value()
		s := fmt.Sprintf("<a  href='/r/%s' target='red' > key %-20s => val %s  </a> <br>", k, k, v)
		bf.Write([]byte(s + "\n"))

	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		bf.Write([]byte("iter errors accumulated: " + err.Error() + "\n"))
	}

	return bf.String(), nil
}
