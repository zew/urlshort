package storages

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type fileSystemT struct {
	Root string
	sync.RWMutex
}

// NewFS creates and returns a new file system
func NewFS(root string) (*fileSystemT, error) {
	s := &fileSystemT{}
	s.Root = root
	log.Printf("filesystem creation finished")
	return s, os.MkdirAll(s.Root, 0744)
}

// Enc returns the current number of files plus one - base36 encoded
func (s *fileSystemT) enc() (string, error) {
	files, err := ioutil.ReadDir(s.Root)
	if err != nil {
		return "", err
	}
	encode := strconv.FormatUint(uint64(len(files)+1), 36)
	return encode, nil
}

// Save writes the url to a file and returns the new files ordinal base36 encoded
func (s *fileSystemT) Save(url string) (string, error) {
	s.Lock()
	encode, err := s.enc()
	if err != nil {
		s.Unlock()
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(s.Root, encode), []byte(url), 0744)
	if err != nil {
		s.Unlock()
		return "", err
	}
	s.Unlock()
	// log.Printf("encoded %v to %v and saved it to file", url, encode)
	return encode, nil
}

// Load takes base36 encoded number, retrieves the corresponding file contents.
// We relent the Lock/Unlock around Readfile.
// Even if file creation coincided with file reading
// the result would be one crippled redirect
func (s *fileSystemT) Load(encode string) (string, error) {
	urlBytes, err := ioutil.ReadFile(filepath.Join(s.Root, encode))
	return string(urlBytes), err
}

// Dump dummy
func (s *fileSystemT) Dump(from, to int) (string, error) {
	return "", nil
}
