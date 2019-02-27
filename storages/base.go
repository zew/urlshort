// Package storages allows multiple implementation on how to store short URLs.
package storages

// IStore is the interface for saving and retrieving URLs
type IStore interface {
	Save(string) (string, error)
	Load(string) (string, error)
	Dump(int, int) (string, error)
}
